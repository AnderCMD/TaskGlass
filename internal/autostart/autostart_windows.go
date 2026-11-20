//go:build windows

// Package autostart toggles "start with Windows" via the per-user Run key.
package autostart

import (
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey      = `Software\Microsoft\Windows\CurrentVersion\Run`
	approvedKey = `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`
	valueName   = "TaskGlass"
)

// command is the Run-key value for exe: the path in plain double quotes.
//
// It must NOT use fmt's %q: that is Go-string syntax and doubles every
// backslash ("C:\\Dir\\app.exe"), which Windows does not resolve when it
// launches the entry at sign-in, so the app silently never started.
func command(exe string) string { return `"` + exe + `"` }

// Enabled reports whether autostart is registered for exe.
func Enabled(exe string) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(valueName)
	return err == nil && sameCommand(v, exe)
}

// sameCommand compares a stored value with exe, also recognising the
// doubled-backslash form written by TaskGlass 1.1.1 and earlier.
func sameCommand(v, exe string) bool {
	return v == command(exe) || v == `"`+strings.ReplaceAll(exe, `\`, `\\`)+`"`
}

// Set registers or removes autostart for exe.
func Set(exe string, on bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if on {
		if err := k.SetStringValue(valueName, command(exe)); err != nil {
			return err
		}
		clearDisabled()
		return nil
	}
	if err := k.DeleteValue(valueName); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}

// Repair rewrites an entry written by older versions (broken doubled
// backslashes) for exe. It does nothing when autostart is not registered or
// points at another program, and reports whether it changed anything.
func Repair(exe string) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(valueName)
	if err != nil || v == command(exe) || !sameCommand(v, exe) {
		return false // absent, already right, or deliberately pointing elsewhere
	}
	return k.SetStringValue(valueName, command(exe)) == nil
}

// clearDisabled drops the StartupApproved flag so a previously disabled entry
// (Task Manager > Startup apps) is enabled again when the user turns it on here.
func clearDisabled() {
	k, err := registry.OpenKey(registry.CURRENT_USER, approvedKey, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer k.Close()
	_ = k.DeleteValue(valueName)
}
