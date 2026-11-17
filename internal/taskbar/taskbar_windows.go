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

type applied struct {
	ap config.Appearance
	at time.Time
}

// Manager tracks what has been applied to each taskbar so it only re-applies
// when something changed or the refresh interval elapsed.
type Manager struct {
	log  *slog.Logger
	seen map[winapi.HWND]applied
}

// New creates a Manager.
func New(log *slog.Logger) *Manager {
	return &Manager{log: log, seen: map[winapi.HWND]applied{}}
}

// Tick samples the desktop and updates every taskbar. Call it periodically
// from a single goroutine.
func (m *Manager) Tick(c config.Config) {
	bars := winapi.Taskbars()
	alive := make(map[winapi.HWND]bool, len(bars))
	var maxed map[uintptr]bool
	if c.Maximized.Enabled {
		maxed = winapi.MaximizedMonitors()
	}
	start := c.StartMenu.Enabled && winapi.StartMenuActive()
	refresh := time.Duration(c.RefreshMs) * time.Millisecond

	for _, h := range bars {
		alive[h] = true
		ctx := state.Context{StartMenu: start, Maximized: maxed[winapi.Monitor(h)]}
		ap, kind := state.Resolve(c, ctx)

		prev, ok := m.seen[h]
		if ok && prev.ap == ap && time.Since(prev.at) < refresh {
			continue
		}
		if err := apply(h, ap); err != nil {
			m.log.Warn("apply failed", "hwnd", h, "err", err)
			continue
		}
		if !ok || prev.ap != ap {
			m.log.Debug("applied", "hwnd", h, "rule", kind, "accent", ap.Accent)
		}
		m.seen[h] = applied{ap: ap, at: time.Now()}
	}
	for h := range m.seen {
		if !alive[h] {
			delete(m.seen, h) // taskbar went away (Explorer restart, monitor unplugged)
		}
	}
}

// Restore puts every taskbar back to the stock Windows look.
func (m *Manager) Restore() {
	for _, h := range winapi.Taskbars() {
		if err := apply(h, config.Appearance{Accent: config.AccentNormal, Color: "#000000"}); err != nil {
			m.log.Warn("restore failed", "hwnd", h, "err", err)
		}
		winapi.Redraw(h)
	}
	m.seen = map[winapi.HWND]applied{}
}

func apply(h winapi.HWND, ap config.Appearance) error {
	state, color := toPolicy(ap)
	return winapi.SetAccent(h, state, color)
}

// toPolicy converts an Appearance to an accent state and 0xAABBGGRR colour.
func toPolicy(ap config.Appearance) (winapi.AccentState, uint32) {
	r, g, b, _ := config.ParseColor(ap.Color)
	alpha := uint32(ap.Opacity) * 255 / 100
	var st winapi.AccentState
	switch ap.Accent {
	case config.AccentNormal:
		return winapi.AccentDisabled, 0
	case config.AccentClear:
		st = winapi.AccentEnableTransparent
		if alpha == 0 {
			alpha = 1 // alpha 0 renders opaque black on some builds
		}
	case config.AccentBlur:
		st = winapi.AccentEnableBlurBehind
	case config.AccentAcrylic:
		st = winapi.AccentEnableAcrylicBlurBG
		if alpha == 0 {
			alpha = 1
		}
	case config.AccentOpaque:
		st = winapi.AccentEnableGradient
	}
	return st, alpha<<24 | uint32(b)<<16 | uint32(g)<<8 | uint32(r)
}
