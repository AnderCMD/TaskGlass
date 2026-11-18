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
	procSendNotifyMsgW  = user32.NewProc("SendNotifyMessageW")
	procSetWindowsHook  = user32.NewProc("SetWindowsHookExW")
	procUnhookWindows   = user32.NewProc("UnhookWindowsHookEx")

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

// WindowThreadPID returns the thread and process that own hwnd (0, 0 if unknown).
func WindowThreadPID(hwnd HWND) (tid, pid uint32) {
	t, _, _ := procGetWindowPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return uint32(t), pid
}

const (
	whCallWndProc = 4
	wmNull        = 0
)

// SetCallWndProcHook installs a WH_CALLWNDPROC hook whose procedure lives in
// module mod, on thread tid. Windows loads mod into that thread's process the
// first time the hook fires; see NotifyNull.
func SetCallWndProcHook(proc uintptr, mod windows.Handle, tid uint32) (uintptr, error) {
	h, _, err := procSetWindowsHook.Call(whCallWndProc, proc, uintptr(mod), uintptr(tid))
	if h == 0 {
		return 0, err
	}
	return h, nil
}

// Unhook removes a hook installed with SetCallWndProcHook.
func Unhook(h uintptr) { procUnhookWindows.Call(h) }

// NotifyNull posts WM_NULL to hwnd without waiting, which is enough to make a
// WH_CALLWNDPROC hook on its thread fire.
func NotifyNull(hwnd HWND) { procSendNotifyMsgW.Call(hwnd, wmNull, 0, 0) }

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

// processExe returns the lower-cased executable file name of process pid, or "".
func processExe(pid uint32) string {
	if pid == 0 {
		return ""
	}
	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(proc)
	buf := make([]uint16, windows.MAX_PATH)
	n := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(proc, 0, &buf[0], &n); err != nil {
		return ""
	}
	return strings.ToLower(filepath.Base(windows.UTF16ToString(buf[:n])))
}

// Foreground returns the lower-cased executable name of the focused window and
// the monitor it is on. It returns ("", 0) when there is no usable foreground
// window or it belongs to the desktop shell.
func Foreground() (exe string, monitor uintptr) {
	h, _, _ := procGetForeground.Call()
	if h == 0 || shellClasses[className(h)] {
		return "", 0
	}
	var pid uint32
	procGetWindowPID.Call(h, uintptr(unsafe.Pointer(&pid)))
	return processExe(pid), Monitor(h)
}

// StartMenuActive reports whether the Start menu or Search currently has focus.
func StartMenuActive() bool {
	h, _, _ := procGetForeground.Call()
	if h == 0 || className(h) != "Windows.UI.Core.CoreWindow" {
		return false
	}
	var pid uint32
	procGetWindowPID.Call(h, uintptr(unsafe.Pointer(&pid)))
	return startHosts[processExe(pid)]
}

// FindProcessWindow returns a top-level window owned by a process whose
// executable file name is exe (lower case). It prefers the UWP CoreWindow, and
// otherwise any window the process owns, including hidden and cloaked ones.
func FindProcessWindow(exe string) (HWND, bool) {
	var core, any HWND
	enumTopLevel(func(h HWND) bool {
		var pid uint32
		procGetWindowPID.Call(h, uintptr(unsafe.Pointer(&pid)))
		if pid == 0 || processExe(pid) != exe {
			return true
		}
		if className(h) == "Windows.UI.Core.CoreWindow" {
			core = h
			return false
		}
		if any == 0 {
			any = h
		}
		return true
	})
	if core != 0 {
		return core, true
	}
	return any, any != 0
}

var (
	procGetGUIThreadInfo   = user32.NewProc("GetGUIThreadInfo")
	procPostThreadMessageW = user32.NewProc("PostThreadMessageW")
)

const whGetMessage = 3

// SetGetMessageHook installs a WH_GETMESSAGE hook whose procedure lives in
// module mod, on thread tid. Windows loads mod into that thread's process the
// first time the thread retrieves a message; see PostThreadNull.
func SetGetMessageHook(proc uintptr, mod windows.Handle, tid uint32) (uintptr, error) {
	h, _, err := procSetWindowsHook.Call(whGetMessage, proc, uintptr(mod), uintptr(tid))
	if h == 0 {
		return 0, err
	}
	return h, nil
}

// PostThreadNull posts WM_NULL to a thread's message queue, which makes a
// WH_GETMESSAGE hook on it fire.
func PostThreadNull(tid uint32) {
	procPostThreadMessageW.Call(uintptr(tid), wmNull, 0, 0)
}

