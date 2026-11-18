// Package config defines TaskGlass' on-disk configuration: its schema,
// defaults, validation and atomic persistence. It has no OS dependencies so it
// can be unit-tested anywhere.
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CurrentVersion is the schema version written to new config files.
const CurrentVersion = 1

// SchemaURL is written to new config files so editors offer completion.
const SchemaURL = "https://raw.githubusercontent.com/AnderCMD/TaskGlass/main/docs/config.schema.json"

// MaxAppRules bounds the number of per-application rules.
const MaxAppRules = 64

// Accent is the visual effect applied to the taskbar.
type Accent string

const (
	// AccentNormal leaves the taskbar untouched (stock Windows look).
	AccentNormal Accent = "normal"
	// AccentClear makes the taskbar background fully transparent.
	AccentClear Accent = "clear"
	// AccentBlur applies a gaussian blur behind a tinted background.
	AccentBlur Accent = "blur"
	// AccentAcrylic applies the acrylic (blur + noise) effect.
	AccentAcrylic Accent = "acrylic"
	// AccentOpaque paints a solid colour.
	AccentOpaque Accent = "opaque"
)

// Accents lists every valid accent, in menu order.
var Accents = []Accent{AccentClear, AccentBlur, AccentAcrylic, AccentOpaque, AccentNormal}

// Valid reports whether a is a known accent.
func (a Accent) Valid() bool {
	for _, k := range Accents {
		if a == k {
			return true
		}
	}
	return false
}

// Theme selects which Windows colour mode drives "auto" colours and the tray icon.
type Theme string

const (
	ThemeAuto  Theme = "auto"  // follow Windows (system theme)
	ThemeLight Theme = "light" // always treat the system as light
	ThemeDark  Theme = "dark"  // always treat the system as dark
)

// Themes lists every valid theme, in menu order.
var Themes = []Theme{ThemeAuto, ThemeLight, ThemeDark}

// Valid reports whether t is a known theme.
func (t Theme) Valid() bool {
	for _, k := range Themes {
		if t == k {
			return true
		}
	}
	return false
}

// Light reports whether the system should be treated as light, given what
// Windows currently reports.
func (t Theme) Light(systemIsLight bool) bool {
	switch t {
	case ThemeLight:
		return true
	case ThemeDark:
		return false
	}
	return systemIsLight
}

// TrayIcon selects the notification-area icon style.
type TrayIcon string

const (
	TrayAuto  TrayIcon = "auto"  // monochrome glyph matching the Windows theme
	TrayColor TrayIcon = "color" // full-colour logo
	TrayLight TrayIcon = "light" // always the white glyph
	TrayDark  TrayIcon = "dark"  // always the dark glyph
)

// Valid reports whether t is a known tray icon style.
func (t TrayIcon) Valid() bool {
	switch t {
	case TrayAuto, TrayColor, TrayLight, TrayDark:
		return true
	}
	return false
}

// Taskbars selects which taskbars are styled.
type Taskbars string

const (
	TaskbarsAll     Taskbars = "all"
	TaskbarsPrimary Taskbars = "primary"
)

// ColorAuto is the colour value that follows the Windows theme.
const ColorAuto = "auto"

// Colours used when a colour is "auto": they match the stock Windows 11 taskbar.
const (
	autoDark  = "#1C1C1C"
	autoLight = "#F3F3F3"
)

// Appearance describes how the taskbar should look in a given state.
type Appearance struct {
	Accent  Accent `json:"accent"`
	Color   string `json:"color"`   // "#RRGGBB" or "auto"
	Opacity int    `json:"opacity"` // 0 (invisible) .. 100 (fully opaque)
}

// RGB returns the tint colour for the given theme.
func (a Appearance) RGB(light bool) (r, g, b uint8) {
	c := a.Color
	if strings.EqualFold(strings.TrimSpace(c), ColorAuto) {
		c = autoDark
		if light {
			c = autoLight
		}
	}
	r, g, b, _ = ParseColor(c)
	return r, g, b
}

// Rule is an Appearance that only applies when Enabled is true.
type Rule struct {
	Enabled bool `json:"enabled"`
	Appearance
}

// AppRule styles the taskbar while a specific application has focus.
type AppRule struct {
	// Exe is the executable file name, matched case-insensitively
	// (for example "chrome.exe").
	Exe string `json:"exe"`
	Appearance
}

