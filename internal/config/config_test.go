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

func TestParseRejectsBadInput(t *testing.T) {
	cases := map[string]string{
		"unknown field":   `{"version":1,"nope":1}`,
		"bad accent":      `{"version":1,"default":{"accent":"glow","color":"#000000","opacity":0}}`,
		"bad opacity":     `{"version":1,"default":{"accent":"clear","color":"#000000","opacity":101}}`,
		"bad version":     `{"version":2}`,
		"poll too low":    `{"version":1,"poll_ms":1}`,
		"refresh < poll":  `{"version":1,"poll_ms":500,"refresh_ms":100}`,
		"bad log level":   `{"version":1,"log_level":"loud"}`,
		"malformed json":  `{`,
		"bad nested rule": `{"version":1,"maximized":{"enabled":true,"accent":"clear","color":"nope","opacity":1}}`,
	}
	for name, doc := range cases {
		if _, err := Parse([]byte(doc)); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestParsePartialKeepsDefaults(t *testing.T) {
	c, err := Parse([]byte(`{"version":1,"poll_ms":200,"refresh_ms":2000}`))
	if err != nil {
		t.Fatal(err)
	}
	if c.PollMs != 200 || c.Default.Accent != AccentClear {
		t.Fatalf("unexpected config: %+v", c)
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
	if err := Save(p, c); err != nil {
		t.Fatal(err)
	}
	got, err := Load(p)
	if err != nil || got.Default.Accent != AccentAcrylic {
		t.Fatalf("round trip failed: %+v %v", got, err)
	}
	// no stray temp files
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
