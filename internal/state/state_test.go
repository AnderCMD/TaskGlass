package state

import (
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

func TestResolve(t *testing.T) {
	base := config.Default()
	base.Maximized.Enabled = true
	base.StartMenu.Enabled = true
	base.Apps = []config.AppRule{{Exe: "code.exe", Appearance: config.Appearance{Accent: config.AccentBlur, Color: "auto", Opacity: 30}}}

	noMax := base
	noMax.Maximized.Enabled = false
	noStart := base
	noStart.StartMenu.Enabled = false
	off := base
	off.Enabled = false

	tests := []struct {
		name string
		cfg  config.Config
		ctx  Context
		want Kind
	}{
		{"idle", base, Context{}, KindDefault},
		{"maximized", base, Context{Maximized: true}, KindMaximized},
		{"start", base, Context{StartMenu: true}, KindStartMenu},
		{"start beats maximized", base, Context{StartMenu: true, Maximized: true}, KindStartMenu},
		{"maximized rule disabled", noMax, Context{Maximized: true}, KindDefault},
		{"start rule disabled falls to maximized", noStart, Context{StartMenu: true, Maximized: true}, KindMaximized},
		{"start rule disabled alone", noStart, Context{StartMenu: true}, KindDefault},
		{"app rule", base, Context{ForegroundExe: "Code.exe"}, KindApp},
		{"app rule beats maximized", base, Context{ForegroundExe: "code.exe", Maximized: true}, KindApp},
		{"start beats app rule", base, Context{ForegroundExe: "code.exe", StartMenu: true}, KindStartMenu},
		{"unmatched app falls through", base, Context{ForegroundExe: "notepad.exe", Maximized: true}, KindMaximized},
		{"master switch off", off, Context{Maximized: true, StartMenu: true, ForegroundExe: "code.exe"}, KindDisabled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := Resolve(tt.cfg, tt.ctx)
			if got != tt.want {
				t.Fatalf("got %s want %s", got, tt.want)
			}
		})
	}
}

func TestDisabledResolvesToNormal(t *testing.T) {
	c := config.Default()
	c.Enabled = false
	ap, kind := Resolve(c, Context{})
	if ap.Accent != config.AccentNormal || kind != KindDisabled {
		t.Fatalf("got %+v %s", ap, kind)
	}
}

func TestResolveReturnsMatchingAppearance(t *testing.T) {
	c := config.Default()
	c.Maximized.Enabled = true
	ap, _ := Resolve(c, Context{Maximized: true})
	if ap != c.Maximized.Appearance {
		t.Fatalf("got %+v", ap)
	}
}

func TestFirstMatchingAppRuleWins(t *testing.T) {
	c := config.Default()
	first := config.Appearance{Accent: config.AccentBlur, Color: "auto", Opacity: 10}
	second := config.Appearance{Accent: config.AccentOpaque, Color: "auto", Opacity: 100}
	c.Apps = []config.AppRule{{Exe: "a.exe", Appearance: first}, {Exe: "a.exe", Appearance: second}}
	ap, _ := Resolve(c, Context{ForegroundExe: "a.exe"})
	if ap != first {
		t.Fatalf("got %+v, want the first rule", ap)
	}
}

func TestNeedsWindowScan(t *testing.T) {
	c := config.Default()
	c.Maximized.Enabled = false
	if m, f := NeedsWindowScan(c); m || f {
		t.Fatal("nothing to scan for")
	}
	c.Maximized.Enabled = true
	c.Apps = []config.AppRule{{Exe: "a.exe", Appearance: c.Default}}
	if m, f := NeedsWindowScan(c); !m || !f {
		t.Fatal("both scans needed")
	}
	c.Enabled = false
	if m, f := NeedsWindowScan(c); m || f {
		t.Fatal("disabled config needs no scan")
	}
}

func TestWindowOnlyAppRuleNeverRestylesTheTaskbar(t *testing.T) {
	c := config.Default()
	c.Maximized.Enabled = false
	op := config.Percent(80)
	c.Apps = []config.AppRule{{Exe: "code.exe", WindowOpacity: &op}}
	if _, kind := Resolve(c, Context{ForegroundExe: "code.exe"}); kind != KindDefault {
		t.Fatalf("got %s, want default", kind)
	}
	if _, fg := NeedsWindowScan(c); fg {
		t.Fatal("a window-only rule needs no foreground tracking")
	}
	c.Apps = append(c.Apps, config.AppRule{Exe: "code.exe", Appearance: config.Appearance{Accent: config.AccentBlur, Color: "auto", Opacity: 30}})
	if _, kind := Resolve(c, Context{ForegroundExe: "code.exe"}); kind != KindApp {
		t.Fatalf("got %s, want app", kind)
	}
}
