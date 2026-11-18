//go:build windows

// Package winglass gives windows a uniform level of opacity: File Explorer
// windows and the windows of applications listed in the config.
//
// Explorer draws its file list and navigation pane itself, so they ignore the
// system backdrop; the title bar, tabs and command bar react to it but the
// content does not. What does reach everything is the window-level alpha of a
// layered window, which Windows applies to the whole window from outside the
// owning process. That is what this package uses: no code is loaded into
// the target processes.
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

// Wanted returns the alpha every rule asks for, keyed by the lower-case
// executable name ("explorer.exe" for File Explorer windows). Rules at 100 %
// are left out. The first app rule for an executable wins.
func Wanted(c config.Config) map[string]uint8 {
	out := map[string]uint8{}
	if !c.Enabled {
		return out
	}
	if c.Surfaces.Explorer.Enabled && c.Surfaces.Explorer.Opacity < 100 {
		out["explorer.exe"] = Alpha(int(c.Surfaces.Explorer.Opacity))
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
		var wins []winapi.HWND
		if exe == "explorer.exe" {
			wins = winapi.FileExplorerWindows()
		} else {
			wins = winapi.AppWindows(exe)
		}
		for _, h := range wins {
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

// RestoreAll makes every File Explorer window opaque, including ones changed
// by an earlier TaskGlass process that was killed. File Explorer windows are
// not layered by default, so this is safe to run at any time.
func RestoreAll() {
	for _, h := range winapi.FileExplorerWindows() {
		winapi.ClearWindowAlpha(h, true)
	}
}
