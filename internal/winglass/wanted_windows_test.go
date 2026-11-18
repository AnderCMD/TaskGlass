//go:build windows

package winglass

import (
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

func pct(v int) *config.Percent { p := config.Percent(v); return &p }

func TestWanted(t *testing.T) {
	c := config.Default()
	c.Surfaces.Explorer = config.Surface{Enabled: true, Opacity: 85}
	c.Apps = []config.AppRule{
		{Exe: "Code.exe", WindowOpacity: pct(90)},
		{Exe: "code.exe", WindowOpacity: pct(10)},      // duplicate: the first rule wins
		{Exe: "chrome.exe"},                            // no window opacity
		{Exe: "explorer.exe", WindowOpacity: pct(100)}, // an app rule can cancel the surface
	}
	got := Wanted(c)
	if len(got) != 1 || got["code.exe"] != Alpha(90) {
		t.Fatalf("unexpected: %v", got)
	}
	c.Apps = nil
	got = Wanted(c)
	if got["explorer.exe"] != Alpha(85) || len(got) != 1 {
		t.Fatalf("surface only: %v", got)
	}
	c.Surfaces.Explorer.Opacity = 100
	if len(Wanted(c)) != 0 {
		t.Fatal("100 % means untouched")
	}
	c.Surfaces.Explorer.Opacity = 50
	c.Enabled = false
	if len(Wanted(c)) != 0 {
		t.Fatal("the master switch turns everything off")
	}
}
