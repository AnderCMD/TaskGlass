//go:build windows

// Package winglass gives the windows of applications listed in the config a
// uniform level of opacity (the window_opacity of an app rule).
//
// It uses the window-level alpha of a layered window, which Windows applies to
// the whole window from outside the owning process, so no code is loaded into
// the target applications.
package winglass

import (
	"log/slog"
	"strings"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

type state struct {
	addedLayered bool  // WS_EX_LAYERED was added by us (and must be removed on restore)
	alpha        uint8 // last alpha applied
}

// Manager applies and restores window opacity.
type Manager struct {
	log  *slog.Logger
	seen map[winapi.HWND]state
}

// New creates a Manager.
func New(log *slog.Logger) *Manager {
	return &Manager{log: log, seen: map[winapi.HWND]state{}}
}

// Alpha converts an opacity percent (0..100) to a window alpha. 0 percent is
// clamped to a small value so a window never becomes impossible to see.
func Alpha(opacityPercent int) uint8 {
	if opacityPercent >= 100 {
		return 255
	}
	if opacityPercent < 10 {
		opacityPercent = 10
	}
	return uint8(opacityPercent * 255 / 100)
}

// Wanted returns the alpha every app rule asks for, keyed by the lower-case
// executable name. Rules at 100 % are left out. The first rule for an
// executable wins.
func Wanted(c config.Config) map[string]uint8 {
	out := map[string]uint8{}
	if !c.Enabled {
		return out
	}
	seen := map[string]bool{}
	for _, r := range c.Apps {
		exe := strings.ToLower(strings.TrimSpace(r.Exe))
		if seen[exe] {
			continue
		}
		seen[exe] = true
		if r.WindowOpacity == nil {
			continue
		}
		if *r.WindowOpacity >= 100 {
			delete(out, exe)
			continue
		}
		out[exe] = Alpha(int(*r.WindowOpacity))
	}
	return out
}

// Tick applies the wanted opacity to the matching windows and gives back the
// windows that no longer match any rule.
func (m *Manager) Tick(c config.Config) {
	wanted := Wanted(c)
	desired := map[winapi.HWND]uint8{}
	for exe, alpha := range wanted {
		for _, h := range winapi.AppWindows(exe) {
			desired[h] = alpha
		}
	}
	for h, alpha := range desired {
		st, ok := m.seen[h]
		if ok && st.alpha == alpha {
			continue
		}
		if !ok {
			st.addedLayered = !winapi.IsLayered(h)
		}
		if err := winapi.SetWindowAlpha(h, alpha); err != nil {
			m.log.Debug("window alpha failed", "hwnd", h, "err", err)
			continue
		}
		st.alpha = alpha
		m.seen[h] = st
	}
	for h, st := range m.seen {
		if _, still := desired[h]; !still {
			if winapi.IsWindow(h) { // rule removed or app closed: give the window back
				winapi.ClearWindowAlpha(h, st.addedLayered)
			}
			delete(m.seen, h)
		}
	}
}

// Restore makes every window we touched opaque again.
func (m *Manager) Restore() {
	for h, st := range m.seen {
		if winapi.IsWindow(h) {
			winapi.ClearWindowAlpha(h, st.addedLayered)
		}
		delete(m.seen, h)
	}
}
