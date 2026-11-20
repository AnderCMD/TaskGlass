//go:build windows

// Package ui is TaskGlass' native settings window: real Win32 controls (via
// walk) next to a live preview of the taskbar and the hidden-icons flyout.
//
// All windows live on one dedicated OS thread, as Win32 requires; Show may be
// called from any goroutine.
package ui

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/win"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/i18n"
	"github.com/AnderCMD/TaskGlass/internal/theme"
	"github.com/AnderCMD/TaskGlass/internal/update"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

// Host is what the settings window needs from the running application.
type Host interface {
	Config() config.Config
	// Apply validates, saves and applies a whole configuration.
	Apply(config.Config) error
	AutostartEnabled() bool
	SetAutostart(on bool) error
	OpenFile(which string) error // "config" or "log"
	Version() string
	ConfigPath() string
	XamlAvailable() bool
	// CheckUpdate looks for a newer release; InstallUpdate installs the one found.
	CheckUpdate() (update.Info, error)
	InstallUpdate() error
}

// UI owns the settings window.
type UI struct {
	host Host
	// OnError is called when the window cannot be created.
	OnError func(error)

	once sync.Once
	req  chan struct{}

	mu  sync.Mutex
	cur *window
}

// New creates the UI; nothing is shown until Show.
func New(h Host) *UI { return &UI{host: h, req: make(chan struct{}, 1)} }

// Show opens the settings window, or brings the open one to the front.
func (u *UI) Show() {
	u.once.Do(func() { go u.loop() })
	u.mu.Lock()
	cur := u.cur
	u.mu.Unlock()
	if cur != nil {
		cur.mw.Synchronize(func() {
			h := cur.mw.Handle()
			win.ShowWindow(h, win.SW_RESTORE)
			win.SetForegroundWindow(h)
		})
		return
	}
	select {
	case u.req <- struct{}{}:
	default:
	}
}

func (u *UI) loop() {
	runtime.LockOSThread()
	for range u.req {
		w, err := newWindow(u.host)
		if err != nil {
			if u.OnError != nil {
				u.OnError(err)
			}
			continue
		}
		u.mu.Lock()
		u.cur = w
		u.mu.Unlock()
		w.center()
		w.mw.Show()
		w.mw.Run()
		w.dispose()
		u.mu.Lock()
		u.cur = nil
		u.mu.Unlock()
	}
}

// ---- window ---------------------------------------------------------------

const (
	previewW, previewH = 560, 300
	resIconID          = 1 // the icon linked into taskglass.exe
)

type scenario struct {
	label     string
	maximized bool
	ap        *config.Appearance
}

type window struct {
	host Host
	mw   *walk.MainWindow

	cfg, saved      config.Config
	auto, savedAuto bool
	lang            i18n.Lang
	systemLight     bool
	syncing         bool
	retext          []func()
	scenarios       []scenario
	scenarioIdx     int
	bmp             *walk.Bitmap
	toastTimer      *time.Timer
	closeConfirmed  bool

	tabs                                    *walk.TabWidget
	pages                                   []*walk.TabPage
	enabled, autostart, hotkeyOn, xaml, max *walk.CheckBox
	language, themeC, taskbars, trayIcon    *walk.ComboBox
	hotkey                                  *walk.LineEdit
	hotkeyWarn                              *walk.Label
	updMode                                 *walk.ComboBox
	updCheck, updInstall                    *walk.PushButton
	updMsg                                  *walk.Label
	flyMax                                  *walk.CheckBox
	flyMaxOp                                *walk.Slider
	flyMaxOpVal                             *walk.Label
	def, maxEd                              *appEditor
	maxGroup                                *walk.Composite
	flyOn                                   *walk.CheckBox
	flyOp                                   *walk.Slider
	flyOpVal                                *walk.Label
	appList                                 *walk.ListBox
	appExe                                  *walk.LineEdit
	appEd                                   *appEditor
	appEditorBox                            *walk.Composite
	appAdd, appRemove                       *walk.PushButton
	appLimit, appPick                       *walk.Label
	poll, debounce, refresh                 *walk.NumberEdit
	logLevel                                *walk.ComboBox
	scenCombo                               *walk.ComboBox
	view                                    *walk.ImageView
	save, revert, reset                     *walk.PushButton
	status                                  *walk.Label
	ver, cfgPath                            *walk.Label
}

