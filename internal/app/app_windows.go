//go:build windows

// Package app wires the config, taskbar manager, window events and tray menu
// together.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/systray"

	"github.com/AnderCMD/TaskGlass/assets"
	"github.com/AnderCMD/TaskGlass/internal/autostart"
	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/events"
	"github.com/AnderCMD/TaskGlass/internal/i18n"
	"github.com/AnderCMD/TaskGlass/internal/settings"
	"github.com/AnderCMD/TaskGlass/internal/taskbar"
	"github.com/AnderCMD/TaskGlass/internal/theme"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
	"github.com/AnderCMD/TaskGlass/internal/xaml"
)

// Options configures an App.
type Options struct {
	Log     *slog.Logger
	CfgPath string
	LogPath string
	Version string
	Config  config.Config
	// ShowSettings opens the settings window right after start-up.
	ShowSettings bool
}

// App is the running TaskGlass instance.
type App struct {
	Options
	exe string

	mu      sync.RWMutex
	cfg     config.Config
	modTime time.Time

	systemLight atomic.Bool  // Windows reports a light theme
	menuLight   atomic.Int32 // 0 not applied yet, 1 dark, 2 light
	mgr         *taskbar.Manager
	xaml        *xaml.Controller // nil when the TAP DLL is unavailable
	ev          *events.Watcher
	kick        chan struct{} // asks the loop for an immediate update

	srv      *settings.Server
	menu     menuItems
	cancel   context.CancelFunc
	done     chan struct{}
	quitReq  atomic.Bool
	stopOnce sync.Once
}

type menuItems struct {
	title, enabled, hotkey, maximized, xamlFix, autostart *systray.MenuItem
	settings, styleMenu, themeMenu, langMenu              *systray.MenuItem
	surfTray, openCfg, openLog, reload, quit              *systray.MenuItem
	accents                                               map[config.Accent]*systray.MenuItem
	themes                                                map[config.Theme]*systray.MenuItem
	langs                                                 map[i18n.Lang]*systray.MenuItem
}

// New creates an App with an already-loaded configuration.
func New(o Options) *App {
	exe, _ := os.Executable()
	a := &App{Options: o, exe: exe, cfg: o.Config, mgr: taskbar.New(o.Log), kick: make(chan struct{}, 1)}
	if st, err := os.Stat(o.CfgPath); err == nil {
		a.modTime = st.ModTime()
	}
	a.systemLight.Store(theme.SystemIsLight())
	a.srv = settings.New(host{a})
	if autostart.Repair(exe) {
		o.Log.Info("autostart entry repaired")
	}
	if p, ok := xaml.Available(exe); ok {
		a.xaml = xaml.New(o.Log, p)
	} else {
		o.Log.Warn("taskglass_tap.dll not found next to the executable: on Windows 11 the taskbar cannot become fully transparent", "expected", p)
	}
	return a
}

// Run blocks until the user quits (tray menu, --quit or a signal).
func (a *App) Run() {
	a.applyMenuTheme(a.snapshot())
	systray.Run(a.onReady, a.onExit)
	a.shutdown() // Run can return before onExit has finished; wait for the cleanup
}

// Quit asks the app to exit. It is safe to call from any goroutine.
func (a *App) Quit() {
	if a.quitReq.CompareAndSwap(false, true) {
		systray.Quit()
	}
}

func (a *App) snapshot() config.Config {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cfg
}

func (a *App) wake() {
	select {
	case a.kick <- struct{}{}:
	default:
	}
}

// update applies fn to a copy of the config, saves it and publishes it.
func (a *App) update(fn func(*config.Config)) {
	if err := a.apply(fn); err != nil {
		a.Log.Error("save config", "err", err)
	}
}

// apply is update that reports why a change was rejected (invalid values or a
// file that cannot be written).
func (a *App) apply(fn func(*config.Config)) error {
	a.mu.Lock()
	next := a.cfg
	next.Apps = append([]config.AppRule{}, a.cfg.Apps...)
	fn(&next)
	if err := config.Save(a.CfgPath, next); err != nil {
		a.mu.Unlock()
		return err
	}
	a.cfg = next
	if st, err := os.Stat(a.CfgPath); err == nil {
		a.modTime = st.ModTime()
	}
	a.mu.Unlock()
	a.applySideEffects(next)
	a.syncMenu()
	a.wake()
	return nil
}

