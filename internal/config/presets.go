package config

import (
	"fmt"
	"sort"
	"strings"
)

// Preset is a named look. Applying one replaces only the visual settings
// (taskbar appearances, rules for maximized windows and the Start menu, and
// the glass surfaces); your app rules, hotkey, theme and everything else stay.
type Preset struct {
	Name        string
	Description string
	apply       func(c *Config)
}

func app(accent Accent, opacity Percent) Appearance {
	return Appearance{Accent: accent, Color: ColorAuto, Opacity: opacity}
}

var presets = []Preset{
	{
		Name:        "liquid",
		Description: "Frosted taskbar over the desktop, translucent menus, flyouts and Explorer windows (closest to macOS-style glass)",
		apply: func(c *Config) {
			c.Enabled = true
			c.Default = app(AccentAcrylic, 28)
			c.Maximized = Rule{Enabled: true, Appearance: app(AccentAcrylic, 72)}
			c.StartMenu = Rule{Enabled: true, Appearance: app(AccentAcrylic, 45)}
			c.Surfaces = Surfaces{
				TrayOverflow: Surface{Enabled: true, Opacity: 55},
				StartMenu:    Surface{Enabled: true, Opacity: 55},
				Explorer:     Surface{Enabled: true, Opacity: 88},
			}
			c.XamlFix = true
		},
	},
	{
		Name:        "clear",
		Description: "Fully transparent taskbar, flyouts and Start menu translucent, Explorer untouched",
		apply: func(c *Config) {
			c.Enabled = true
			c.Default = app(AccentClear, 0)
			c.Maximized = Rule{Enabled: true, Appearance: app(AccentOpaque, 100)}
			c.StartMenu = Rule{Enabled: false, Appearance: app(AccentAcrylic, 60)}
			c.Surfaces = Surfaces{
				TrayOverflow: Surface{Enabled: true, Opacity: 60},
				StartMenu:    Surface{Enabled: true, Opacity: 60},
				Explorer:     Surface{Enabled: false, Opacity: 85},
			}
			c.XamlFix = true
		},
	},
	{
		Name:        "subtle",
		Description: "Light blur on the taskbar and a barely translucent shell",
		apply: func(c *Config) {
			c.Enabled = true
			c.Default = app(AccentBlur, 18)
			c.Maximized = Rule{Enabled: true, Appearance: app(AccentOpaque, 100)}
			c.StartMenu = Rule{Enabled: false, Appearance: app(AccentAcrylic, 60)}
			c.Surfaces = Surfaces{
				TrayOverflow: Surface{Enabled: true, Opacity: 85},
				StartMenu:    Surface{Enabled: true, Opacity: 88},
				Explorer:     Surface{Enabled: true, Opacity: 96},
			}
			c.XamlFix = true
		},
	},
	{
		Name:        "stock",
		Description: "Windows as it ships: nothing is changed",
		apply: func(c *Config) {
			c.Enabled = false
		},
	},
}

// Presets returns the built-in presets in menu order.
func Presets() []Preset { return append([]Preset(nil), presets...) }

// PresetNames lists the preset names, sorted.
func PresetNames() []string {
	names := make([]string, 0, len(presets))
	for _, p := range presets {
		names = append(names, p.Name)
	}
	sort.Strings(names)
	return names
}

// ApplyPreset returns c with the named preset applied.
func ApplyPreset(c Config, name string) (Config, error) {
	for _, p := range presets {
		if strings.EqualFold(p.Name, strings.TrimSpace(name)) {
			next := c
			next.Apps = append([]AppRule(nil), c.Apps...)
			p.apply(&next)
			if err := next.Validate(); err != nil {
				return c, fmt.Errorf("preset %s: %w", p.Name, err)
			}
			return next, nil
		}
	}
	return c, fmt.Errorf("unknown preset %q (choose from %s)", name, strings.Join(PresetNames(), ", "))
}
