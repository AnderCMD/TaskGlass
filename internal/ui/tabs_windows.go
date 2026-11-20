//go:build windows

package ui

import (
	"fmt"
	"os/exec"
	"time"
	"unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/i18n"
	"github.com/AnderCMD/TaskGlass/internal/preview"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

var (
	languageChoices = []choice{{"auto", "language.auto"}, {"en", "language.en"}, {"es", "language.es"}}
	themeChoices    = []choice{{"auto", "theme.auto"}, {"light", "theme.light"}, {"dark", "theme.dark"}}
	taskbarChoices  = []choice{{"all", "general.taskbars.all"}, {"primary", "general.taskbars.primary"}}
	trayChoices     = []choice{{"auto", "general.trayIcon.auto"}, {"color", "general.trayIcon.color"}, {"light", "general.trayIcon.light"}, {"dark", "general.trayIcon.dark"}}
	updateChoices   = []choice{{"notify", "update.mode.notify"}, {"auto", "update.mode.auto"}, {"off", "update.mode.off"}}
	logChoices      = []choice{{"debug", ""}, {"info", ""}, {"warn", ""}, {"error", ""}}
)

func unsafeSizeofChooseColor() uintptr { return unsafe.Sizeof(win.CHOOSECOLOR{}) }

// page adds a scrolling tab page and returns the container to fill.
func (w *window) page(key string) walk.Container {
	tp, _ := walk.NewTabPage()
	tp.SetLayout(walk.NewVBoxLayout())
	if l, ok := tp.Layout().(*walk.BoxLayout); ok {
		l.SetMargins(walk.Margins{})
	}
	w.text(func(s string) { tp.SetTitle(s) }, key)
	_ = w.tabs.Pages().Add(tp)
	body, _ := walk.NewComposite(tp)
	l := walk.NewVBoxLayout()
	l.SetMargins(walk.Margins{HNear: 14, VNear: 12, HFar: 14, VFar: 12})
	l.SetSpacing(6)
	body.SetLayout(l)
	return body
}

func (w *window) buildTabs(parent walk.Container) {
	w.tabs, _ = walk.NewTabWidget(parent)
	w.tabs.SetAlignment(walk.AlignHFarVNear)

	// --- General
	g := w.page("ui.tab.general")
	w.enabled = w.check(g, "general.enabled", func() bool { return w.cfg.Enabled }, func(v bool) { w.cfg.Enabled = v })
	w.help(g, "general.enabled.help")
	w.autostart = w.check(g, "general.autostart", func() bool { return w.auto }, func(v bool) { w.auto = v })
	w.help(g, "general.autostart.help")
	w.caption(g, "general.language")
	w.language = w.combo(g, languageChoices, func() string { return string(w.cfg.Language) }, func(id string) { w.cfg.Language = i18n.Lang(id) })
	w.caption(g, "general.theme")
	w.themeC = w.combo(g, themeChoices, func() string { return string(w.cfg.Theme) }, func(id string) { w.cfg.Theme = config.Theme(id) })
	w.help(g, "general.theme.help")
	w.caption(g, "general.taskbars")
	w.taskbars = w.combo(g, taskbarChoices, func() string { return string(w.cfg.Taskbars) }, func(id string) { w.cfg.Taskbars = config.Taskbars(id) })
	w.caption(g, "general.trayIcon")
	w.trayIcon = w.combo(g, trayChoices, func() string { return string(w.cfg.TrayIcon) }, func(id string) { w.cfg.TrayIcon = config.TrayIcon(id) })
	w.caption(g, "general.updates")
	w.updMode = w.combo(g, updateChoices, func() string { return string(w.cfg.Updates) }, func(id string) { w.cfg.Updates = config.UpdateMode(id) })
	w.help(g, "update.help")
	ur := hbox(g, 8)
	w.updCheck = w.button(ur, "update.check", w.checkUpdates)
	w.updInstall = w.button(ur, "update.install", w.installUpdate)
	w.updInstall.SetVisible(false)
	w.updMsg, _ = walk.NewLabel(ur)
	walk.NewHSpacer(ur)
	w.caption(g, "general.hotkey")
	w.hotkeyOn = w.check(g, "general.hotkey.enable", func() bool { return w.cfg.HotkeyEnabled }, func(v bool) { w.cfg.HotkeyEnabled = v })
	w.hotkey, _ = walk.NewLineEdit(g)
	w.hotkey.TextChanged().Attach(func() {
		if w.syncing {
			return
		}
		w.cfg.Hotkey = w.hotkey.Text()
		w.changed()
	})
	w.hotkeyWarn, _ = walk.NewLabel(g)
	w.help(g, "general.hotkey.help")
	walk.NewVSpacer(g)

	// --- Appearance
	a := w.page("ui.tab.appearance")
	w.def = w.newAppEditor(a, func() *config.Appearance { return &w.cfg.Default })
	w.caption(a, "appear.preset")
	pr := hbox(a, 6)
	for _, p := range []struct {
		key    string
		accent config.Accent
		op     config.Percent
	}{
		{"appear.preset.glass", config.AccentClear, 0},
		{"appear.preset.frosted", config.AccentBlur, 35},
		{"appear.preset.acrylic", config.AccentAcrylic, 55},
		{"appear.preset.solid", config.AccentOpaque, 100},
	} {
		p := p
		w.button(pr, p.key, func() {
			w.cfg.Default.Accent, w.cfg.Default.Opacity = p.accent, p.op
			w.pull()
			w.changed()
		})
	}
	walk.NewHSpacer(pr)
	w.xaml = w.check(a, "appear.xaml", func() bool { return w.cfg.XamlFix }, func(v bool) { w.cfg.XamlFix = v })
	w.help(a, "appear.xaml.help")
	if !w.host.XamlAvailable() {
		l, _ := walk.NewTextLabel(a)
		w.text(func(s string) { l.SetText(s) }, "appear.xaml.missing")
	}
	w.help(a, "appear.accentNote")
	walk.NewVSpacer(a)

	// --- Maximized window
	m := w.page("ui.tab.maximized")
	w.max = w.check(m, "appear.maxRule", func() bool { return w.cfg.Maximized.Enabled }, func(v bool) { w.cfg.Maximized.Enabled = v })
	w.help(m, "appear.maxRule.help")
	w.maxGroup = vbox(m, walk.Margins{VNear: 6}, 6)
	w.maxEd = w.newAppEditor(w.maxGroup, func() *config.Appearance { return &w.cfg.Maximized.Appearance })
	walk.NewVSpacer(m)

	// --- Hidden icons
	t := w.page("ui.tab.tray")
	w.flyOn = w.check(t, "tray.surface", func() bool { return w.cfg.Surfaces.TrayOverflow.Enabled }, func(v bool) { w.cfg.Surfaces.TrayOverflow.Enabled = v })
	w.help(t, "tray.surface.help")
	w.caption(t, "tray.opacity")
	w.flyOp, w.flyOpVal = w.slider(t, func() int { return int(w.cfg.Surfaces.TrayOverflow.Opacity) }, func(v int) { w.cfg.Surfaces.TrayOverflow.Opacity = config.Percent(v) })
	w.help(t, "tray.opacity.help")
	w.flyMax = w.check(t, "tray.maxRule", func() bool { return w.cfg.Surfaces.TrayOverflow.Maximized.Enabled }, func(v bool) { w.cfg.Surfaces.TrayOverflow.Maximized.Enabled = v })
	w.help(t, "tray.maxRule.help")
	w.caption(t, "tray.maxOpacity")
	w.flyMaxOp, w.flyMaxOpVal = w.slider(t, func() int { return int(w.cfg.Surfaces.TrayOverflow.Maximized.Opacity) }, func(v int) { w.cfg.Surfaces.TrayOverflow.Maximized.Opacity = config.Percent(v) })
	walk.NewVSpacer(t)

	// --- Per app
	p := w.page("ui.tab.apps")
	w.help(p, "apps.help")
	split := hbox(p, 14)
	left := vbox(split, walk.Margins{}, 6)
	w.appList, _ = walk.NewListBox(left)
	w.appList.SetMinMaxSize(walk.Size{Width: 190, Height: 220}, walk.Size{Width: 240})
	w.appList.CurrentIndexChanged().Attach(func() {
		if w.syncing {
			return
		}
		w.pullApp()
	})
	br := hbox(left, 6)
	w.appAdd = w.button(br, "apps.add", func() {
		w.cfg.Apps = append(w.cfg.Apps, config.AppRule{Appearance: config.Appearance{Accent: config.AccentOpaque, Color: config.ColorAuto, Opacity: 100}})
		w.rebuildApps(len(w.cfg.Apps) - 1)
		w.changed()
		w.appExe.SetFocus()
	})
	w.appRemove = w.button(br, "apps.remove", func() {
		if i := w.appList.CurrentIndex(); i >= 0 && i < len(w.cfg.Apps) {
			w.cfg.Apps = append(w.cfg.Apps[:i], w.cfg.Apps[i+1:]...)
			w.rebuildApps(i)
			w.changed()
		}
	})
	w.appLimit, _ = walk.NewLabel(left)
	w.appEditorBox = vbox(split, walk.Margins{}, 6)
	w.appPick, _ = walk.NewLabel(w.appEditorBox)
	w.text(func(s string) { w.appPick.SetText(s) }, "apps.pick")
	w.caption(w.appEditorBox, "apps.exe")
	w.appExe, _ = walk.NewLineEdit(w.appEditorBox)
	w.appExe.TextChanged().Attach(func() {
		if w.syncing {
			return
		}
		if i := w.appList.CurrentIndex(); i >= 0 && i < len(w.cfg.Apps) {
			w.cfg.Apps[i].Exe = w.appExe.Text()
			w.refreshAppNames(i)
			w.changed()
		}
	})
	w.appEd = w.newAppEditor(w.appEditorBox, func() *config.Appearance {
		if i := w.appList.CurrentIndex(); i >= 0 && i < len(w.cfg.Apps) {
			return &w.cfg.Apps[i].Appearance
		}
		return nil
	})
	walk.NewVSpacer(w.appEditorBox)
	split.Layout().(*walk.BoxLayout).SetStretchFactor(w.appEditorBox, 2)

	// --- Advanced
	d := w.page("ui.tab.advanced")
	w.help(d, "adv.help")
	num := func(key string, min, max float64, get func() int, set func(int)) *walk.NumberEdit {
		w.caption(d, key)
		n, _ := walk.NewNumberEdit(d)
		_ = n.SetRange(min, max)
		_ = n.SetDecimals(0)
		n.ValueChanged().Attach(func() {
			if w.syncing {
				return
			}
			set(int(n.Value()))
			w.changed()
		})
		return n
	}
	w.poll = num("adv.poll", 50, 10000, func() int { return w.cfg.PollMs }, func(v int) { w.cfg.PollMs = v })
	w.debounce = num("adv.debounce", 0, 1000, func() int { return w.cfg.DebounceMs }, func(v int) { w.cfg.DebounceMs = v })
	w.refresh = num("adv.refresh", 100, 60000, func() int { return w.cfg.RefreshMs }, func(v int) { w.cfg.RefreshMs = v })
	w.caption(d, "adv.log")
	w.logLevel = w.combo(d, logChoices, func() string { return w.cfg.LogLevel }, func(id string) { w.cfg.LogLevel = id })
	ob := hbox(d, 8)
	w.button(ob, "adv.openConfig", func() { _ = w.host.OpenFile("config") })
	w.button(ob, "adv.openLog", func() { _ = w.host.OpenFile("log") })
	walk.NewHSpacer(ob)
	walk.NewVSpacer(d)

	// --- About
	ab := w.page("ui.tab.about")
	w.caption(ab, "about.version")
	w.ver, _ = walk.NewLabel(ab)
	w.ver.SetText("TaskGlass " + w.host.Version())
	w.help(ab, "about.license")
	ll, _ := walk.NewLinkLabel(ab)
	ll.SetText(`<a href="https://github.com/AnderCMD/TaskGlass">github.com/AnderCMD/TaskGlass</a>`)
	ll.LinkActivated().Attach(func(l *walk.LinkLabelLink) {
		_ = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", l.URL()).Start()
	})
	w.cfgPath, _ = walk.NewLabel(ab)
	w.cfgPath.SetText(w.host.ConfigPath())
	walk.NewVSpacer(ab)
}

func (w *window) buildPreview(parent walk.Container) {
	col := vbox(parent, walk.Margins{}, 8)
	col.SetAlignment(walk.AlignHNearVNear)
	col.SetMinMaxSize(walk.Size{Width: previewW + 10}, walk.Size{Width: previewW + 40})
	w.caption(col, "ui.preview")
	row := hbox(col, 8)
	w.caption(row, "ui.scenario")
	w.scenCombo, _ = walk.NewComboBox(row)
	w.scenCombo.CurrentIndexChanged().Attach(func() {
		if w.syncing {
			return
		}
		if i := w.scenCombo.CurrentIndex(); i >= 0 {
			w.scenarioIdx = i
			w.render()
		}
	})
	w.view, _ = walk.NewImageView(col)
	w.view.SetMode(walk.ImageViewModeIdeal)
	w.view.SetMinMaxSize(walk.Size{Width: previewW, Height: previewH}, walk.Size{Width: previewW, Height: previewH})
	w.help(col, "ui.previewNote")
	walk.NewVSpacer(col)
}

func (w *window) buildFooter(parent walk.Container) {
	f := hbox(parent, 8)
	w.save = w.button(f, "ui.save", func() { w.saveNow() })
	w.revert = w.button(f, "ui.revert", func() {
		w.cfg, w.auto = clone(w.saved), w.savedAuto
		w.lang = i18n.Resolve(w.cfg.Language, i18n.SystemLocale())
		w.pull()
		w.applyText()
		w.afterChange()
	})
	w.status, _ = walk.NewLabel(f)
	walk.NewHSpacer(f)
	w.reset = w.button(f, "ui.reset", func() {
		if walk.MsgBox(w.mw, w.t("ui.title"), w.t("ui.resetConfirm"), walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
			return
		}
		w.cfg = norm(config.Default())
		w.lang = i18n.Resolve(w.cfg.Language, i18n.SystemLocale())
		w.pull()
		w.applyText()
		w.afterChange()
	})
}

// ---- state ------------------------------------------------------------------

func (w *window) applyText() {
	for _, f := range w.retext {
		f()
	}
	title := w.t("ui.title")
	if w.dirty() {
		title += " *"
	}
	w.mw.SetTitle(title)
	w.def.refresh()
	w.maxEd.refresh()
	w.appEd.refresh()
	w.appLimit.SetText(fmt.Sprintf(w.t("apps.limit"), config.MaxAppRules))
}

// pull copies the whole draft config into the controls.
func (w *window) pull() {
	w.syncing = true
	c := w.cfg
	w.enabled.SetChecked(c.Enabled)
	w.autostart.SetChecked(w.auto)
	selectID(w.language, languageChoices, string(c.Language))
	selectID(w.themeC, themeChoices, string(c.Theme))
	selectID(w.taskbars, taskbarChoices, string(c.Taskbars))
	selectID(w.trayIcon, trayChoices, string(c.TrayIcon))
	selectID(w.updMode, updateChoices, string(c.Updates))
	w.hotkeyOn.SetChecked(c.HotkeyEnabled)
	_ = w.hotkey.SetText(c.Hotkey)
	w.xaml.SetChecked(c.XamlFix)
	w.xaml.SetEnabled(w.host.XamlAvailable())
	w.max.SetChecked(c.Maximized.Enabled)
	w.flyOn.SetChecked(c.Surfaces.TrayOverflow.Enabled)
	w.flyOp.SetValue(int(c.Surfaces.TrayOverflow.Opacity))
	w.flyOpVal.SetText(fmt.Sprintf("%d %%", c.Surfaces.TrayOverflow.Opacity))
	w.flyMax.SetChecked(c.Surfaces.TrayOverflow.Maximized.Enabled)
	w.flyMaxOp.SetValue(int(c.Surfaces.TrayOverflow.Maximized.Opacity))
	w.flyMaxOpVal.SetText(fmt.Sprintf("%d %%", c.Surfaces.TrayOverflow.Maximized.Opacity))
	_ = w.poll.SetValue(float64(c.PollMs))
	_ = w.debounce.SetValue(float64(c.DebounceMs))
	_ = w.refresh.SetValue(float64(c.RefreshMs))
	selectID(w.logLevel, logChoices, c.LogLevel)
	w.syncing = false
	w.def.pull()
	w.maxEd.pull()
	w.rebuildApps(0)
	w.afterChange()
}

func (w *window) appNames() []string {
	names := make([]string, len(w.cfg.Apps))
	for i, r := range w.cfg.Apps {
		names[i] = r.Exe
		if names[i] == "" {
			names[i] = "…"
		}
	}
	return names
}

func (w *window) rebuildApps(sel int) {
	w.syncing = true
	_ = w.appList.SetModel(w.appNames())
	if sel >= len(w.cfg.Apps) {
		sel = len(w.cfg.Apps) - 1
	}
	w.appList.SetCurrentIndex(sel)
	w.syncing = false
	w.pullApp()
	w.appAdd.SetEnabled(len(w.cfg.Apps) < config.MaxAppRules)
}

// refreshAppNames updates the list captions while keeping the selection.
func (w *window) refreshAppNames(sel int) {
	w.syncing = true
	_ = w.appList.SetModel(w.appNames())
	w.appList.SetCurrentIndex(sel)
	w.syncing = false
}

func (w *window) pullApp() {
	i := w.appList.CurrentIndex()
	has := i >= 0 && i < len(w.cfg.Apps)
	w.syncing = true
	if has {
		_ = w.appExe.SetText(w.cfg.Apps[i].Exe)
	} else {
		_ = w.appExe.SetText("")
	}
	w.syncing = false
	w.appExe.SetEnabled(has)
	w.appRemove.SetEnabled(has)
	w.appPick.SetVisible(!has)
	w.appEd.root.SetEnabled(has)
	w.appEd.pull()
	w.rebuildScenarios()
	w.render()
}

// changed runs after every edit made by the user.
func (w *window) changed() {
	if l := i18n.Resolve(w.cfg.Language, i18n.SystemLocale()); l != w.lang {
		w.lang = l
		w.applyText()
	}
	w.afterChange()
}

func (w *window) afterChange() {
	w.maxGroup.SetEnabled(w.cfg.Maximized.Enabled)
	w.def.refresh()
	w.maxEd.refresh()
	w.appEd.refresh()
	valid := true
	if _, _, err := config.ParseHotkey(w.cfg.Hotkey); err != nil {
		valid = false
		w.hotkeyWarn.SetText("⚠ " + w.t("ui.hotkeyInvalid"))
	} else {
		w.hotkeyWarn.SetText("")
	}
	d := w.dirty()
	w.save.SetEnabled(d && valid)
	w.revert.SetEnabled(d)
	title := w.t("ui.title")
	if d {
		title += " *"
		w.status.SetText("")
	}
	w.mw.SetTitle(title)
	w.rebuildScenarios()
	w.render()
}

func (w *window) rebuildScenarios() {
	var sc []scenario
	sc = append(sc, scenario{label: w.t("ui.previewTaskbar"), ap: &w.cfg.Default})
	fo := w.cfg.Surfaces.TrayOverflow
	if w.cfg.Maximized.Enabled || (fo.Enabled && fo.Maximized.Enabled) {
		ap := &w.cfg.Default
		if w.cfg.Maximized.Enabled {
			ap = &w.cfg.Maximized.Appearance
		}
		sc = append(sc, scenario{label: w.t("ui.previewMax"), maximized: true, ap: ap})
	}
	for i := range w.cfg.Apps {
		if w.cfg.Apps[i].Exe != "" {
			sc = append(sc, scenario{label: fmt.Sprintf(w.t("ui.previewApp"), w.cfg.Apps[i].Exe), ap: &w.cfg.Apps[i].Appearance})
		}
	}
	w.scenarios = sc
	if w.scenarioIdx >= len(sc) {
		w.scenarioIdx = 0
	}
	names := make([]string, len(sc))
	for i, s := range sc {
		names[i] = s.label
	}
	w.syncing = true
	_ = w.scenCombo.SetModel(names)
	w.scenCombo.SetCurrentIndex(w.scenarioIdx)
	w.syncing = false
}

func (w *window) render() {
	if len(w.scenarios) == 0 || w.view == nil {
		return
	}
	sc := w.scenarios[w.scenarioIdx]
	avail := w.host.XamlAvailable()
	fly := w.cfg.Surfaces.TrayOverflow
	fly.Opacity = fly.Effective(sc.maximized) // the maximized/fullscreen rule shows in its own scenario
	img := preview.Render(preview.Input{
		Scenario:       preview.Scenario{Maximized: sc.maximized, Appearance: *sc.ap},
		Disabled:       !w.cfg.Enabled,
		Light:          w.light(),
		XamlFix:        w.cfg.XamlFix && avail,
		BlurViaAcrylic: winapi.BlurNeedsAcrylic(),
		Flyout:         fly,
		FlyoutActive:   w.cfg.Enabled && w.cfg.XamlFix && w.cfg.Surfaces.TrayOverflow.Enabled && avail,
	}, previewW, previewH)
	bmp, err := walk.NewBitmapFromImage(img)
	if err != nil {
		return
	}
	_ = w.view.SetImage(bmp)
	if w.bmp != nil {
		w.bmp.Dispose()
	}
	w.bmp = bmp
}

// saveNow applies the draft; it reports whether it succeeded.
func (w *window) saveNow() bool {
	c := clone(w.cfg)
	kept := c.Apps[:0]
	for _, r := range c.Apps {
		if r.Exe != "" {
			kept = append(kept, r)
		}
	}
	c.Apps = kept
	if err := w.host.Apply(c); err != nil {
		walk.MsgBox(w.mw, w.t("ui.title"), fmt.Sprintf(w.t("ui.saveFailed"), err.Error()), walk.MsgBoxOK|walk.MsgBoxIconError)
		return false
	}
	if w.auto != w.savedAuto {
		if err := w.host.SetAutostart(w.auto); err != nil {
			walk.MsgBox(w.mw, w.t("ui.title"), fmt.Sprintf(w.t("ui.saveFailed"), err.Error()), walk.MsgBoxOK|walk.MsgBoxIconError)
			return false
		}
	}
	w.cfg = norm(w.host.Config())
	w.saved = clone(w.cfg)
	w.auto = w.host.AutostartEnabled()
	w.savedAuto = w.auto
	w.pull()
	w.status.SetText("✓ " + w.t("ui.saved"))
	if w.toastTimer != nil {
		w.toastTimer.Stop()
	}
	w.toastTimer = time.AfterFunc(2500*time.Millisecond, func() {
		w.mw.Synchronize(func() {
			if !w.dirty() {
				w.status.SetText("")
			}
		})
	})
	return true
}
