//go:build windows

package winapi

import (
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDestroyWindow    = user32.NewProc("DestroyWindow")
	procShowWindow       = user32.NewProc("ShowWindow")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procPeekMessageW     = user32.NewProc("PeekMessageW")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procGetModuleHandleW = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetModuleHandleW")
)

type wndClassEx struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

type point struct{ X, Y int32 }
type msg struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

var anchorProc = syscall.NewCallback(func(hwnd, m, wp, lp uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(hwnd, m, wp, lp)
	return r
})

// TaskbarAnchor makes the taskbar's XAML tree emit live "element added/removed"
// events by briefly showing an off-screen, non-activating window that owns a
// taskbar button (hold is how long it stays). It blocks for about hold.
func TaskbarAnchor(hold time.Duration) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	const (
		wsPopup        = 0x80000000
		wsExAppWindow  = 0x00040000
		wsExNoActivate = 0x08000000
		swShowNoActive = 4
		pmRemove       = 1
		className      = "TaskGlassAnchor"
	)
	var off int32 = -32000 // far off-screen, sign-extended into the call below
	inst, _, _ := procGetModuleHandleW.Call(0)
	cls, _ := windows.UTF16PtrFromString(className)
	wc := wndClassEx{Size: uint32(unsafe.Sizeof(wndClassEx{})), WndProc: anchorProc, Instance: inst, ClassName: cls}
	if r, _, e := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 && e != windows.ERROR_CLASS_ALREADY_EXISTS {
		return e
	}
	title, _ := windows.UTF16PtrFromString("TaskGlass")
	h, _, e := procCreateWindowExW.Call(wsExAppWindow|wsExNoActivate, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)),
		wsPopup, uintptr(off), uintptr(off), 1, 1, 0, 0, inst, 0)
	if h == 0 {
		return e
	}
	defer procDestroyWindow.Call(h)
	procShowWindow.Call(h, swShowNoActive)

	var m msg
	for end := time.Now().Add(hold); time.Now().Before(end); time.Sleep(10 * time.Millisecond) {
		for {
			r, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, pmRemove)
			if r == 0 {
				break
			}
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		}
	}
	return nil
}
