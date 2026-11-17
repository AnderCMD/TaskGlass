//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

// attachConsole connects stdout/stderr to the console that launched us.
//
// taskglass.exe is a GUI-subsystem program (so double-clicking it never
// flashes a console window), which means a shell does not give it a console.
// Command-line flags such as --validate still need to print, so when stdout or
// stderr are not already redirected we attach to the parent's console.
func attachConsole() {
	need := func(std uint32) bool {
		h, err := windows.GetStdHandle(std)
		return err != nil || h == 0 || h == windows.InvalidHandle
	}
	if !need(windows.STD_OUTPUT_HANDLE) && !need(windows.STD_ERROR_HANDLE) {
		return // already redirected to a file or pipe
	}
	const attachParentProcess = ^uint32(0)
	if err := windows.NewLazySystemDLL("kernel32.dll").NewProc("AttachConsole").Find(); err != nil {
		return
	}
	r, _, _ := windows.NewLazySystemDLL("kernel32.dll").NewProc("AttachConsole").Call(uintptr(attachParentProcess))
	if r == 0 {
		return
	}
	if f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		if need(windows.STD_OUTPUT_HANDLE) {
			os.Stdout = f
		}
		if need(windows.STD_ERROR_HANDLE) {
			os.Stderr = f
		}
	}
}
