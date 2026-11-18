//go:build windows

package winapi

import (
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shell32              = windows.NewLazySystemDLL("shell32.dll")
	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")
	procLoadIconW        = user32.NewProc("LoadIconW")
)

// notifyIconData mirrors NOTIFYICONDATAW (Windows Vista and later layout).
type notifyIconData struct {
	CbSize            uint32
	HWnd              uintptr
	UID               uint32
	UFlags            uint32
	UCallbackMessage  uint32
	HIcon             uintptr
	SzTip             [128]uint16
	DwState           uint32
	DwStateMask       uint32
	SzInfo            [256]uint16
	UTimeoutOrVersion uint32
	SzInfoTitle       [64]uint16
	DwInfoFlags       uint32
	GuidItem          windows.GUID
	HBalloonIcon      uintptr
}

// TrayProbe adds a tray icon and removes it again after hold. Explorer builds
// and tears down a tray item for it, which makes the shell's XAML trees (the
// tray, including the overflow flyout) emit the live events that TaskGlass'
// helper DLL needs in order to find them. The icon is added and removed in a
// fraction of a second.
func TrayProbe(hold time.Duration) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	const (
		hwndMessage = ^uintptr(2) // HWND_MESSAGE (-3)
		nifIcon     = 0x2
		nifTip      = 0x4
		nimAdd      = 0
		nimDelete   = 2
		idiApp      = 32512
	)
	inst, _, _ := procGetModuleHandleW.Call(0)
	cls, _ := windows.UTF16PtrFromString("TaskGlassTrayProbe")
	wc := wndClassEx{Size: uint32(unsafe.Sizeof(wndClassEx{})), WndProc: anchorProc, Instance: inst, ClassName: cls}
	if r, _, e := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 && e != windows.ERROR_CLASS_ALREADY_EXISTS {
		return e
	}
	title, _ := windows.UTF16PtrFromString("TaskGlass")
	h, _, e := procCreateWindowExW.Call(0, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)), 0, 0, 0, 0, 0, hwndMessage, 0, inst, 0)
	if h == 0 {
		return e
	}
	defer procDestroyWindow.Call(h)

	icon, _, _ := procLoadIconW.Call(0, idiApp)
	var nid notifyIconData
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = h
	nid.UID = 0x7461 // arbitrary, only needs to be unique within our window
	nid.UFlags = nifIcon | nifTip
	nid.HIcon = icon
	copy(nid.SzTip[:], windows.StringToUTF16("TaskGlass"))
	if r, _, e := procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&nid))); r == 0 {
		return e
	}
	time.Sleep(hold)
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&nid)))
	return nil
}
