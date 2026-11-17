//go:build windows

// Package winapi is a thin, dependency-free wrapper over the handful of Win32
// calls TaskGlass needs. Everything that touches the OS lives here.
package winapi

import (
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// HWND is a window handle.
type HWND = uintptr

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	dwmapi = windows.NewLazySystemDLL("dwmapi.dll")

	procSetWCA          = user32.NewProc("SetWindowCompositionAttribute")
	procFindWindowW     = user32.NewProc("FindWindowW")
	procFindWindowExW   = user32.NewProc("FindWindowExW")
	procEnumWindows     = user32.NewProc("EnumWindows")
	procIsWindowVisible = user32.NewProc("IsWindowVisible")
	procIsIconic        = user32.NewProc("IsIconic")
	procIsZoomed        = user32.NewProc("IsZoomed")
	procIsWindow        = user32.NewProc("IsWindow")
	procGetWindowLongW  = user32.NewProc("GetWindowLongW")
	procMonitorFromWnd  = user32.NewProc("MonitorFromWindow")
	procGetForeground   = user32.NewProc("GetForegroundWindow")
	procGetClassNameW   = user32.NewProc("GetClassNameW")
	procGetWindowPID    = user32.NewProc("GetWindowThreadProcessId")
	procSendMessageW    = user32.NewProc("SendMessageW")

	procDwmGetWindowAttribute = dwmapi.NewProc("DwmGetWindowAttribute")
)

const (
	wcaAccentPolicy        = 19
	gwlExStyle             = -20
	wsExToolWindow         = 0x00000080
	monitorDefaultToNull   = 0
	dwmwaCloaked           = 14
	wmDwmCompositionChange = 0x031E
)

// AccentState mirrors the undocumented ACCENT_STATE enum.
type AccentState uint32

const (
	AccentDisabled            AccentState = 0
	AccentEnableGradient      AccentState = 1
	AccentEnableTransparent   AccentState = 2
	AccentEnableBlurBehind    AccentState = 3
	AccentEnableAcrylicBlurBG AccentState = 4
)

// accentPolicy mirrors ACCENT_POLICY.
type accentPolicy struct {
	State  uint32
	Flags  uint32
	Color  uint32 // 0xAABBGGRR
	AnimID uint32
}

// windowCompositionAttribData mirrors WINDOWCOMPOSITIONATTRIBDATA.
type windowCompositionAttribData struct {
	Attrib uint32
	Data   uintptr
	Size   uintptr
}

// SetAccent applies an accent policy to hwnd. color is 0xAABBGGRR.
func SetAccent(hwnd HWND, state AccentState, color uint32) error {
	flags := uint32(0)
	if state != AccentDisabled && state != AccentEnableGradient {
		flags = 2 // draw all borders
	}
	p := accentPolicy{State: uint32(state), Flags: flags, Color: color}
	d := windowCompositionAttribData{
		Attrib: wcaAccentPolicy,
		Data:   uintptr(unsafe.Pointer(&p)),
		Size:   unsafe.Sizeof(p),
	}
	r, _, err := procSetWCA.Call(hwnd, uintptr(unsafe.Pointer(&d)))
	if r == 0 {
		return err
	}
	return nil
}

// Redraw asks the window to re-evaluate its composition, which makes Explorer
// repaint the stock taskbar after an accent policy is disabled.
func Redraw(hwnd HWND) {
	procSendMessageW.Call(hwnd, wmDwmCompositionChange, 0, 0)
}

// IsWindow reports whether hwnd still refers to a live window.
func IsWindow(hwnd HWND) bool {
	r, _, _ := procIsWindow.Call(hwnd)
	return r != 0
}

func utf16(s string) *uint16 {
	p, _ := windows.UTF16PtrFromString(s)
	return p
}

// Taskbars returns the primary taskbar and every secondary (other monitor)
// taskbar currently present.
func Taskbars() []HWND {
	var out []HWND
	if h, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(utf16("Shell_TrayWnd"))), 0); h != 0 {
		out = append(out, h)
	}
	var prev uintptr
	for {
		h, _, _ := procFindWindowExW.Call(0, prev, uintptr(unsafe.Pointer(utf16("Shell_SecondaryTrayWnd"))), 0)
		if h == 0 {
			break
		}
		out = append(out, h)
		prev = h
	}
	return out
}

// Monitor returns the monitor handle that contains hwnd.
func Monitor(hwnd HWND) uintptr {
	m, _, _ := procMonitorFromWnd.Call(hwnd, monitorDefaultToNull)
	return m
}

var (
	enumMu sync.Mutex
	enumFn func(HWND) bool
	enumCB = syscall.NewCallback(func(h, _ uintptr) uintptr {
		if enumFn(h) {
			return 1
		}
		return 0
	})
)

// enumTopLevel calls fn for every top-level window until fn returns false.
func enumTopLevel(fn func(HWND) bool) {
	enumMu.Lock()
	defer enumMu.Unlock()
	enumFn = fn
	procEnumWindows.Call(enumCB, 0)
	enumFn = nil
}

func className(hwnd HWND) string {
	var buf [256]uint16
	n, _, _ := procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf[:n])
}

func cloaked(hwnd HWND) bool {
	var v uint32
	hr, _, _ := procDwmGetWindowAttribute.Call(hwnd, dwmwaCloaked, uintptr(unsafe.Pointer(&v)), unsafe.Sizeof(v))
	return hr == 0 && v != 0
}

// shellClasses are windows that are maximized-sized but are not "apps".
var shellClasses = map[string]bool{
	"Progman":                true,
	"WorkerW":                true,
	"Shell_TrayWnd":          true,
	"Shell_SecondaryTrayWnd": true,
}

// MaximizedMonitors returns the set of monitors that currently show a
// maximized, visible, non-minimized, non-cloaked application window.
func MaximizedMonitors() map[uintptr]bool {
	out := map[uintptr]bool{}
	enumTopLevel(func(h HWND) bool {
		if v, _, _ := procIsWindowVisible.Call(h); v == 0 {
			return true
		}
		if v, _, _ := procIsIconic.Call(h); v != 0 {
			return true
		}
		if v, _, _ := procIsZoomed.Call(h); v == 0 {
			return true
		}
		ex, _, _ := procGetWindowLongW.Call(h, uintptr(gwlExStyle&0xFFFFFFFF))
		if ex&wsExToolWindow != 0 {
			return true
		}
		if shellClasses[className(h)] || cloaked(h) {
			return true
		}
		if m := Monitor(h); m != 0 {
			out[m] = true
		}
		return true
	})
	return out
}

var startHosts = map[string]bool{
	"startmenuexperiencehost.exe": true,
	"searchhost.exe":              true,
	"searchapp.exe":               true,
}

// StartMenuActive reports whether the Start menu or Search currently has focus.
func StartMenuActive() bool {
	h, _, _ := procGetForeground.Call()
	if h == 0 || className(h) != "Windows.UI.Core.CoreWindow" {
		return false
	}
	var pid uint32
	procGetWindowPID.Call(h, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return false
	}
	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(proc)
	buf := make([]uint16, windows.MAX_PATH)
	n := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(proc, 0, &buf[0], &n); err != nil {
		return false
	}
	return startHosts[strings.ToLower(filepath.Base(windows.UTF16ToString(buf[:n])))]
}
