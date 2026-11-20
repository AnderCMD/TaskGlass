//go:build windows

package settings

import "github.com/AnderCMD/TaskGlass/internal/i18n"

func systemLocale() string { return i18n.SystemLocale() }
