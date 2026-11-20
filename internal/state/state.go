// Package state decides which appearance a taskbar should have given what is
// currently happening on its monitor. It is pure logic with no OS calls.
package state

import "github.com/AnderCMD/TaskGlass/internal/config"

// Kind identifies which configuration rule produced an appearance.
type Kind string

const (
	KindDisabled  Kind = "disabled"
	KindDefault   Kind = "default"
	KindMaximized Kind = "maximized"
	KindApp       Kind = "app"
)

// Context is the observed desktop state for one taskbar.
type Context struct {
	// Maximized is true when the taskbar's monitor has a maximized window.
	Maximized bool
	// ForegroundExe is the executable file name of the focused window when
	// that window is on this taskbar's monitor, or "" otherwise.
	ForegroundExe string
}

// Resolve picks the appearance for ctx. Priority: master switch,
// per-application rule, maximized window, default. Disabled rules are skipped.
func Resolve(c config.Config, ctx Context) (config.Appearance, Kind) {
	if !c.Enabled {
		return config.Appearance{Accent: config.AccentNormal, Color: "#000000"}, KindDisabled
	}
	if ctx.ForegroundExe != "" {
		for _, r := range c.Apps {
			if r.Matches(ctx.ForegroundExe) {
				return r.Appearance, KindApp
			}
		}
	}
	if ctx.Maximized && c.Maximized.Enabled {
		return c.Maximized.Appearance, KindMaximized
	}
	return c.Default, KindDefault
}

// NeedsWindowScan reports whether the configuration depends on knowing about
// maximized windows or the focused application, so idle setups can skip the
// (cheap, but not free) window enumeration entirely.
func NeedsWindowScan(c config.Config) (maximized, foreground bool) {
	if !c.Enabled {
		return false, false
	}
	flyout := c.Surfaces.TrayOverflow.Enabled && c.Surfaces.TrayOverflow.Maximized.Enabled
	return c.Maximized.Enabled || flyout, len(c.Apps) > 0
}