func newWindow(h Host) (*window, error) {
	w := &window{host: h, systemLight: theme.SystemIsLight()}
	w.cfg = norm(h.Config())
	w.saved = clone(w.cfg)
	w.auto = h.AutostartEnabled()
	w.savedAuto = w.auto
	w.lang = i18n.Resolve(w.cfg.Language, i18n.SystemLocale())

	mw, err := walk.NewMainWindow()
	if err != nil {
		return nil, err
	}
	w.mw = mw
	if ic, err := walk.NewIconFromResourceId(resIconID); err == nil {
		mw.SetIcon(ic)
	}
	mw.SetMinMaxSize(walk.Size{Width: 1000, Height: 560}, walk.Size{})
	mw.SetSize(walk.Size{Width: 1080, Height: 580})
	root := walk.NewVBoxLayout()
	root.SetMargins(walk.Margins{HNear: 12, VNear: 12, HFar: 12, VFar: 10})
	root.SetSpacing(10)
	mw.SetLayout(root)

	body, _ := walk.NewComposite(mw)
	bl := walk.NewHBoxLayout()
	bl.SetMargins(walk.Margins{})
	bl.SetSpacing(14)
	body.SetLayout(bl)

	w.buildTabs(body)
	w.buildPreview(body)
	bl.SetStretchFactor(w.tabs, 11)

	walk.NewVSpacer(mw) // keeps the footer at the bottom when the window grows
	w.buildFooter(mw)

	mw.Closing().Attach(func(canceled *bool, _ walk.CloseReason) {
		if w.dirty() && !w.closeConfirmed {
			switch walk.MsgBox(mw, w.t("ui.title"), w.t("ui.closeUnsaved"), walk.MsgBoxYesNoCancel|walk.MsgBoxIconQuestion) {
			case walk.DlgCmdYes:
				if !w.saveNow() {
					*canceled = true
				}
			case walk.DlgCmdCancel:
				*canceled = true
			}
		}
	})

	w.pull()
	w.applyText()
	return w, nil
}

func (w *window) center() {
	sw := int(win.GetSystemMetrics(win.SM_CXSCREEN))
	sh := int(win.GetSystemMetrics(win.SM_CYSCREEN))
	b := w.mw.BoundsPixels()
	w.mw.SetBoundsPixels(walk.Rectangle{X: (sw - b.Width) / 2, Y: (sh - b.Height) / 2, Width: b.Width, Height: b.Height})
}

func (w *window) dispose() {
	if w.toastTimer != nil {
		w.toastTimer.Stop()
	}
	if w.bmp != nil {
		w.bmp.Dispose()
	}
	w.mw.Dispose()
}

// ---- helpers --------------------------------------------------------------

func norm(c config.Config) config.Config {
	if c.Apps == nil {
		c.Apps = []config.AppRule{}
	}
	return c
}

func clone(c config.Config) config.Config {
	c.Apps = append([]config.AppRule{}, c.Apps...)
	return c
}

func (w *window) t(key string) string { return i18n.T(w.lang, key) }

// text registers a setter that is re-run whenever the language changes.
func (w *window) text(set func(string), key string) {
	f := func() { set(w.t(key)) }
	w.retext = append(w.retext, f)
	f()
}

func (w *window) light() bool { return w.cfg.Theme.Light(w.systemLight) }

func (w *window) dirty() bool {
	a, _ := json.Marshal(w.cfg)
	b, _ := json.Marshal(w.saved)
	return string(a) != string(b) || w.auto != w.savedAuto
}

func hexOf(r, g, b uint8) string { return fmt.Sprintf("#%02X%02X%02X", r, g, b) }

func label(parent walk.Container, wrap bool) *walk.TextLabel {
	l, _ := walk.NewTextLabel(parent)
	_ = wrap
	return l
}

func vbox(parent walk.Container, m walk.Margins, spacing int) *walk.Composite {
	c, _ := walk.NewComposite(parent)
	l := walk.NewVBoxLayout()
	l.SetMargins(m)
	l.SetSpacing(spacing)
	c.SetLayout(l)
	return c
}

