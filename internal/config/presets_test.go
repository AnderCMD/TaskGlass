package config

import (
	"strings"
	"testing"
)

func TestEveryPresetIsValidAndKeepsPersonalSettings(t *testing.T) {
	op := Percent(80)
	base := Default()
	base.Hotkey = "Ctrl+Alt+H"
	base.HotkeyEnabled = true
	base.Theme = ThemeLight
	base.Apps = []AppRule{{Exe: "code.exe", WindowOpacity: &op}}
	for _, p := range Presets() {
		got, err := ApplyPreset(base, p.Name)
		if err != nil {
			t.Fatalf("%s: %v", p.Name, err)
		}
		if got.Hotkey != "Ctrl+Alt+H" || !got.HotkeyEnabled || got.Theme != ThemeLight || len(got.Apps) != 1 {
			t.Errorf("%s dropped personal settings: %+v", p.Name, got)
		}
		if p.Name != "stock" && !got.Enabled {
			t.Errorf("%s must enable TaskGlass", p.Name)
		}
		if p.Description == "" {
			t.Errorf("%s needs a description", p.Name)
		}
	}
}

func TestPresetDoesNotMutateTheInput(t *testing.T) {
	base := Default()
	before := base.Default
	if _, err := ApplyPreset(base, "liquid"); err != nil {
		t.Fatal(err)
	}
	if base.Default != before {
		t.Fatal("ApplyPreset must work on a copy")
	}
}

func TestLiquidIsTranslucentEverywhere(t *testing.T) {
	got, err := ApplyPreset(Default(), " LIQUID ")
	if err != nil {
		t.Fatal(err)
	}
	if got.Default.Accent != AccentAcrylic || got.Surfaces.Explorer.Opacity >= 100 || got.Surfaces.StartMenu.Opacity >= 100 || got.Surfaces.TrayOverflow.Opacity >= 100 {
		t.Fatalf("liquid should frost the taskbar and make every surface translucent: %+v", got)
	}
}

func TestUnknownPreset(t *testing.T) {
	_, err := ApplyPreset(Default(), "neon")
	if err == nil || !strings.Contains(err.Error(), "liquid") {
		t.Fatalf("the error must list the valid presets, got %v", err)
	}
}