// lang is the language the UI currently uses.
func (a *App) lang() i18n.Lang { return i18n.Resolve(a.snapshot().Language, i18n.SystemLocale()) }

// OpenSettings shows the settings window (starting its local server on first use).
func (a *App) OpenSettings() {
	if a.srv == nil {
		return
	}
	url, err := a.srv.URL()
	if err != nil {
		a.Log.Error("settings server", "err", err)
		return
	}
	profile := filepath.Join(os.TempDir(), "TaskGlass-ui")
	if dir, err := os.UserCacheDir(); err == nil {
		profile = filepath.Join(dir, "TaskGlass", "ui-profile")
	}
	if err := settings.Show(url, profile); err != nil {
		a.Log.Error("open settings window", "err", err)
	}
}

// host adapts App to settings.Host.
type host struct{ a *App }

func (h host) Config() config.Config { return h.a.snapshot() }
func (h host) Apply(c config.Config) error {
	next := c
	next.Apps = append([]config.AppRule{}, c.Apps...)
	return h.a.apply(func(cur *config.Config) { *cur = next })
}
func (h host) AutostartEnabled() bool { return autostart.Enabled(h.a.exe) }
func (h host) SetAutostart(on bool) error {
	if err := autostart.Set(h.a.exe, on); err != nil {
		return err
	}
	h.a.syncMenu()
	return nil
}
func (h host) OpenFile(which string) error {
	path := h.a.CfgPath
	if which == "log" {
		path = h.a.LogPath
	}
	return exec.Command("notepad.exe", path).Start()
}
func (h host) Version() string        { return h.a.Version }
func (h host) ConfigPath() string     { return h.a.CfgPath }
func (h host) XamlAvailable() bool    { return h.a.xaml != nil }
func (h host) BlurNeedsAcrylic() bool { return winapi.BlurNeedsAcrylic() }

// reload re-reads the config file if it changed on disk.
func (a *App) reload() {
	st, err := os.Stat(a.CfgPath)
	if err != nil {
		return
	}
	a.mu.RLock()
	same := st.ModTime().Equal(a.modTime)
	a.mu.RUnlock()
	if same {
		return
	}
	c, err := config.Load(a.CfgPath)
	a.mu.Lock()
	a.modTime = st.ModTime() // remember even on error so we log once per edit
	if err != nil {
		a.mu.Unlock()
		a.Log.Error("reload config: keeping previous settings", "err", err)
		return
	}
	a.cfg = c
	a.mu.Unlock()
	a.Log.Info("config reloaded")
	a.applySideEffects(c)
	a.syncMenu()
	a.wake()
}

// applySideEffects handles settings that are more than "what to paint".
func (a *App) applySideEffects(c config.Config) {
	if a.ev != nil {
		// Only grab the shortcut when the user turned it on.
		if !c.HotkeyEnabled || c.Hotkey == "" {
			a.ev.SetHotkey(0, 0)
		} else if mods, vk, err := config.ParseHotkey(c.Hotkey); err == nil {
			a.ev.SetHotkey(mods, vk)
		}
	}
	a.updateTrayIcon(c)
}

// applyMenuTheme gives the tray icon's right-click menu the light or dark look
// of Windows (apps mode), or the forced one when the config pins a theme.
func (a *App) applyMenuTheme(c config.Config) {
	light := c.Theme.Light(theme.AppsAreLight())
	v := int32(0)
	if light {
		v = 1
	}
	if a.menuLight.Swap(v+1) != v+1 {
		theme.ApplyMenuTheme(light)
	}
}

func (a *App) isLight(c config.Config) bool { return c.Theme.Light(a.systemLight.Load()) }

func (a *App) updateTrayIcon(c config.Config) {
	var icon []byte
	switch c.TrayIcon {
	case config.TrayColor:
		icon = assets.Color
	case config.TrayLight:
		icon = assets.GlyphLight
	case config.TrayDark:
		icon = assets.GlyphDark
	default: // auto: contrast with the taskbar
		if a.isLight(c) {
			icon = assets.GlyphDark
		} else {
			icon = assets.GlyphLight
		}
	}
	systray.SetIcon(icon)
}

