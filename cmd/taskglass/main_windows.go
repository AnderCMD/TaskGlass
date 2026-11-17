//go:build windows

// Command taskglass makes the Windows taskbar transparent, blurred or acrylic.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/AnderCMD/TaskGlass/internal/app"
	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/taskbar"
)

// version is overridden at build time: -ldflags "-X main.version=1.2.3".
var version = "dev"

func main() {
	var (
		showVersion = flag.Bool("version", false, "print version and exit")
		cfgFlag     = flag.String("config", "", "path to config.json (default %APPDATA%\\TaskGlass\\config.json)")
		reset       = flag.Bool("reset", false, "rewrite the config file with defaults and exit")
		verbose     = flag.Bool("verbose", false, "log to stderr as well as the log file")
		restore     = flag.Bool("restore", false, "put the taskbar back to the stock look and exit (use after a crash)")
	)
	flag.Parse()

	if *restore {
		taskbar.New(slog.New(slog.NewTextHandler(os.Stderr, nil))).Restore()
		return
	}

	if *showVersion {
		fmt.Println("taskglass", version)
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

	// Single instance.
	if _, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(`Local\TaskGlass.Singleton`)); err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			return
		}
		fatal("mutex: %v", err)
	}

	cfg, cfgErr := config.Load(cfgPath)
	if cfgErr != nil {
		cfg = config.Default() // keep running; the error is logged below
	}

	log, closeLog := newLogger(cfg.LogLevel, *verbose)
	defer closeLog()
	if cfgErr != nil {
		log.Error("config invalid, using defaults until it is fixed", "path", cfgPath, "err", cfgErr)
	}

	a := app.New(log, cfgPath, version, cfg)
	a.Run()
}

func newLogger(level string, verbose bool) (*slog.Logger, func()) {
	var lv slog.Level
	_ = lv.UnmarshalText([]byte(strings.ToUpper(level)))
	opts := &slog.HandlerOptions{Level: lv}

	dir, err := os.UserCacheDir()
	if err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, opts)), func() {}
	}
	dir = filepath.Join(dir, "TaskGlass")
	_ = os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, "taskglass.log")
	// Keep the log bounded: start fresh once it passes 1 MiB.
	if st, err := os.Stat(path); err == nil && st.Size() > 1<<20 {
		_ = os.Rename(path, path+".old")
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return slog.New(slog.NewTextHandler(os.Stderr, opts)), func() {}
	}
	var h slog.Handler = slog.NewTextHandler(f, opts)
	if verbose {
		h = multi{h, slog.NewTextHandler(os.Stderr, opts)}
	}
	return slog.New(h), func() { f.Close() }
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "taskglass: "+format+"\n", args...)
	os.Exit(1)
}
