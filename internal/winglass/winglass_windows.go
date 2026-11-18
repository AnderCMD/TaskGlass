//go:build windows

// Package winglass gives File Explorer windows a uniform level of opacity.
//
// Explorer draws its file list and navigation pane itself, so they ignore the
// system backdrop; the title bar, tabs and command bar react to it but the
// content does not. What does reach everything is the window-level alpha of a
// layered window, which Windows applies to the whole window from outside the
// owning process. That is what this package uses: no code is loaded into
// Explorer.
package winglass

import (
	"log/slog"

	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

type state struct {
	addedLayered bool  // WS_EX_LAYERED was added by us (and must be removed on restore)
	alpha        uint8 // last alpha applied
}

// Manager applies and restores the opacity of File Explorer windows.
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

// Tick applies the opacity to every File Explorer window when enabled, and
// gives windows back when it is disabled or at 100 %.
func (m *Manager) Tick(enabled bool, opacityPercent int) {
	if !enabled || opacityPercent >= 100 {
		m.Restore()
		return
	}
	alpha := Alpha(opacityPercent)
	live := map[winapi.HWND]bool{}
	for _, h := range winapi.FileExplorerWindows() {
		live[h] = true
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
		if !live[h] { // closed or hidden: nothing to give back
			_ = st
			delete(m.seen, h)
		}
	}
}

// Restore makes every window we touched opaque again.
func (m *Manager) Restore() {
	for h, st := range m.seen {
		winapi.ClearWindowAlpha(h, st.addedLayered)
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
