//go:build windows

// Package theme reads the Windows light/dark colour mode.
package theme

import "golang.org/x/sys/windows/registry"

const personalizeKey = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

// SystemIsLight reports whether Windows uses the light theme for the taskbar
// and Start (the "system" mode, as opposed to the per-app mode). If the value
// cannot be read it assumes dark, the Windows 11 default.
func SystemIsLight() bool {
	return readBool("SystemUsesLightTheme")
}

// AppsAreLight reports whether Windows uses the light theme for apps.
func AppsAreLight() bool {
	return readBool("AppsUseLightTheme")
}

func readBool(name string) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, personalizeKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetIntegerValue(name)
	return err == nil && v != 0
}