// Config is the full TaskGlass configuration.
type Config struct {
	Schema  string `json:"$schema,omitempty"`
	Version int    `json:"version"`
	// Enabled is the master switch. When false the stock taskbar is shown.
	Enabled bool `json:"enabled"`
	// Theme drives "auto" colours and the tray icon.
	Theme Theme `json:"theme"`
	// Taskbars selects which taskbars are styled: all monitors or the primary.
	Taskbars Taskbars `json:"taskbars"`
	// Default is used when no other rule matches.
	Default Appearance `json:"default"`
	// Maximized applies on a monitor that has a maximized window.
	Maximized Rule `json:"maximized"`
	// StartMenu applies while the Start menu / Search is focused.
	StartMenu Rule `json:"start_menu"`
	// Apps apply while the named application is focused (first match wins).
	Apps []AppRule `json:"apps"`
	// TrayIcon is the notification-area icon style.
	TrayIcon TrayIcon `json:"tray_icon"`
	// HotkeyEnabled turns the global hotkey on. It is off by default so
	// TaskGlass never grabs a shortcut you did not ask for.
	HotkeyEnabled bool `json:"hotkey_enabled"`
	// Hotkey is the shortcut that toggles Enabled when HotkeyEnabled is true,
	// e.g. "Ctrl+Alt+G".
	Hotkey string `json:"hotkey"`
	// PollMs is the fallback sampling interval. Window events normally trigger
	// updates instantly; sampling only guards against missed events.
	PollMs int `json:"poll_ms"`
	// DebounceMs coalesces bursts of window events into one update.
	DebounceMs int `json:"debounce_ms"`
	// RefreshMs is how often the effect is re-applied even if nothing changed,
	// to survive Explorer silently resetting it.
	RefreshMs int `json:"refresh_ms"`
	// XamlFix removes Windows 11's own XAML taskbar fill (needs
	// taskglass_tap.dll next to the executable). Without it the taskbar can
	// never become fully transparent on Windows 11.
	XamlFix bool `json:"xaml_fix"`
	// LogLevel is one of debug, info, warn, error.
	LogLevel string `json:"log_level"`
}

// Default returns the built-in configuration.
func Default() Config {
	return Config{
		Schema:   SchemaURL,
		Version:  CurrentVersion,
		Enabled:  true,
		Theme:    ThemeAuto,
		Taskbars: TaskbarsAll,
		Default:  Appearance{Accent: AccentClear, Color: ColorAuto, Opacity: 0},
		Maximized: Rule{
			Enabled:    true,
			Appearance: Appearance{Accent: AccentOpaque, Color: ColorAuto, Opacity: 100},
		},
		StartMenu: Rule{
			Enabled:    false,
			Appearance: Appearance{Accent: AccentAcrylic, Color: ColorAuto, Opacity: 60},
		},
		Apps:     []AppRule{},
		TrayIcon: TrayAuto,
		// Off by default: enabling it is one switch, with a sensible shortcut ready.
		HotkeyEnabled: false,
		Hotkey:        "Ctrl+Alt+G",
		PollMs:        500,
		DebounceMs:    25,
		RefreshMs:     1000,
		XamlFix:       true,
		LogLevel:      "info",
	}
}

// Validate checks the configuration and returns a descriptive error for the
// first problem found.
func (c Config) Validate() error {
	if c.Version != CurrentVersion {
		return fmt.Errorf("unsupported config version %d (want %d)", c.Version, CurrentVersion)
	}
	if !c.Theme.Valid() {
		return fmt.Errorf("theme %q must be auto, light or dark", c.Theme)
	}
	if c.Taskbars != TaskbarsAll && c.Taskbars != TaskbarsPrimary {
		return fmt.Errorf("taskbars %q must be all or primary", c.Taskbars)
	}
	if !c.TrayIcon.Valid() {
		return fmt.Errorf("tray_icon %q must be auto, color, light or dark", c.TrayIcon)
	}
	checks := []struct {
		name string
		a    Appearance
	}{
		{"default", c.Default},
		{"maximized", c.Maximized.Appearance},
		{"start_menu", c.StartMenu.Appearance},
	}
	for _, k := range checks {
		if err := k.a.Validate(); err != nil {
			return fmt.Errorf("%s: %w", k.name, err)
		}
	}
	if len(c.Apps) > MaxAppRules {
		return fmt.Errorf("apps: at most %d rules (got %d)", MaxAppRules, len(c.Apps))
	}
	for i, r := range c.Apps {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("apps[%d]: %w", i, err)
		}
	}
	if c.Hotkey != "" {
		if _, _, err := ParseHotkey(c.Hotkey); err != nil {
			return fmt.Errorf("hotkey: %w", err)
		}
	} else if c.HotkeyEnabled {
		return errors.New("hotkey_enabled is true but hotkey is empty")
	}
	if c.PollMs < 50 || c.PollMs > 10000 {
		return fmt.Errorf("poll_ms %d out of range [50, 10000]", c.PollMs)
	}
	if c.DebounceMs < 0 || c.DebounceMs > 1000 {
		return fmt.Errorf("debounce_ms %d out of range [0, 1000]", c.DebounceMs)
	}
	if c.RefreshMs < 100 || c.RefreshMs > 60000 {
		return fmt.Errorf("refresh_ms %d out of range [100, 60000]", c.RefreshMs)
	}
	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log_level %q must be debug, info, warn or error", c.LogLevel)
	}
	return nil
}

