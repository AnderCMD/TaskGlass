//go:build windows

package theme

import (
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
)

// Win32 menus (the tray icon's right-click menu) only follow dark mode when the
// process opts in through two undocumented uxtheme exports, found by ordinal.
// If they are missing (very old Windows builds) the menu keeps the light look.
const (
	ordSetPreferredAppMode = 135
	ordFlushMenuThemes     = 136

	modeForceDark  = 2
	modeForceLight = 3
)

var (
	menuOnce  sync.Once
	setMode   uintptr
	flushMenu uintptr
)

func loadMenuProcs() {
	h, err := windows.LoadLibrary("uxtheme.dll")
	if err != nil {
		return
	}
	setMode, _ = windows.GetProcAddressByOrdinal(h, ordSetPreferredAppMode)
	flushMenu, _ = windows.GetProcAddressByOrdinal(h, ordFlushMenuThemes)
}

// ApplyMenuTheme makes the process's own popup menus light or dark. It is cheap
// and safe to call repeatedly, from any thread.
func ApplyMenuTheme(light bool) {
	menuOnce.Do(loadMenuProcs)
	if setMode == 0 {
		return
	}
	mode := uintptr(modeForceDark)
	if light {
		mode = modeForceLight
	}
	syscall.SyscallN(setMode, mode)
	if flushMenu != 0 {
		syscall.SyscallN(flushMenu)
	}
}
