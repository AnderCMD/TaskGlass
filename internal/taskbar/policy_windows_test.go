//go:build windows

package taskbar

import (
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
	"github.com/AnderCMD/TaskGlass/internal/winapi"
)

func TestToPolicy(t *testing.T) {
	tests := []struct {
		name  string
		ap    config.Appearance
		light bool
		want  policy
	}{
		{"normal disables", config.Appearance{Accent: config.AccentNormal, Color: "#112233"}, false, policy{winapi.AccentDisabled, 0}},
		{"clear never alpha 0", config.Appearance{Accent: config.AccentClear, Color: "#000000", Opacity: 0}, false, policy{winapi.AccentEnableTransparent, 0x01000000}},
		{"opaque full", config.Appearance{Accent: config.AccentOpaque, Color: "#112233", Opacity: 100}, false, policy{winapi.AccentEnableGradient, 0xFF332211}},
		{"blur half, BGR order", config.Appearance{Accent: config.AccentBlur, Color: "#FF0000", Opacity: 50}, false, blurPolicy(0x7F0000FF)},
		{"acrylic", config.Appearance{Accent: config.AccentAcrylic, Color: "#202020", Opacity: 60}, false, policy{winapi.AccentEnableAcrylicBlurBG, 0x99202020}},
		{"auto colour, dark", config.Appearance{Accent: config.AccentOpaque, Color: "auto", Opacity: 100}, false, policy{winapi.AccentEnableGradient, 0xFF1C1C1C}},
		{"auto colour, light", config.Appearance{Accent: config.AccentOpaque, Color: "auto", Opacity: 100}, true, policy{winapi.AccentEnableGradient, 0xFFF3F3F3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toPolicy(tt.ap, tt.light); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

// blurPolicy is what blur resolves to on this machine.
func blurPolicy(classic uint32) policy {
	if blurViaAcrylic {
		return policy{winapi.AccentEnableAcrylicBlurBG, 0x014040FF} // #FF0000 mixed 25% to white, BGR
	}
	return policy{winapi.AccentEnableBlurBehind, classic}
}