// Validate checks a single appearance.
func (a Appearance) Validate() error {
	if !a.Accent.Valid() {
		return fmt.Errorf("unknown accent %q", a.Accent)
	}
	if !strings.EqualFold(strings.TrimSpace(a.Color), ColorAuto) {
		if _, _, _, err := ParseColor(a.Color); err != nil {
			return err
		}
	}
	if a.Opacity < 0 || a.Opacity > 100 {
		return fmt.Errorf("opacity %d out of range [0, 100]", a.Opacity)
	}
	return nil
}

// Validate checks an application rule.
func (r AppRule) Validate() error {
	exe := strings.TrimSpace(r.Exe)
	if exe == "" {
		return errors.New("exe is required")
	}
	if strings.ContainsAny(exe, `\/:*?"<>|`) {
		return fmt.Errorf("exe %q must be a bare file name like chrome.exe", r.Exe)
	}
	return r.Appearance.Validate()
}

// Matches reports whether the rule applies to the executable file name exe.
func (r AppRule) Matches(exe string) bool {
	return strings.EqualFold(strings.TrimSpace(r.Exe), exe)
}

// ParseColor parses "#RRGGBB" (the leading '#' is optional).
func ParseColor(s string) (r, g, b uint8, err error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid colour %q: want #RRGGBB or auto", s)
	}
	v, perr := strconv.ParseUint(h, 16, 32)
	if perr != nil {
		return 0, 0, 0, fmt.Errorf("invalid colour %q: want #RRGGBB or auto", s)
	}
	return uint8(v >> 16), uint8(v >> 8), uint8(v), nil
}

// Hotkey modifier flags, identical to the Win32 MOD_* values.
const (
	ModAlt   = 0x1
	ModCtrl  = 0x2
	ModShift = 0x4
	ModWin   = 0x8
)

// ParseHotkey parses strings such as "Ctrl+Alt+G" or "Win+Shift+F12" into
// Win32 modifier flags and a virtual-key code. At least one modifier is
// required so ordinary typing is never intercepted.
func ParseHotkey(s string) (mods uint32, vk uint32, err error) {
	parts := strings.Split(s, "+")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return 0, 0, fmt.Errorf("invalid hotkey %q", s)
		}
		if i < len(parts)-1 {
			switch strings.ToLower(p) {
			case "ctrl", "control":
				mods |= ModCtrl
			case "alt":
				mods |= ModAlt
			case "shift":
				mods |= ModShift
			case "win", "windows", "super":
				mods |= ModWin
			default:
				return 0, 0, fmt.Errorf("unknown modifier %q in hotkey %q", p, s)
			}
			continue
		}
		up := strings.ToUpper(p)
		switch {
		case len(up) == 1 && ((up[0] >= 'A' && up[0] <= 'Z') || (up[0] >= '0' && up[0] <= '9')):
			vk = uint32(up[0]) // VK codes for A-Z and 0-9 equal their ASCII values
		case len(up) >= 2 && len(up) <= 3 && up[0] == 'F':
			n, perr := strconv.Atoi(up[1:])
			if perr != nil || n < 1 || n > 24 {
				return 0, 0, fmt.Errorf("invalid function key %q in hotkey %q", p, s)
			}
			vk = uint32(0x70 + n - 1) // VK_F1 = 0x70
		default:
			return 0, 0, fmt.Errorf("unsupported key %q in hotkey %q (use A-Z, 0-9 or F1-F24)", p, s)
		}
	}
	if mods == 0 {
		return 0, 0, fmt.Errorf("hotkey %q needs at least one of Ctrl, Alt, Shift, Win", s)
	}
	return mods, vk, nil
}

// Parse decodes and validates a config document. Unknown fields are rejected
// so typos are reported instead of silently ignored.
func Parse(data []byte) (Config, error) {
	c := Default()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if err := c.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}
	return c, nil
}

// Load reads the config at path. If the file does not exist it is created with
// defaults. A malformed file returns an error and leaves the file untouched.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		c := Default()
		if serr := Save(path, c); serr != nil {
			return c, serr
		}
		return c, nil
	}
	if err != nil {
		return Config{}, err
	}
	return Parse(data)
}

// Save writes c to path atomically (temp file + rename).
func Save(path string, c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
