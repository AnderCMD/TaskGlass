//go:build windows

// Package app wires the config, taskbar manager and tray menu together.
package app

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"fyne.io/systray"

	"github.com/AnderCMD/TaskGlass/assets"
	"github.com/AnderCMD/TaskGlass/internal/autostart"
	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/taskbar"
	"github.com/AnderCMD/TaskGlass/internal/xaml"
)

// App is the running TaskGlass instance.
type App struct {
	log     *slog.Logger
	cfgPath string
	exe     string
	version string

	mu      sync.RWMutex
	cfg     config.Config
	modTime time.Time

	mgr    *taskbar.Manager
	xaml   *xaml.Controller // nil when the TAP DLL is unavailable
	cancel context.CancelFunc
	done   chan struct{}
}

// New creates an App with an already-loaded configuration.
func New(log *slog.Logger, cfgPath, version string, cfg config.Config) *App {
	exe, _ := os.Executable()
	a := &App{log: log, cfgPath: cfgPath, exe: exe, version: version, cfg: cfg, mgr: taskbar.New(log)}
	if st, err := os.Stat(cfgPath); err == nil {
		a.modTime = st.ModTime()
	}
	if p, ok := xaml.Available(exe); ok {
		a.xaml = xaml.New(log, p)
	} else {
		log.Warn("taskglass_tap.dll not found next to the executable: on Windows 11 the taskbar cannot become fully transparent", "expected", p)
	}
	return a
}

// Run blocks until the user quits from the tray menu.
func (a *App) Run() {
	systray.Run(a.onReady, a.onExit)
}

func (a *App) snapshot() config.Config {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cfg
}

func (a *App) update(fn func(*config.Config)) {
	a.mu.Lock()
	next := a.cfg
	fn(&next)
	if err := config.Save(a.cfgPath, next); err != nil {
		a.mu.Unlock()
		a.log.Error("save config", "err", err)
		return
	}
	a.cfg = next
	if st, err := os.Stat(a.cfgPath); err == nil {
		a.modTime = st.ModTime()
	}
	a.mu.Unlock()
}

func (a *App) reload() {
	st, err := os.Stat(a.cfgPath)
	if err != nil {
		return
	}
	a.mu.RLock()
	same := st.ModTime().Equal(a.modTime)
	a.mu.RUnlock()
	if same {
		return
	}
	c, err := config.Load(a.cfgPath)
	a.mu.Lock()
	a.modTime = st.ModTime() // remember even on error so we log once per edit
	if err != nil {
		a.mu.Unlock()
		a.log.Error("reload config: keeping previous settings", "err", err)
		return
	}
	a.cfg = c
	a.mu.Unlock()
	a.log.Info("config reloaded")
}

func (a *App) loop(ctx context.Context) {
	defer close(a.done)
	defer a.mgr.Restore()
	if a.xaml != nil {
		if err := a.xaml.Start(); err != nil {
			a.log.Warn("xaml fix disabled", "err", err)
			a.xaml = nil
		} else {
			defer a.xaml.Stop() // runs before Restore: show the XAML fill, then reset the accent
		}
	}
	var tick *time.Ticker
	var interval time.Duration
	defer func() {
		if tick != nil {
			tick.Stop()
		}
	}()
	lastReload := time.Now()
	for {
		c := a.snapshot()
		if d := time.Duration(c.PollMs) * time.Millisecond; d != interval {
			if tick != nil {
				tick.Stop()
			}
			interval, tick = d, time.NewTicker(d)
		}
		styled := a.mgr.Tick(c)
		if a.xaml != nil {
			a.xaml.SetHide(c.XamlFix && styled)
		}
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		if time.Since(lastReload) > time.Second {
			a.reload()
			lastReload = time.Now()
		}
	}
}

func (a *App) onReady() {
	systray.SetIcon(assets.Icon)
	systray.SetTitle("TaskGlass")
	systray.SetTooltip("TaskGlass " + a.version)

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.done = make(chan struct{})
	go a.loop(ctx)

	a.buildMenu(ctx)
	a.log.Info("started", "version", a.version, "config", a.cfgPath)
}

func (a *App) buildMenu(ctx context.Context) {
	cur := a.snapshot().Default.Accent
	title := systray.AddMenuItem("TaskGlass "+a.version, "")
	title.Disable()
	systray.AddSeparator()

	accentItems := map[config.Accent]*systray.MenuItem{}
	labels := map[config.Accent]string{
		config.AccentClear:   "Clear",
		config.AccentBlur:    "Blur",
		config.AccentAcrylic: "Acrylic",
		config.AccentOpaque:  "Opaque",
		config.AccentNormal:  "Normal (Windows default)",
	}
	for _, ac := range config.Accents {
		it := systray.AddMenuItemCheckbox(labels[ac], "Default taskbar style", ac == cur)
		accentItems[ac] = it
		go func(ac config.Accent, it *systray.MenuItem) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-it.ClickedCh:
					a.update(func(c *config.Config) { c.Default.Accent = ac })
					for k, other := range accentItems {
						if k == ac {
							other.Check()
						} else {
							other.Uncheck()
						}
					}
				}
			}
		}(ac, it)
	}

	systray.AddSeparator()
	c := a.snapshot()
	maxItem := systray.AddMenuItemCheckbox("Opaque when a window is maximized", "", c.Maximized.Enabled)
	startItem := systray.AddMenuItemCheckbox("Styled Start menu", "Use the start_menu rule while Start is open", c.StartMenu.Enabled)
	autoItem := systray.AddMenuItemCheckbox("Start with Windows", "", autostart.Enabled(a.exe))
	systray.AddSeparator()
	openCfg := systray.AddMenuItem("Open config file", "")
	quit := systray.AddMenuItem("Quit TaskGlass", "")

	toggle := func(it *systray.MenuItem, set func(c *config.Config, on bool)) {
		on := !it.Checked()
		a.update(func(c *config.Config) { set(c, on) })
		if on {
			it.Check()
		} else {
			it.Uncheck()
		}
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-maxItem.ClickedCh:
				toggle(maxItem, func(c *config.Config, on bool) { c.Maximized.Enabled = on })
			case <-startItem.ClickedCh:
				toggle(startItem, func(c *config.Config, on bool) { c.StartMenu.Enabled = on })
			case <-autoItem.ClickedCh:
				on := !autoItem.Checked()
				if err := autostart.Set(a.exe, on); err != nil {
					a.log.Error("autostart", "err", err)
					continue
				}
				if on {
					autoItem.Check()
				} else {
					autoItem.Uncheck()
				}
			case <-openCfg.ClickedCh:
				// notepad is always present and never interprets the file.
				if err := exec.Command("notepad.exe", a.cfgPath).Start(); err != nil {
					a.log.Error("open config", "err", err)
				}
			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (a *App) onExit() {
	if a.cancel != nil {
		a.cancel()
		<-a.done // loop restores the stock taskbar before we return
	}
	a.log.Info("stopped")
}