// loop is the heart of the app: it updates the taskbars whenever a window
// event, hotkey, config change or theme change asks for it, and at least every
// poll_ms as a safety net.
func (a *App) loop(ctx context.Context) {
	defer close(a.done)
	defer a.mgr.Restore()
	if a.xaml != nil {
		if err := a.xaml.Start(); err != nil {
			a.Log.Warn("xaml fix disabled", "err", err)
			a.xaml = nil
		} else {
			defer a.xaml.Stop() // runs before Restore: show the XAML fill, then reset the accent
		}
	}

	housekeeping := time.NewTicker(time.Second)
	defer housekeeping.Stop()
	timer := time.NewTimer(0)
	defer timer.Stop()
	panics := newPanicGuard(5, time.Minute)

	for {
		c := a.snapshot()
		if !a.safely(panics, func() {
			styled := a.mgr.Tick(c, a.isLight(c))
			if a.xaml != nil {
				on := c.Enabled && c.XamlFix
				a.xaml.SetHide(on && styled)
				a.xaml.SetOverflow(xaml.Encode(on && c.Surfaces.TrayOverflow.Enabled, int(c.Surfaces.TrayOverflow.Opacity)))
			}
		}) {
			a.Log.Error("too many internal errors; restoring the taskbar and exiting")
			a.Quit()
			<-ctx.Done()
			return
		}

		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(time.Duration(c.PollMs) * time.Millisecond)

		var changed, hotkey <-chan struct{}
		if a.ev != nil {
			changed, hotkey = a.ev.Changed, a.ev.Hotkey
		}
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		case <-a.kick:
		case <-hotkey:
			a.update(func(c *config.Config) { c.Enabled = !c.Enabled })
			a.Log.Info("hotkey toggled", "enabled", a.snapshot().Enabled)
		case <-changed:
			// Coalesce a burst of window events into one update.
			if d := time.Duration(c.DebounceMs) * time.Millisecond; d > 0 {
				time.Sleep(d)
			}
			select {
			case <-a.ev.Changed:
			default:
			}
		case <-housekeeping.C:
			a.reload()
			a.applyMenuTheme(a.snapshot())
			if l := theme.SystemIsLight(); l != a.systemLight.Swap(l) {
				a.Log.Info("windows theme changed", "light", l)
				a.updateTrayIcon(a.snapshot())
			}
		}
	}
}

// safely runs fn, converting a panic into a logged error. It reports false
// when panics are happening too often to keep going.
func (a *App) safely(g *panicGuard, fn func()) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			a.Log.Error("panic recovered", "panic", fmt.Sprint(r), "stack", string(debug.Stack()))
			a.mgr.Restore()
			ok = g.allow()
		}
	}()
	fn()
	return true
}

func (a *App) onReady() {
	systray.SetTitle("TaskGlass")
	systray.SetTooltip("TaskGlass " + a.Version)
	a.updateTrayIcon(a.snapshot())

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.done = make(chan struct{})

	if ev, err := events.Start(a.Log); err != nil {
		a.Log.Warn("window events unavailable, falling back to polling", "err", err)
	} else {
		a.ev = ev
		a.applySideEffects(a.snapshot())
	}

	a.buildMenu(ctx)
	go a.loop(ctx)
	a.Log.Info("started", "version", a.Version, "config", a.CfgPath)
	if a.Options.ShowSettings {
		go a.OpenSettings()
	}
}

func (a *App) onExit() { a.shutdown() }

// shutdown stops the loop, which restores the stock taskbar, and releases the
// event hooks. It is idempotent and blocks until the work is done: the tray
// library may return from Run before (or while) it calls onExit, and the
// process must not exit before the taskbar has been restored.
func (a *App) shutdown() {
	a.stopOnce.Do(func() {
		if a.srv != nil {
			a.srv.Close()
		}
		if a.cancel != nil {
			a.cancel()
			select {
			case <-a.done: // the loop restores the stock taskbar before closing done
			case <-time.After(10 * time.Second):
				a.Log.Error("shutdown timed out")
			}
		}
		if a.ev != nil {
			a.ev.Stop()
		}
		a.Log.Info("stopped")
	})
}

