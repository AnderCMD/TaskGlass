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

// Appearance describes how the taskbar should look in a given state.
type Appearance struct {
	Accent  Accent `json:"accent"`
	Color   string `json:"color"`   // "#RRGGBB"
	Opacity int    `json:"opacity"` // 0 (invisible) .. 100 (fully opaque)
}

// Rule is an Appearance that only applies when Enabled is true.
type Rule struct {
	Enabled bool `json:"enabled"`
	Appearance
}

// Config is the full TaskGlass configuration.
type Config struct {
	Version int `json:"version"`
	// Default is used when no other rule matches.
	Default Appearance `json:"default"`
	// Maximized applies on a monitor that has a maximized window.
	Maximized Rule `json:"maximized"`
	// StartMenu applies while the Start menu / Search is focused.
	StartMenu Rule `json:"start_menu"`
	// PollMs is how often window state is sampled.
	PollMs int `json:"poll_ms"`
	// RefreshMs is how often the effect is re-applied even if nothing changed,
	// to survive Explorer silently resetting it.
	RefreshMs int `json:"refresh_ms"`
	// LogLevel is one of debug, info, warn, error.
	LogLevel string `json:"log_level"`
	// XamlFix removes Windows 11's own XAML taskbar fill (needs
	// taskglass_tap.dll next to the executable). Without it the taskbar can
	// never become fully transparent on Windows 11.
	XamlFix bool `json:"xaml_fix"`
}

// Default returns the built-in configuration.
func Default() Config {
	return Config{
		Version: CurrentVersion,
		Default: Appearance{Accent: AccentClear, Color: "#000000", Opacity: 0},
		Maximized: Rule{
			Enabled:    true,
			Appearance: Appearance{Accent: AccentOpaque, Color: "#000000", Opacity: 100},
		},
		StartMenu: Rule{
			Enabled:    false,
			Appearance: Appearance{Accent: AccentAcrylic, Color: "#202020", Opacity: 60},
		},
		PollMs:    100,
		RefreshMs: 1000,
		LogLevel:  "info",
		XamlFix:   true,
	}
}

// Validate checks the configuration and returns a descriptive error for the
// first problem found.
func (c Config) Validate() error {
	if c.Version != CurrentVersion {
		return fmt.Errorf("unsupported config version %d (want %d)", c.Version, CurrentVersion)
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
	if c.PollMs < 16 || c.PollMs > 5000 {
		return fmt.Errorf("poll_ms %d out of range [16, 5000]", c.PollMs)
	}
	if c.RefreshMs < c.PollMs || c.RefreshMs > 60000 {
		return fmt.Errorf("refresh_ms %d out of range [poll_ms, 60000]", c.RefreshMs)
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
	if _, _, _, err := ParseColor(a.Color); err != nil {
		return err
	}
	if a.Opacity < 0 || a.Opacity > 100 {
		return fmt.Errorf("opacity %d out of range [0, 100]", a.Opacity)
	}
	return nil
}

// ParseColor parses "#RRGGBB" (the leading '#' is optional).
func ParseColor(s string) (r, g, b uint8, err error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid colour %q: want #RRGGBB", s)
	}
	v, perr := strconv.ParseUint(h, 16, 32)
	if perr != nil {
		return 0, 0, 0, fmt.Errorf("invalid colour %q: want #RRGGBB", s)
	}
	return uint8(v >> 16), uint8(v >> 8), uint8(v), nil
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
