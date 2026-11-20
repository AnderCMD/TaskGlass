//go:build windows

// Package taskbar applies configured appearances to every Windows taskbar and
// restores the stock look on exit.
package taskbar

import (
	"log/slog"
	"time"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/state"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

// policy is what was actually sent to a taskbar: the resolved accent state
// and colour, so theme changes that alter an "auto" colour are noticed.
type policy struct {
	state winapi.AccentState
	color uint32
}

type applied struct {
	p  policy
	at time.Time
}

// Manager tracks what has been applied to each taskbar so it only re-applies
// when something changed or the refresh interval elapsed.
type Manager struct {
	log  *slog.Logger
	seen map[winapi.HWND]applied
	// anyMaximized is whether any monitor had a maximized or fullscreen window
	// at the last Tick that needed to know.
	anyMaximized bool
}

// New creates a Manager.
func New(log *slog.Logger) *Manager {
	return &Manager{log: log, seen: map[winapi.HWND]applied{}}
}

// Tick samples the desktop and updates every selected taskbar. Call it from a
// single goroutine. light says whether the Windows theme is light (it picks
// "auto" colours). It reports whether any taskbar currently has a non-stock
// appearance, which decides whether the XAML fill must be hidden too.
func (m *Manager) Tick(c config.Config, light bool) (styled bool) {
	all := winapi.Taskbars()
	bars := all
	if c.Taskbars == config.TaskbarsPrimary && len(bars) > 1 {
		bars = bars[:1]
	}
	selected := make(map[winapi.HWND]bool, len(bars))
	alive := make(map[winapi.HWND]bool, len(all))
	for _, h := range all {
		alive[h] = true
	}

	needMax, needFg := state.NeedsWindowScan(c)
	var maxed map[uintptr]bool
	if needMax {
		maxed = winapi.MaximizedMonitors()
		m.anyMaximized = len(maxed) > 0
	} else {
		m.anyMaximized = false
	}
	var fgExe string
	var fgMon uintptr
	if needFg {
		fgExe, fgMon = winapi.Foreground()
	}
	refresh := time.Duration(c.RefreshMs) * time.Millisecond

	for _, h := range bars {
		selected[h] = true
		mon := winapi.Monitor(h)
		ctx := state.Context{Maximized: maxed[mon]}
		if needFg && fgMon == mon {
			ctx.ForegroundExe = fgExe
		}
		ap, kind := state.Resolve(c, ctx)
		if ap.Accent != config.AccentNormal {
			styled = true
		}
		want := toPolicy(ap, light)

		prev, ok := m.seen[h]
		if ok && prev.p == want && time.Since(prev.at) < refresh {
			continue
		}
		if err := winapi.SetAccent(h, want.state, want.color); err != nil {
			m.log.Warn("apply failed", "hwnd", h, "err", err)
			continue
		}
		if !ok || prev.p != want {
			m.log.Debug("applied", "hwnd", h, "rule", kind, "accent", ap.Accent)
		}
		m.seen[h] = applied{p: want, at: time.Now()}
	}

	// Forget taskbars that went away (Explorer restart, monitor unplugged) and
	// give back any that are no longer selected (taskbars switched to "primary").
	for h := range m.seen {
		switch {
		case !alive[h]:
			delete(m.seen, h)
		case !selected[h]:
			m.restoreOne(h)
			delete(m.seen, h)
		}
	}
	return styled
}

// AnyMaximized reports whether a window was maximized or fullscreen on any
// monitor at the last Tick. Call it from the same goroutine as Tick.
func (m *Manager) AnyMaximized() bool { return m.anyMaximized }

// Restore puts every taskbar back to the stock Windows look.
func (m *Manager) Restore() {
	for _, h := range winapi.Taskbars() {
		m.restoreOne(h)
	}
	m.seen = map[winapi.HWND]applied{}
}

func (m *Manager) restoreOne(h winapi.HWND) {
	if err := winapi.SetAccent(h, winapi.AccentDisabled, 0); err != nil {
		m.log.Warn("restore failed", "hwnd", h, "err", err)
	}
	winapi.Redraw(h)
}

// blurViaAcrylic is true where the classic blur-behind accent is unusable
// (Windows 11, build 22000 and later).
var blurViaAcrylic = winapi.BlurNeedsAcrylic()

// lighten mixes a colour towards white by t (0..1).
func lighten(r, g, b uint8, t float64) (uint8, uint8, uint8) {
	f := func(c uint8) uint8 { return uint8(float64(c) + (255-float64(c))*t + 0.5) }
	return f(r), f(g), f(b)
}

// toPolicy converts an Appearance to an accent state and 0xAABBGGRR colour.
func toPolicy(ap config.Appearance, light bool) policy {
	if ap.Accent == config.AccentNormal {
		return policy{winapi.AccentDisabled, 0}
	}
	r, g, b := ap.RGB(light)
	alpha := uint32(ap.Opacity) * 255 / 100
	var st winapi.AccentState
	switch ap.Accent {
	case config.AccentClear:
		st = winapi.AccentEnableTransparent
		if alpha == 0 {
			alpha = 1 // alpha 0 renders opaque black on some builds
		}
	case config.AccentBlur:
		if blurViaAcrylic {
			// Windows 11 draws the classic blur-behind accent as solid black, so
			// blur is produced by the acrylic backdrop. That backdrop takes its
			// brightness from the tint colour and ignores the alpha, so the
			// strength is applied by mixing the tint towards white: 0 is a light
			// frosted glass, 100 the full tint colour.
			st = winapi.AccentEnableAcrylicBlurBG
			mix := 0.5 * (1 - float64(ap.Opacity)/100)
			r, g, b = lighten(r, g, b, mix)
			alpha = 1
			break
		}
		st = winapi.AccentEnableBlurBehind
		if alpha == 0 {
			alpha = 1 // same quirk as clear: a fully zero tint can render opaque black
		}
	case config.AccentAcrylic:
		st = winapi.AccentEnableAcrylicBlurBG
		if alpha == 0 {
			alpha = 1
		}
	case config.AccentOpaque:
		st = winapi.AccentEnableGradient
	}
	return policy{st, alpha<<24 | uint32(b)<<16 | uint32(g)<<8 | uint32(r)}
}
