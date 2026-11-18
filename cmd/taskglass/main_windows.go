//go:build windows

// Command taskglass makes the Windows taskbar transparent, blurred or acrylic.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/AnderCMD/TaskGlass/internal/app"
	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/taskbar"
	"github.com/AnderCMD/TaskGlass/internal/winglass"
)

// version is overridden at build time: -ldflags "-X main.version=1.2.3".
var version = "dev"

const (
	mutexName = `Local\TaskGlass.Singleton`
	quitName  = `Local\TaskGlass.Quit`
)

func main() {
	var (
		showVersion = flag.Bool("version", false, "print version and exit")
		cfgFlag     = flag.String("config", "", "path to config.json (default %APPDATA%\\TaskGlass\\config.json)")
		reset       = flag.Bool("reset", false, "rewrite the config file with defaults and exit")
		validate    = flag.Bool("validate", false, "check the config file and exit (non-zero on errors)")
		quit        = flag.Bool("quit", false, "ask a running TaskGlass to exit")
		verbose     = flag.Bool("verbose", false, "log to stderr as well as the log file")
		restore     = flag.Bool("restore", false, "put the taskbar back to the stock look and exit (use after a crash)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "TaskGlass %s: a transparent taskbar for Windows 10 and 11\n\nUsage: taskglass [flags]\n\n", version)
		flag.PrintDefaults()
	}
	if len(os.Args) > 1 {
		attachConsole() // so flags like --validate can print from a GUI-subsystem exe
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(versionString())
		return
	}
	if *restore {
		taskbar.New(slog.New(slog.NewTextHandler(os.Stderr, nil))).Restore()
		winglass.RestoreAll()
		return
	}
	if *quit {
		if err := signalQuit(); err != nil {
			fatal("TaskGlass is not running (%v)", err)
		}
		fmt.Println("asked TaskGlass to quit")
		return
	}

	cfgPath := *cfgFlag
	if cfgPath == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			fatal("cannot locate config dir: %v", err)
		}
		cfgPath = filepath.Join(dir, "TaskGlass", "config.json")
	}
	if *reset {
		if err := config.Save(cfgPath, config.Default()); err != nil {
			fatal("reset: %v", err)
		}
		fmt.Println("config reset:", cfgPath)
		return
	}
	if *validate {
		if _, err := os.Stat(cfgPath); err != nil {
			fatal("%v", err)
		}
		if _, err := config.Load(cfgPath); err != nil {
			fatal("%s: %v", cfgPath, err)
		}
		fmt.Println("config OK:", cfgPath)
		return
	}

	// Single instance.
	if _, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName)); err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			return
		}
		fatal("mutex: %v", err)
	}

	cfg, cfgErr := config.Load(cfgPath)
	if cfgErr != nil {
		cfg = config.Default() // keep running; the error is logged below
	}

	log, logPath, closeLog := newLogger(cfg.LogLevel, *verbose)
	defer closeLog()
	if cfgErr != nil {
		log.Error("config invalid, using defaults until it is fixed", "path", cfgPath, "err", cfgErr)
	}

	a := app.New(app.Options{Log: log, CfgPath: cfgPath, LogPath: logPath, Version: version, Config: cfg})
	watchForQuit(a)
	a.Run()
}

// watchForQuit exits cleanly on Ctrl+C / termination signals and when
// `taskglass --quit` is run from another process.
func watchForQuit(a *app.App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		a.Quit()
	}()

	ev, err := windows.CreateEvent(nil, 1, 0, windows.StringToUTF16Ptr(quitName))
	if err != nil {
		return
	}
	go func() {
		for {
			r, _ := windows.WaitForSingleObject(ev, windows.INFINITE)
			if r == windows.WAIT_OBJECT_0 {
				a.Quit()
				return
			}
		}
	}()
}

func signalQuit() error {
	ev, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, windows.StringToUTF16Ptr(quitName))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(ev)
	return windows.SetEvent(ev)
}

func versionString() string {
	rev := ""
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				rev = " (" + s.Value[:7] + ")"
			}
		}
	}
	return fmt.Sprintf("taskglass %s%s %s %s/%s", version, rev, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func newLogger(level string, verbose bool) (*slog.Logger, string, func()) {
	var lv slog.Level
	_ = lv.UnmarshalText([]byte(strings.ToUpper(level)))
	opts := &slog.HandlerOptions{Level: lv}

	dir, err := os.UserCacheDir()
	if err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, opts)), "", func() {}
	}
	dir = filepath.Join(dir, "TaskGlass")
	_ = os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, "taskglass.log")
	// Keep the log bounded: start fresh once it passes 1 MiB.
	if st, err := os.Stat(path); err == nil && st.Size() > 1<<20 {
		_ = os.Remove(path + ".old")
		_ = os.Rename(path, path+".old")
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, opts)), "", func() {}
	}
	var h slog.Handler = slog.NewTextHandler(f, opts)
	if verbose {
		h = multi{h, slog.NewTextHandler(os.Stderr, opts)}
	}
	return slog.New(h), path, func() { f.Close() }
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "taskglass: "+format+"\n", args...)
	os.Exit(1)
}