// ---- tray menu ---------------------------------------------------------

func check(it *systray.MenuItem, on bool) {
	if on {
		it.Check()
	} else {
		it.Uncheck()
	}
}

// syncMenu refreshes check marks and (re)translates every menu entry, so a
// language change takes effect immediately.
func (a *App) syncMenu() {
	m := a.menu
	if m.enabled == nil {
		return
	}
	c := a.snapshot()
	l := a.lang()
	t := func(k string) string { return i18n.T(l, k) }
	set := func(it *systray.MenuItem, title, tip string) {
		it.SetTitle(title)
		if tip != "" {
			it.SetTooltip(tip)
		}
	}
	check(m.enabled, c.Enabled)
	check(m.maximized, c.Maximized.Enabled)
	check(m.xamlFix, c.XamlFix)
	check(m.surfTray, c.Surfaces.TrayOverflow.Enabled)
	check(m.hotkey, c.HotkeyEnabled)
	check(m.autostart, autostart.Enabled(a.exe))
	for k, it := range m.accents {
		check(it, k == c.Default.Accent)
		it.SetTitle(t("accent." + string(k)))
	}
	for k, it := range m.themes {
		check(it, k == c.Theme)
		it.SetTitle(t("theme." + string(k)))
	}
	for k, it := range m.langs {
		check(it, k == c.Language)
		it.SetTitle(t("language." + string(k)))
	}
	title := t("tray.enabled")
	if c.HotkeyEnabled && c.Hotkey != "" {
		title += "\t" + c.Hotkey
	}
	set(m.enabled, title, t("tray.enabled.tip"))
	set(m.settings, t("tray.settings"), t("tray.settings.tip"))
	set(m.styleMenu, t("tray.style"), "")
	set(m.themeMenu, t("tray.theme"), "")
	set(m.langMenu, t("tray.language"), "")
	set(m.maximized, t("tray.maximized"), "")
	set(m.hotkey, fmt.Sprintf(t("tray.hotkey"), c.Hotkey), t("tray.hotkey.tip"))
	set(m.surfTray, t("tray.flyout"), t("tray.flyout.tip"))
	set(m.xamlFix, t("tray.xaml"), t("tray.xaml.tip"))
	set(m.autostart, t("tray.autostart"), "")
	set(m.openCfg, t("tray.openconfig"), "")
	set(m.openLog, t("tray.openlog"), "")
	set(m.reload, t("tray.reload"), "")
	set(m.quit, t("tray.quit"), "")
}

