//go:build windows

package winglass

import (
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

func pct(v int) *config.Percent { p := config.Percent(v); return &p }

func TestWanted(t *testing.T) {
	c := config.Default()
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
	if len(Wanted(c)) != 0 {
		t.Fatal("no rules, nothing wanted")
	}
	c.Apps = []config.AppRule{{Exe: "code.exe", WindowOpacity: pct(50)}}
	c.Enabled = false
	if len(Wanted(c)) != 0 {
		t.Fatal("the master switch turns everything off")
	}
}
