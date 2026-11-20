package preview

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

func base(ac config.Accent, op config.Percent) Input {
	return Input{
		Scenario:       Scenario{Appearance: config.Appearance{Accent: ac, Color: "auto", Opacity: op}},
		XamlFix:        true,
		BlurViaAcrylic: true,
		Flyout:         config.Surface{Enabled: true, Opacity: 10},
		FlyoutActive:   true,
	}
}

func barRegion(img *image.RGBA) image.Rectangle {
	b := img.Bounds()
	return image.Rect(0, b.Max.Y-40, 60, b.Max.Y-4) // left of the icons
}

func TestAccentsLookDifferent(t *testing.T) {
	const w, h = 560, 300
	clear := Color(Render(base(config.AccentClear, 0), w, h), barRegion(Render(base(config.AccentClear, 0), w, h)))
	blur := Color(Render(base(config.AccentBlur, 35), w, h), image.Rect(0, h-40, 60, h-4))
	opaque := Color(Render(base(config.AccentOpaque, 100), w, h), image.Rect(0, h-40, 60, h-4))
	if opaque.R != 28 || opaque.G != 28 || opaque.B != 28 {
		t.Errorf("opaque auto/dark should be #1C1C1C, got %v", opaque)
	}
	if clear == blur || clear == opaque {
		t.Errorf("styles are indistinguishable: clear=%v blur=%v opaque=%v", clear, blur, opaque)
	}
	if blur.R < 20 && blur.G < 20 && blur.B < 20 {
		t.Errorf("blur must not be black: %v", blur)
	}
}

func TestMaximizedAndFlyoutDraw(t *testing.T) {
	in := base(config.AccentClear, 0)
	in.Scenario.Maximized = true
	img := Render(in, 560, 300)
	if c := Color(img, image.Rect(100, 60, 200, 100)); c.R < 200 {
		t.Errorf("maximized window should be light, got %v", c)
	}
}

// TestWritePNGs writes sample renders when TG_PREVIEW_DIR is set (manual check).
func TestWritePNGs(t *testing.T) {
	dir := os.Getenv("TG_PREVIEW_DIR")
	if dir == "" {
		t.Skip()
	}
	cases := map[string]Input{
		"clear": base(config.AccentClear, 0), "blur": base(config.AccentBlur, 35),
		"acrylic": base(config.AccentAcrylic, 55), "opaque": base(config.AccentOpaque, 100), "normal": base(config.AccentNormal, 100),
	}
	for name, in := range cases {
		f, _ := os.Create(filepath.Join(dir, name+".png"))
		png.Encode(f, Render(in, 560, 300))
		f.Close()
	}
}
