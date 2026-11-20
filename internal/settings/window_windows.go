//go:build windows

package settings

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"golang.org/x/sys/windows/registry"
)

// edgePath finds Microsoft Edge (shipped with Windows 10/11).
func edgePath() string {
	for _, root := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		k, err := registry.OpenKey(root, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe`, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		v, _, err := k.GetStringValue("")
		k.Close()
		if err == nil {
			if _, err := os.Stat(v); err == nil {
				return v
			}
		}
	}
	for _, env := range []string{"ProgramFiles(x86)", "ProgramFiles"} {
		if base := os.Getenv(env); base != "" {
			p := filepath.Join(base, `Microsoft\Edge\Application\msedge.exe`)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
}

// Show opens url in a standalone app-style window (no tabs or address bar).
// It uses Edge with a private profile so it neither joins the user's browser
// session nor shows first-run prompts. If Edge is missing or cannot be started
// it tries the default browser, first through the shell and then through
// rundll32. When every way fails the error lists each attempt, so a security
// product that blocks TaskGlass from starting programs is visible in the log.
func Show(url, profileDir string) error {
	var errs []error
	if edge := edgePath(); edge != "" {
		cmd := exec.Command(edge,
			"--app="+url,
			"--user-data-dir="+profileDir,
			"--window-size=1240,820",
			"--no-first-run",
			"--no-default-browser-check",
			"--disable-extensions",
		)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		err := cmd.Start()
		if err == nil {
			go cmd.Wait()
			return nil
		}
		errs = append(errs, fmt.Errorf("edge: %w", err))
	}
	if err := windows.ShellExecute(0, windows.StringToUTF16Ptr("open"), windows.StringToUTF16Ptr(url), nil, nil, windows.SW_SHOWNORMAL); err == nil {
		return nil
	} else {
		errs = append(errs, fmt.Errorf("shell: %w", err))
	}
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url)
	if err := cmd.Start(); err != nil {
		return errors.Join(append(errs, fmt.Errorf("rundll32: %w", err))...)
	}
	go cmd.Wait()
	return nil
}

var (
	user32          = windows.NewLazySystemDLL("user32.dll")
	kernel32        = windows.NewLazySystemDLL("kernel32.dll")
	pOpenClipboard  = user32.NewProc("OpenClipboard")
	pCloseClipboard = user32.NewProc("CloseClipboard")
	pEmptyClipboard = user32.NewProc("EmptyClipboard")
	pSetClipboard   = user32.NewProc("SetClipboardData")
	pMessageBox     = user32.NewProc("MessageBoxW")
	pGlobalAlloc    = kernel32.NewProc("GlobalAlloc")
	pGlobalLock     = kernel32.NewProc("GlobalLock")
	pGlobalUnlock   = kernel32.NewProc("GlobalUnlock")
	pMoveMemory     = kernel32.NewProc("RtlMoveMemory")
)

// CopyToClipboard puts text on the clipboard. It works without starting any
// program, so it still works when something blocks TaskGlass from launching one.
func CopyToClipboard(text string) bool {
	u, err := windows.UTF16FromString(text)
	if err != nil {
		return false
	}
	if r, _, _ := pOpenClipboard.Call(0); r == 0 {
		return false
	}
	defer pCloseClipboard.Call()
	pEmptyClipboard.Call()
	size := uintptr(len(u) * 2)
	h, _, _ := pGlobalAlloc.Call(0x0002, size) // GMEM_MOVEABLE
	if h == 0 {
		return false
	}
	p, _, _ := pGlobalLock.Call(h)
	if p == 0 {
		return false
	}
	pMoveMemory.Call(p, uintptr(unsafe.Pointer(&u[0])), size)
	pGlobalUnlock.Call(h)
	r, _, _ := pSetClipboard.Call(13, h) // CF_UNICODETEXT
	return r != 0
}

// Notify shows a modal warning box (blocks until it is closed).
func Notify(title, text string) {
	t, _ := windows.UTF16PtrFromString(title)
	m, _ := windows.UTF16PtrFromString(text)
	pMessageBox.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x30|0x40000) // MB_ICONWARNING|MB_TOPMOST
}