func hbox(parent walk.Container, spacing int) *walk.Composite {
	c, _ := walk.NewComposite(parent)
	l := walk.NewHBoxLayout()
	l.SetMargins(walk.Margins{})
	l.SetSpacing(spacing)
	c.SetLayout(l)
	return c
}

// row adds "label ........ control" style rows; help is an optional dim line.
func (w *window) caption(parent walk.Container, key string) {
	l, _ := walk.NewTextLabel(parent)
	w.text(func(s string) { l.SetText(s) }, key)
}

func (w *window) check(parent walk.Container, key string, get func() bool, set func(bool)) *walk.CheckBox {
	cb, _ := walk.NewCheckBox(parent)
	cb.SetAlignment(walk.AlignHNearVNear)
	w.text(func(s string) { cb.SetText(s) }, key)
	cb.CheckedChanged().Attach(func() {
		if w.syncing {
			return
		}
		set(cb.Checked())
		w.changed()
	})
	return cb
}

func (w *window) help(parent walk.Container, key string) *walk.TextLabel {
	l, _ := walk.NewTextLabel(parent)
	w.text(func(s string) { l.SetText(s) }, key)
	return l
}

type choice struct {
	id  string
	key string
}

// combo builds a ComboBox whose items are translated and re-translated on language change.
func (w *window) combo(parent walk.Container, items []choice, get func() string, set func(string)) *walk.ComboBox {
	cb, _ := walk.NewComboBox(parent)
	fill := func() {
		w.syncing = true
		names := make([]string, len(items))
		sel := -1
		for i, it := range items {
			names[i] = it.id // raw label unless it has a translation key
			if it.key != "" {
				names[i] = w.t(it.key)
			}
			if it.id == get() {
				sel = i
			}
		}
		_ = cb.SetModel(names)
		cb.SetCurrentIndex(sel)
		w.syncing = false
	}
	w.retext = append(w.retext, fill)
	cb.CurrentIndexChanged().Attach(func() {
		if w.syncing {
			return
		}
		if i := cb.CurrentIndex(); i >= 0 && i < len(items) {
			set(items[i].id)
			w.changed()
		}
	})
	return cb
}

func selectID(cb *walk.ComboBox, items []choice, id string) {
	for i, it := range items {
		if it.id == id {
			cb.SetCurrentIndex(i)
			return
		}
	}
}

func (w *window) slider(parent walk.Container, get func() int, set func(int)) (*walk.Slider, *walk.Label) {
	row := hbox(parent, 8)
	s, _ := walk.NewSlider(row)
	s.SetRange(0, 100)
	s.SetTracking(true)
	val, _ := walk.NewLabel(row)
	val.SetMinMaxSize(walk.Size{Width: 40}, walk.Size{Width: 40})
	s.ValueChanged().Attach(func() {
		val.SetText(fmt.Sprintf("%d %%", s.Value()))
		if w.syncing {
			return
		}
		set(s.Value())
		w.changed()
	})
	return s, val
}

func (w *window) button(parent walk.Container, key string, click func()) *walk.PushButton {
	b, _ := walk.NewPushButton(parent)
	w.text(func(s string) { b.SetText(s) }, key)
	b.Clicked().Attach(click)
	return b
}

// ---- appearance editor (accent, colour, strength) -------------------------

var accentChoices = func() []choice {
	var out []choice
	for _, a := range config.Accents {
		out = append(out, choice{string(a), "accent." + string(a)})
	}
	return out
}()

type appEditor struct {
	w      *window
	get    func() *config.Appearance
	root   *walk.Composite
	accent *walk.ComboBox
	desc   *walk.TextLabel
	auto   *walk.CheckBox
	swatch *walk.Composite
	hexLbl *walk.Label
	pick   *walk.PushButton
	slider *walk.Slider
	value  *walk.Label
	note   *walk.TextLabel
	sw     *walk.SolidColorBrush
}

