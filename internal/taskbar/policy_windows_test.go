//go:build windows

package taskbar

import (
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

func TestToPolicy(t *testing.T) {
	tests := []struct {
		name      string
		ap        config.Appearance
		wantState winapi.AccentState
		wantColor uint32
	}{
		{"normal disables", config.Appearance{Accent: config.AccentNormal, Color: "#112233"}, winapi.AccentDisabled, 0},
		{"clear never alpha 0", config.Appearance{Accent: config.AccentClear, Color: "#000000", Opacity: 0}, winapi.AccentEnableTransparent, 0x01000000},
		{"opaque full", config.Appearance{Accent: config.AccentOpaque, Color: "#112233", Opacity: 100}, winapi.AccentEnableGradient, 0xFF332211},
		{"blur half, BGR order", config.Appearance{Accent: config.AccentBlur, Color: "#FF0000", Opacity: 50}, winapi.AccentEnableBlurBehind, 0x7F0000FF},
		{"acrylic", config.Appearance{Accent: config.AccentAcrylic, Color: "#202020", Opacity: 60}, winapi.AccentEnableAcrylicBlurBG, 0x99202020},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c := toPolicy(tt.ap)
			if s != tt.wantState || c != tt.wantColor {
				t.Fatalf("got state=%d color=%#x, want state=%d color=%#x", s, c, tt.wantState, tt.wantColor)
			}
		})
	}
}
