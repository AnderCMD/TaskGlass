//go:build windows

package autostart

import "testing"

func TestCommandIsPlainQuotedPath(t *testing.T) {
	exe := `C:\Program Files\TaskGlass\taskglass.exe`
	if got, want := command(exe), `"C:\Program Files\TaskGlass\taskglass.exe"`; got != want {
		t.Fatalf("command = %s, want %s (backslashes must not be doubled)", got, want)
	}
}

func TestSameCommandRecognisesLegacyDoubledBackslashes(t *testing.T) {
	exe := `C:\Tools\TaskGlass\taskglass.exe`
	legacy := `"C:\Tools\TaskGlass\taskglass.exe"`
	if !sameCommand(legacy, exe) || !sameCommand(command(exe), exe) {
		t.Fatal("both forms must be recognised")
	}
	if sameCommand(`"D:\Other\app.exe"`, exe) {
		t.Fatal("a different program must not match")
	}
}