// ProcessGUIThreads returns the process id and every GUI thread (a thread that
// has a message queue) of the process whose executable file name is exe, even
// when it owns no top-level windows. ok is false if the process is not running
// or has no GUI thread.
func ProcessGUIThreads(exe string) (pid uint32, tids []uint32, ok bool) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, nil, false
	}
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	for e := windows.Process32First(snap, &pe); e == nil; e = windows.Process32Next(snap, &pe) {
		if strings.EqualFold(windows.UTF16ToString(pe.ExeFile[:]), exe) {
			pid = pe.ProcessID
			break
		}
	}
	windows.CloseHandle(snap)
	if pid == 0 {
		return 0, nil, false
	}
	ts, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return pid, nil, false
	}
	defer windows.CloseHandle(ts)
	var te windows.ThreadEntry32
	te.Size = uint32(unsafe.Sizeof(te))
	type guiThreadInfo struct {
		Size                                               uint32
		Flags                                              uint32
		Active, Focus, Capture, MenuOwner, MoveSize, Caret uintptr
		Rc                                                 [4]int32
	}
	for e := windows.Thread32First(ts, &te); e == nil; e = windows.Thread32Next(ts, &te) {
		if te.OwnerProcessID != pid {
			continue
		}
		var gi guiThreadInfo
		gi.Size = uint32(unsafe.Sizeof(gi))
		// Succeeds only for threads that already have a message queue.
		if r, _, _ := procGetGUIThreadInfo.Call(uintptr(te.ThreadID), uintptr(unsafe.Pointer(&gi))); r != 0 {
			tids = append(tids, te.ThreadID)
		}
	}
	return pid, tids, len(tids) > 0
}

var (
	procGetWindowLongPtrW     = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW     = user32.NewProc("SetWindowLongPtrW")
	procSetLayeredWindowAttrs = user32.NewProc("SetLayeredWindowAttributes")
	procGetLayeredWindowAttrs = user32.NewProc("GetLayeredWindowAttributes")
)

const (
	wsExLayered = 0x00080000
	lwaAlpha    = 0x2
)

func exStyle(h HWND) uintptr {
	r, _, _ := procGetWindowLongPtrW.Call(h, ^uintptr(19)) // GWL_EXSTYLE = -20
	return r
}

// IsLayered reports whether hwnd has the WS_EX_LAYERED style.
func IsLayered(h HWND) bool { return exStyle(h)&wsExLayered != 0 }

// SetWindowAlpha makes the whole window uniformly translucent: alpha 255 is
// opaque, 0 invisible. It adds WS_EX_LAYERED if needed.
func SetWindowAlpha(h HWND, alpha uint8) error {
	if ex := exStyle(h); ex&wsExLayered == 0 {
		procSetWindowLongPtrW.Call(h, ^uintptr(19), ex|wsExLayered)
	}
	r, _, err := procSetLayeredWindowAttrs.Call(h, 0, uintptr(alpha), lwaAlpha)
	if r == 0 {
		return err
	}
	return nil
}

// ClearWindowAlpha makes the window opaque again. When removeStyle is true it
// also drops WS_EX_LAYERED (use it only if TaskGlass added the style).
func ClearWindowAlpha(h HWND, removeStyle bool) {
	if !IsLayered(h) {
		return
	}
	procSetLayeredWindowAttrs.Call(h, 0, 255, lwaAlpha)
	if removeStyle {
		procSetWindowLongPtrW.Call(h, ^uintptr(19), exStyle(h)&^wsExLayered)
	}
}

const (
	gwOwner         = 4
	wsExToolWindowF = 0x00000080
)

var procGetWindow = user32.NewProc("GetWindow")

type rect struct{ L, T, R, B int32 }

var procGetWindowRect = user32.NewProc("GetWindowRect")

// AppWindows returns the visible top-level application windows of processes
// whose executable file name is exe (lower case). Tool windows, owned windows
// (dialogs, popups), shell windows and tiny windows are skipped.
func AppWindows(exe string) []HWND {
	var out []HWND
	enumTopLevel(func(h HWND) bool {
		if v, _, _ := procIsWindowVisible.Call(h); v == 0 {
			return true
		}
		if c := className(h); shellClasses[c] || c == "CabinetWClass" {
			return true
		}
		if ex := exStyle(h); ex&wsExToolWindowF != 0 {
			return true
		}
		if o, _, _ := procGetWindow.Call(h, gwOwner); o != 0 {
			return true
		}
		var r rect
		if ok, _, _ := procGetWindowRect.Call(h, uintptr(unsafe.Pointer(&r))); ok == 0 || r.R-r.L < 120 || r.B-r.T < 60 {
			return true
		}
		var pid uint32
		procGetWindowPID.Call(h, uintptr(unsafe.Pointer(&pid)))
		if processExe(pid) == exe {
			out = append(out, h)
		}
		return true
	})
	return out
}