func (a *App) buildMenu(ctx context.Context) {
	c := a.snapshot()
	m := &a.menu
	m.title = systray.AddMenuItem("TaskGlass "+a.Version, "")
	m.title.Disable()
	systray.AddSeparator()

	m.settings = systray.AddMenuItem("Settings", "")
	m.enabled = systray.AddMenuItemCheckbox("Enabled", "", c.Enabled)
	systray.AddSeparator()

	m.styleMenu = systray.AddMenuItem("Style", "")
	m.accents = map[config.Accent]*systray.MenuItem{}
	for _, ac := range config.Accents {
		m.accents[ac] = m.styleMenu.AddSubMenuItemCheckbox(string(ac), "", ac == c.Default.Accent)
	}
	m.themeMenu = systray.AddMenuItem("Theme", "")
	m.themes = map[config.Theme]*systray.MenuItem{}
	for _, th := range config.Themes {
		m.themes[th] = m.themeMenu.AddSubMenuItemCheckbox(string(th), "", th == c.Theme)
	}
	m.langMenu = systray.AddMenuItem("Language", "")
	m.langs = map[i18n.Lang]*systray.MenuItem{}
	for _, lg := range []i18n.Lang{i18n.Auto, i18n.English, i18n.Spanish} {
		m.langs[lg] = m.langMenu.AddSubMenuItemCheckbox(string(lg), "", lg == c.Language)
	}

	m.maximized = systray.AddMenuItemCheckbox("Maximized", "", c.Maximized.Enabled)
	m.hotkey = systray.AddMenuItemCheckbox("Hotkey", "", c.HotkeyEnabled)
	m.surfTray = systray.AddMenuItemCheckbox("Flyout", "", c.Surfaces.TrayOverflow.Enabled)
	m.xamlFix = systray.AddMenuItemCheckbox("Transparent", "", c.XamlFix)
	if a.xaml == nil {
		m.xamlFix.Disable()
	}
	systray.AddSeparator()
	m.autostart = systray.AddMenuItemCheckbox("Autostart", "", autostart.Enabled(a.exe))
	m.openCfg = systray.AddMenuItem("Config", "")
	m.openLog = systray.AddMenuItem("Log", "")
	m.reload = systray.AddMenuItem("Reload", "")
	systray.AddSeparator()
	m.quit = systray.AddMenuItem("Quit", "")

	// A left click on the tray icon opens the settings window; the menu stays
	// on the right click.
	systray.SetOnTapped(a.OpenSettings)

	spawn := func(ch <-chan struct{}, fn func()) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					fn()
				}
			}
		}()
	}
	flip := func(it *systray.MenuItem, set func(c *config.Config, on bool)) func() {
		return func() {
			on := !it.Checked()
			a.update(func(c *config.Config) { set(c, on) })
		}
	}
	spawn(m.settings.ClickedCh, a.OpenSettings)
	spawn(m.enabled.ClickedCh, flip(m.enabled, func(c *config.Config, on bool) { c.Enabled = on }))
	spawn(m.hotkey.ClickedCh, flip(m.hotkey, func(c *config.Config, on bool) { c.HotkeyEnabled = on }))
	spawn(m.maximized.ClickedCh, flip(m.maximized, func(c *config.Config, on bool) { c.Maximized.Enabled = on }))
	spawn(m.surfTray.ClickedCh, flip(m.surfTray, func(c *config.Config, on bool) { c.Surfaces.TrayOverflow.Enabled = on }))
	spawn(m.xamlFix.ClickedCh, flip(m.xamlFix, func(c *config.Config, on bool) { c.XamlFix = on }))
	for ac, it := range m.accents {
		ac := ac
		spawn(it.ClickedCh, func() { a.update(func(c *config.Config) { c.Default.Accent = ac }) })
	}
	for th, it := range m.themes {
		th := th
		spawn(it.ClickedCh, func() { a.update(func(c *config.Config) { c.Theme = th }) })
	}
	for lg, it := range m.langs {
		lg := lg
		spawn(it.ClickedCh, func() { a.update(func(c *config.Config) { c.Language = lg }) })
	}
	spawn(m.autostart.ClickedCh, func() {
		on := !m.autostart.Checked()
		if err := autostart.Set(a.exe, on); err != nil {
			a.Log.Error("autostart", "err", err)
			return
		}
		check(m.autostart, on)
	})
	// notepad is always present and never interprets the file it opens.
	open := func(path string) func() {
		return func() {
			if err := exec.Command("notepad.exe", path).Start(); err != nil {
				a.Log.Error("open file", "path", filepath.Base(path), "err", err)
			}
		}
	}
	spawn(m.openCfg.ClickedCh, open(a.CfgPath))
	spawn(m.openLog.ClickedCh, open(a.LogPath))
	spawn(m.reload.ClickedCh, func() {
		a.mu.Lock()
		a.modTime = time.Time{} // force a re-read
		a.mu.Unlock()
		a.reload()
	})
	go func() {
		select {
		case <-ctx.Done():
		case <-m.quit.ClickedCh:
			a.Quit()
		}
	}()
	a.syncMenu()
}

// ---- panic guard -------------------------------------------------------

// panicGuard allows up to max panics inside a sliding window.
type panicGuard struct {
	max    int
	window time.Duration
	times  []time.Time
}

func newPanicGuard(max int, window time.Duration) *panicGuard {
	return &panicGuard{max: max, window: window}
}

// allow records a panic and reports whether the app may keep running.
func (g *panicGuard) allow() bool {
	now := time.Now()
	keep := g.times[:0]
	for _, t := range g.times {
		if now.Sub(t) < g.window {
			keep = append(keep, t)
		}
	}
	g.times = append(keep, now)
	return len(g.times) <= g.max
}
