package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultIsValid(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestParseColor(t *testing.T) {
	r, g, b, err := ParseColor("#1A2b3C")
	if err != nil || r != 0x1A || g != 0x2B || b != 0x3C {
		t.Fatalf("got %d %d %d %v", r, g, b, err)
	}
	for _, bad := range []string{"", "#12345", "#GGGGGG", "red", "#1234567"} {
		if _, _, _, err := ParseColor(bad); err == nil {
			t.Errorf("ParseColor(%q) should fail", bad)
		}
	}
}

func TestAutoColorFollowsTheme(t *testing.T) {
	a := Appearance{Accent: AccentOpaque, Color: "auto", Opacity: 100}
	if r, g, b := a.RGB(false); r != 0x1C || g != 0x1C || b != 0x1C {
		t.Errorf("dark auto = %d %d %d", r, g, b)
	}
	if r, g, b := a.RGB(true); r != 0xF3 || g != 0xF3 || b != 0xF3 {
		t.Errorf("light auto = %d %d %d", r, g, b)
	}
	fixed := Appearance{Accent: AccentOpaque, Color: "#102030", Opacity: 100}
	if r, g, b := fixed.RGB(true); r != 0x10 || g != 0x20 || b != 0x30 {
		t.Errorf("fixed colour must ignore the theme, got %d %d %d", r, g, b)
	}
}

func TestThemeLight(t *testing.T) {
	cases := []struct {
		theme  Theme
		system bool
		want   bool
	}{
		{ThemeAuto, true, true}, {ThemeAuto, false, false},
		{ThemeLight, false, true}, {ThemeDark, true, false},
	}
	for _, c := range cases {
		if got := c.theme.Light(c.system); got != c.want {
			t.Errorf("%s.Light(%v) = %v, want %v", c.theme, c.system, got, c.want)
		}
	}
}

func TestParseHotkey(t *testing.T) {
	good := map[string]struct{ mods, vk uint32 }{
		"Ctrl+Alt+G":        {ModCtrl | ModAlt, 'G'},
		"win+shift+f12":     {ModWin | ModShift, 0x7B},
		" Control + 5 ":     {ModCtrl, '5'},
		"Ctrl+Alt+Shift+F1": {ModCtrl | ModAlt | ModShift, 0x70},
	}
	for s, want := range good {
		mods, vk, err := ParseHotkey(s)
		if err != nil || mods != want.mods || vk != want.vk {
			t.Errorf("ParseHotkey(%q) = %#x %#x %v, want %#x %#x", s, mods, vk, err, want.mods, want.vk)
		}
	}
	for _, bad := range []string{"G", "Ctrl+", "Ctrl+Alt", "Ctrl+F25", "Ctrl+Foo+G", "Ctrl++G", "Ctrl+Enter", "Ctrl+F0"} {
		if _, _, err := ParseHotkey(bad); err == nil {
			t.Errorf("ParseHotkey(%q) should fail", bad)
		}
	}
}

func TestParseRejectsBadInput(t *testing.T) {
	cases := map[string]string{
		"unknown field":    `{"version":1,"nope":1}`,
		"bad accent":       `{"version":1,"default":{"accent":"glow","color":"#000000","opacity":0}}`,
		"bad opacity":      `{"version":1,"default":{"accent":"clear","color":"#000000","opacity":101}}`,
		"bad version":      `{"version":2}`,
		"poll too low":     `{"version":1,"poll_ms":1}`,
		"refresh too low":  `{"version":1,"refresh_ms":10}`,
		"debounce too big": `{"version":1,"debounce_ms":5000}`,
		"bad log level":    `{"version":1,"log_level":"loud"}`,
		"malformed json":   `{`,
		"bad nested rule":  `{"version":1,"maximized":{"enabled":true,"accent":"clear","color":"nope","opacity":1}}`,
		"bad theme":        `{"version":1,"theme":"sepia"}`,
		"bad taskbars":     `{"version":1,"taskbars":"some"}`,
		"bad tray icon":    `{"version":1,"tray_icon":"rainbow"}`,
		"bad hotkey":       `{"version":1,"hotkey":"G"}`,
		"app without exe":  `{"version":1,"apps":[{"exe":"","accent":"clear","color":"auto","opacity":0}]}`,
		"app with path":    `{"version":1,"apps":[{"exe":"C:\\x\\a.exe","accent":"clear","color":"auto","opacity":0}]}`,
		"app bad accent":   `{"version":1,"apps":[{"exe":"a.exe","accent":"x","color":"auto","opacity":0}]}`,
	}
	for name, doc := range cases {
		if _, err := Parse([]byte(doc)); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestTooManyAppRules(t *testing.T) {
	c := Default()
	for i := 0; i <= MaxAppRules; i++ {
		c.Apps = append(c.Apps, AppRule{Exe: "a.exe", Appearance: c.Default})
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected an error for too many app rules")
	}
}

func TestAppRuleMatching(t *testing.T) {
	r := AppRule{Exe: " Chrome.EXE "}
	if !r.Matches("chrome.exe") || r.Matches("firefox.exe") {
		t.Fatal("matching must be case-insensitive and exact")
	}
}

func TestParsePartialKeepsDefaults(t *testing.T) {
	c, err := Parse([]byte(`{"version":1,"poll_ms":200,"refresh_ms":2000}`))
	if err != nil {
		t.Fatal(err)
	}
	if c.PollMs != 200 || c.Default.Accent != AccentClear || !c.Enabled || !c.XamlFix || c.Theme != ThemeAuto {
		t.Fatalf("unexpected config: %+v", c)
	}
}

func TestHotkeyIsOffByDefaultButReadyToEnable(t *testing.T) {
	d := Default()
	if d.HotkeyEnabled {
		t.Fatal("the hotkey must be disabled by default")
	}
	if d.Hotkey == "" {
		t.Fatal("a ready-made shortcut must be preconfigured so enabling is one switch")
	}
	if _, _, err := ParseHotkey(d.Hotkey); err != nil {
		t.Fatalf("default hotkey must be valid: %v", err)
	}
	c, err := Parse([]byte(`{"version":1}`))
	if err != nil || c.HotkeyEnabled {
		t.Fatalf("a config without the key must keep it off: %+v %v", c, err)
	}
	c, err = Parse([]byte(`{"version":1,"hotkey_enabled":true}`))
	if err != nil || !c.HotkeyEnabled || c.Hotkey != "Ctrl+Alt+G" {
		t.Fatalf("enabling alone must work with the default shortcut: %+v %v", c, err)
	}
}

func TestHotkeyEnabledNeedsAShortcut(t *testing.T) {
	if _, err := Parse([]byte(`{"version":1,"hotkey_enabled":true,"hotkey":""}`)); err == nil {
		t.Fatal("hotkey_enabled with an empty hotkey must be rejected")
	}
	// disabled with an empty shortcut is fine
	if _, err := Parse([]byte(`{"version":1,"hotkey_enabled":false,"hotkey":""}`)); err != nil {
		t.Fatalf("disabled + empty must be accepted: %v", err)
	}
	// a bad shortcut is reported even while disabled, so typos surface early
	if _, err := Parse([]byte(`{"version":1,"hotkey_enabled":false,"hotkey":"G"}`)); err == nil {
		t.Fatal("an invalid shortcut must be rejected even when disabled")
	}
}

func TestSurfacesDefaultsAndValidation(t *testing.T) {
	d := Default().Surfaces
	for name, sf := range map[string]Surface{"tray": d.TrayOverflow, "start": d.StartMenu} {
		if !sf.Enabled || sf.Opacity <= 0 || sf.Opacity >= 100 {
			t.Errorf("%s default should be enabled with a partial opacity, got %+v", name, sf)
		}
	}
	c, err := Parse([]byte(`{"version":1,"surfaces":{"tray_overflow":{"enabled":false,"opacity":50}}}`))
	if err != nil || c.Surfaces.TrayOverflow.Enabled || c.Surfaces.TrayOverflow.Opacity != 50 || !c.Surfaces.StartMenu.Enabled {
		t.Fatalf("partial surfaces must keep the other defaults: %+v %v", c.Surfaces, err)
	}
	for _, bad := range []string{
		`{"version":1,"surfaces":{"start_menu_flyout":{"enabled":true,"opacity":101}}}`,
		`{"version":1,"surfaces":{"tray_overflow":{"enabled":true,"opacity":-1}}}`,
		`{"version":1,"surfaces":{"start_menu_flyout":{"enabled":true,"opacity":900}}}`,
		`{"version":1,"surfaces":{"taskbar":{"enabled":true,"opacity":5}}}`,
	} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Errorf("expected an error for %s", bad)
		}
	}
}

func TestOpacityAcceptsDecimals(t *testing.T) {
	c, err := Parse([]byte(`{"version":1,"surfaces":{"start_menu_flyout":{"enabled":true,"opacity":99.5},"tray_overflow":{"enabled":true,"opacity":60.4}},"default":{"accent":"blur","color":"auto","opacity":12.6}}`))
	if err != nil {
		t.Fatal(err)
	}
	if c.Surfaces.StartMenu.Opacity != 100 || c.Surfaces.TrayOverflow.Opacity != 60 || c.Default.Opacity != 13 {
		t.Fatalf("decimals must be rounded: %+v", c)
	}
	for _, bad := range []string{`"abc"`, `true`, `null`, `[1]`, `1e999`, `100.6`, `-0.6`} {
		if _, err := Parse([]byte(`{"version":1,"surfaces":{"start_menu_flyout":{"enabled":true,"opacity":` + bad + `}}}`)); err == nil {
			t.Errorf("opacity %s must be rejected", bad)
		}
	}
}

func TestAppWindowOpacity(t *testing.T) {
	c, err := Parse([]byte(`{"version":1,"apps":[{"exe":"Code.exe","window_opacity":92.4},{"exe":"x.exe","accent":"blur","color":"auto","opacity":30}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if c.Apps[0].WindowOpacity == nil || *c.Apps[0].WindowOpacity != 92 || c.Apps[0].StylesTaskbar() {
		t.Fatalf("window-only rule parsed wrong: %+v", c.Apps[0])
	}
	if c.Apps[1].WindowOpacity != nil || !c.Apps[1].StylesTaskbar() {
		t.Fatalf("taskbar rule parsed wrong: %+v", c.Apps[1])
	}
	for _, bad := range []string{
		`{"version":1,"apps":[{"exe":"a.exe"}]}`,
		`{"version":1,"apps":[{"exe":"a.exe","window_opacity":101}]}`,
		`{"version":1,"apps":[{"exe":"a.exe","window_opacity":-1}]}`,
	} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Errorf("expected an error for %s", bad)
		}
	}
}

func TestLegacyExplorerKeyIsAcceptedAndDropped(t *testing.T) {
	c, err := Parse([]byte(`{"version":1,"surfaces":{"explorer":{"enabled":true,"opacity":85},"tray_overflow":{"enabled":true,"opacity":70}}}`))
	if err != nil {
		t.Fatalf("an old config must keep loading: %v", err)
	}
	if c.Surfaces.TrayOverflow.Opacity != 70 {
		t.Fatalf("other surface keys must still apply: %+v", c.Surfaces)
	}
	p := filepath.Join(t.TempDir(), "c.json")
	if err := Save(p, c); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if strings.Contains(string(b), "explorer") {
		t.Fatalf("the legacy key must be dropped on save:\n%s", b)
	}
}

func TestExplicitFalseSticks(t *testing.T) {
	c, err := Parse([]byte(`{"version":1,"enabled":false,"xaml_fix":false}`))
	if err != nil || c.Enabled || c.XamlFix {
		t.Fatalf("explicit false must stick: %+v %v", c, err)
	}
}

func TestSchemaKeyAccepted(t *testing.T) {
	if _, err := Parse([]byte(`{"$schema":"` + SchemaURL + `","version":1}`)); err != nil {
		t.Fatal(err)
	}
}

func TestLoadCreatesDefaultsAndRoundTrips(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sub", "config.json")
	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("config not created: %v", err)
	}
	c.Default.Accent = AccentAcrylic
	c.Apps = []AppRule{{Exe: "code.exe", Appearance: Appearance{Accent: AccentBlur, Color: "#101010", Opacity: 40}}}
	c.Hotkey = "Ctrl+Alt+G"
	if err := Save(p, c); err != nil {
		t.Fatal(err)
	}
	got, err := Load(p)
	if err != nil || got.Default.Accent != AccentAcrylic || len(got.Apps) != 1 || got.Hotkey != "Ctrl+Alt+G" {
		t.Fatalf("round trip failed: %+v %v", got, err)
	}
	entries, _ := os.ReadDir(filepath.Dir(p))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Errorf("leftover temp file %s", e.Name())
		}
	}
}

func TestLoadMalformedLeavesFileUntouched(t *testing.T) {
	p := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(p, []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(p); err == nil {
		t.Fatal("expected error")
	}
	b, _ := os.ReadFile(p)
	if string(b) != "{broken" {
		t.Fatalf("file was modified: %q", b)
	}
}

func TestSaveRejectsInvalid(t *testing.T) {
	c := Default()
	c.PollMs = 0
	if err := Save(filepath.Join(t.TempDir(), "c.json"), c); err == nil {
		t.Fatal("expected validation error")
	}
}
