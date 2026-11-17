//go:build windows

// Package autostart toggles "start with Windows" via the per-user Run key.
package autostart

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey    = `Software\Microsoft\Windows\CurrentVersion\Run`
	valueName = "TaskGlass"
)

func command(exe string) string { return fmt.Sprintf("%q", exe) }

// Enabled reports whether autostart is registered for exe.
func Enabled(exe string) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(valueName)
	return err == nil && v == command(exe)
}

// Set registers or removes autostart for exe.
func Set(exe string, on bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if on {
		return k.SetStringValue(valueName, command(exe))
	}
	if err := k.DeleteValue(valueName); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	return nil
}
