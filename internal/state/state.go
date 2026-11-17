// Package state decides which appearance a taskbar should have given what is
// currently happening on its monitor. It is pure logic with no OS calls.
package state

import "github.com/AnderCMD/TaskGlass/internal/config"

// Kind identifies which configuration rule produced an appearance.
type Kind string

const (
	KindDefault   Kind = "default"
	KindMaximized Kind = "maximized"
	KindStartMenu Kind = "start_menu"
)

// Context is the observed desktop state for one taskbar.
type Context struct {
	// StartMenu is true while the Start menu or Search has focus.
	StartMenu bool
	// Maximized is true when the taskbar's monitor has a maximized window.
	Maximized bool
}

// Resolve picks the appearance for ctx. Priority: Start menu, maximized
// window, default. Disabled rules are skipped.
func Resolve(c config.Config, ctx Context) (config.Appearance, Kind) {
	if ctx.StartMenu && c.StartMenu.Enabled {
		return c.StartMenu.Appearance, KindStartMenu
	}
	if ctx.Maximized && c.Maximized.Enabled {
		return c.Maximized.Appearance, KindMaximized
	}
	return c.Default, KindDefault
}