func (w *window) newAppEditor(parent walk.Container, get func() *config.Appearance) *appEditor {
	orig := get
	get = func() *config.Appearance { // never nil: an unselected rule edits a throw-away copy
		if ap := orig(); ap != nil {
			return ap
		}
		return &config.Appearance{Accent: config.AccentClear, Color: config.ColorAuto}
	}
	e := &appEditor{w: w, get: get}
	e.root = vbox(parent, walk.Margins{}, 6)

	w.caption(e.root, "appear.style")
	e.accent = w.combo(e.root, accentChoices, func() string { return string(get().Accent) }, func(id string) { get().Accent = config.Accent(id) })
	e.desc, _ = walk.NewTextLabel(e.root)
	e.accent.CurrentIndexChanged().Attach(func() { e.refresh() })

	w.caption(e.root, "appear.color")
	crow := hbox(e.root, 8)
	e.auto = w.check(crow, "appear.colorAuto", func() bool { return strings.EqualFold(get().Color, config.ColorAuto) }, func(on bool) {
		if on {
			get().Color = config.ColorAuto
		} else {
			r, g, b := get().RGB(w.light())
			get().Color = hexOf(r, g, b)
		}
	})
	e.swatch, _ = walk.NewComposite(crow)
	e.swatch.SetLayout(walk.NewHBoxLayout())
	e.swatch.SetMinMaxSize(walk.Size{Width: 28, Height: 22}, walk.Size{Width: 28, Height: 22})
	e.hexLbl, _ = walk.NewLabel(crow)
	e.pick = w.button(crow, "appear.pickColor", func() { e.chooseColor() })
	walk.NewHSpacer(crow)

	w.caption(e.root, "appear.opacity")
	e.slider, e.value = w.slider(e.root, func() int { return int(get().Opacity) }, func(v int) { get().Opacity = config.Percent(v) })
	e.note, _ = walk.NewTextLabel(e.root)
	w.text(func(s string) { e.note.SetText(s) }, "appear.blurNote")
	w.text(func(string) { e.refresh() }, "appear.opacity.help")
	return e
}

func (e *appEditor) chooseColor() {
	ap := e.get()
	r, g, b := ap.RGB(e.w.light())
	var custom [16]win.COLORREF
	cc := win.CHOOSECOLOR{
		LStructSize:  uint32(unsafeSizeofChooseColor()),
		HwndOwner:    e.w.mw.Handle(),
		RgbResult:    win.COLORREF(uint32(r) | uint32(g)<<8 | uint32(b)<<16),
		LpCustColors: &custom,
		Flags:        win.CC_RGBINIT | win.CC_FULLOPEN,
	}
	if !win.ChooseColor(&cc) {
		return
	}
	v := uint32(cc.RgbResult)
	ap.Color = hexOf(uint8(v), uint8(v>>8), uint8(v>>16))
	e.w.changed()
}

// refresh updates the controls that depend on the current appearance.
func (e *appEditor) refresh() {
	ap := e.get()
	if ap == nil {
		return
	}
	e.desc.SetText(e.w.t("style." + string(ap.Accent) + ".desc"))
	auto := strings.EqualFold(ap.Color, config.ColorAuto)
	r, g, b := ap.RGB(e.w.light())
	e.hexLbl.SetText(hexOf(r, g, b))
	e.pick.SetEnabled(!auto && ap.Accent != config.AccentNormal)
	e.slider.SetEnabled(ap.Accent != config.AccentNormal)
	e.auto.SetEnabled(ap.Accent != config.AccentNormal)
	e.note.SetVisible(ap.Accent == config.AccentBlur && winapi.BlurNeedsAcrylic())
	if br, err := walk.NewSolidColorBrush(walk.RGB(r, g, b)); err == nil {
		e.swatch.SetBackground(br)
		if e.sw != nil {
			e.sw.Dispose()
		}
		e.sw = br
	}
}

// pull copies the config into the controls.
func (e *appEditor) pull() {
	ap := e.get()
	if ap == nil {
		return
	}
	e.w.syncing = true
	selectID(e.accent, accentChoices, string(ap.Accent))
	e.auto.SetChecked(strings.EqualFold(ap.Color, config.ColorAuto))
	e.slider.SetValue(int(ap.Opacity))
	e.value.SetText(fmt.Sprintf("%d %%", ap.Opacity))
	e.w.syncing = false
	e.refresh()
}
