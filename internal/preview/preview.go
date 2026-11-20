// Package preview draws the settings window's picture of a taskbar and of the
// hidden-icons flyout. It is plain software rendering (no OS calls), so the
// picture can be tested anywhere and shows real blur instead of a CSS guess.
package preview

import (
	"image"
	"image/color"
	"math"

	"github.com/AnderCMD/TaskGlass/internal/config"
)

// Scenario selects which taskbar look is drawn.
type Scenario struct {
	// Maximized draws a maximized window above the taskbar.
	Maximized bool
	// Appearance is the look to draw (already resolved from the rules).
	Appearance config.Appearance
}

// Input is everything the picture depends on.
type Input struct {
	Scenario Scenario
	// Disabled draws the stock Windows taskbar (TaskGlass turned off).
	Disabled bool
	Light    bool
	// XamlFix reports that Windows 11's own grey fill is removed.
	XamlFix bool
	// BlurViaAcrylic mirrors the runtime behaviour on Windows 11, where Blur is
	// produced by the acrylic backdrop (see package taskbar).
	BlurViaAcrylic bool
	Flyout         config.Surface
	FlyoutActive   bool // the flyout setting is on and the helper DLL is available
}

type rgb struct{ r, g, b float64 }

func (c rgb) mix(o rgb, t float64) rgb {
	return rgb{c.r + (o.r-c.r)*t, c.g + (o.g-c.g)*t, c.b + (o.b-c.b)*t}
}

func clamp8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5)
}

func hex(r, g, b uint8) rgb { return rgb{float64(r), float64(g), float64(b)} }

// Render draws the desktop, taskbar and flyout into a w x h image.
func Render(in Input, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	barH := h * 48 / 330
	if barH < 30 {
		barH = 30
	}
	barTop := h - barH

	paintWallpaper(img)
	if in.Scenario.Maximized && !in.Disabled {
		fillRect(img, image.Rect(0, 0, w, barTop), hex(238, 238, 238), 1)
		fillRect(img, image.Rect(0, 0, w, 22), hex(222, 222, 222), 1)
	}

	// What the taskbar sees behind it.
	bar := image.Rect(0, barTop, w, h)
	backdrop := blur(img, bar, 9)

	ap := in.Scenario.Appearance
	if in.Disabled {
		ap = config.Appearance{Accent: config.AccentNormal}
	}
	tint := tintColor(ap, in.Light)
	a := float64(ap.Opacity) / 100
	for y := bar.Min.Y; y < bar.Max.Y; y++ {
		for x := bar.Min.X; x < bar.Max.X; x++ {
			under := at(img, x, y)
			bl := at(backdrop, x, y)
			var c rgb
			switch ap.Accent {
			case config.AccentClear:
				c = under.mix(tint, a)
			case config.AccentOpaque:
				c = under.mix(tint, a)
			case config.AccentBlur:
				if in.BlurViaAcrylic {
					t := lighten(tint, 0.5*(1-a))
					c = bl.mix(t, 0.65)
				} else {
					c = bl.mix(tint, a)
				}
			case config.AccentAcrylic:
				c = bl.mix(tint, 0.65)
				n := (noise(x, y) - 0.5) * 14
				c = rgb{c.r + n, c.g + n, c.b + n}
			default: // normal: the stock Windows 11 taskbar
				stock := hex(28, 28, 28)
				if in.Light {
					stock = hex(243, 243, 243)
				}
				c = bl.mix(stock, 0.82)
			}
			// Windows 11 keeps its own grey fill on top unless the helper removes it.
			if !in.Disabled && !in.XamlFix && (ap.Accent == config.AccentClear || ap.Accent == config.AccentBlur || ap.Accent == config.AccentAcrylic) {
				g := hex(32, 32, 32)
				if in.Light {
					g = hex(243, 243, 243)
				}
				c = c.mix(g, 0.55)
			}
			set(img, x, y, c)
		}
	}
	drawTaskbarContents(img, bar, in.Light)

	// Hidden-icons flyout, anchored above the tray corner.
	fw, fh := w*150/560, h*104/330
	if fw < 120 {
		fw = 120
	}
	if fh < 84 {
		fh = 84
	}
	fr := image.Rect(w-fw-w*48/560, barTop-fh-8, w-w*48/560, barTop-8)
	drawFlyout(img, fr, in)
	return img
}

func tintColor(ap config.Appearance, light bool) rgb {
	r, g, b := ap.RGB(light)
	return rgb{float64(r), float64(g), float64(b)}
}

func lighten(c rgb, t float64) rgb { return c.mix(rgb{255, 255, 255}, t) }

func at(img *image.RGBA, x, y int) rgb {
	i := img.PixOffset(x, y)
	return rgb{float64(img.Pix[i]), float64(img.Pix[i+1]), float64(img.Pix[i+2])}
}

func set(img *image.RGBA, x, y int, c rgb) {
	i := img.PixOffset(x, y)
	img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = clamp8(c.r), clamp8(c.g), clamp8(c.b), 255
}

func fillRect(img *image.RGBA, r image.Rectangle, c rgb, a float64) {
	r = r.Intersect(img.Bounds())
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			set(img, x, y, at(img, x, y).mix(c, a))
		}
	}
}

// roundRect fills a rounded rectangle with soft edges.
func roundRect(img *image.RGBA, r image.Rectangle, rad float64, c rgb, a float64) {
	b := r.Intersect(img.Bounds())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			cov := roundCoverage(float64(x)+0.5, float64(y)+0.5, r, rad)
			if cov > 0 {
				set(img, x, y, at(img, x, y).mix(c, a*cov))
			}
		}
	}
}

func roundCoverage(px, py float64, r image.Rectangle, rad float64) float64 {
	cx, cy := float64(r.Min.X+r.Max.X)/2, float64(r.Min.Y+r.Max.Y)/2
	hw, hh := float64(r.Dx())/2-rad, float64(r.Dy())/2-rad
	dx, dy := math.Abs(px-cx)-hw, math.Abs(py-cy)-hh
	d := math.Hypot(math.Max(dx, 0), math.Max(dy, 0)) + math.Min(math.Max(dx, dy), 0) - rad
	cov := 0.5 - d
	if cov < 0 {
		return 0
	}
	if cov > 1 {
		return 1
	}
	return cov
}

func noise(x, y int) float64 {
	n := uint32(x)*374761393 + uint32(y)*668265263
	n = (n ^ (n >> 13)) * 1274126177
	return float64((n^(n>>16))&0xffff) / 65535
}

func paintWallpaper(img *image.RGBA) {
	b := img.Bounds()
	w, h := float64(b.Dx()), float64(b.Dy())
	blobs := []struct {
		x, y, r float64
		c       rgb
	}{
		{0.18, 0.80, 0.22, hex(255, 209, 102)},
		{0.48, 0.98, 0.28, hex(239, 71, 111)},
		{0.80, 0.90, 0.30, hex(6, 214, 160)},
		{0.92, 0.18, 0.50, hex(17, 138, 178)},
		{0.34, 0.60, 0.18, hex(122, 92, 255)},
	}
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			t := (float64(x)/w + float64(y)/h) / 2
			var c rgb
			if t < 0.45 {
				c = hex(43, 27, 107).mix(hex(91, 42, 134), t/0.45)
			} else {
				c = hex(91, 42, 134).mix(hex(199, 75, 138), (t-0.45)/0.55)
			}
			for _, bl := range blobs {
				d := math.Hypot(float64(x)/w-bl.x, (float64(y)/h-bl.y)*h/w) / bl.r
				if d < 1 {
					c = c.mix(bl.c, (1-d)*(1-d)*1.3)
				}
			}
			set(img, x, y, c)
		}
	}
}

// blur returns a copy of img where region r holds a box-blurred (three passes,
// so close to gaussian) version of the same region.
func blur(img *image.RGBA, r image.Rectangle, radius int) *image.RGBA {
	out := image.NewRGBA(img.Bounds())
	copy(out.Pix, img.Pix)
	r = r.Intersect(img.Bounds())
	w, h := r.Dx(), r.Dy()
	buf := make([]rgb, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			buf[y*w+x] = at(img, r.Min.X+x, r.Min.Y+y)
		}
	}
	tmp := make([]rgb, w*h)
	for pass := 0; pass < 3; pass++ {
		boxH(buf, tmp, w, h, radius)
		boxV(tmp, buf, w, h, radius)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			set(out, r.Min.X+x, r.Min.Y+y, buf[y*w+x])
		}
	}
	return out
}

func boxH(src, dst []rgb, w, h, r int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var s rgb
			n := 0
			for k := -r; k <= r; k++ {
				xx := x + k
				if xx < 0 {
					xx = 0
				}
				if xx >= w {
					xx = w - 1
				}
				p := src[y*w+xx]
				s = rgb{s.r + p.r, s.g + p.g, s.b + p.b}
				n++
			}
			dst[y*w+x] = rgb{s.r / float64(n), s.g / float64(n), s.b / float64(n)}
		}
	}
}

func boxV(src, dst []rgb, w, h, r int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var s rgb
			n := 0
			for k := -r; k <= r; k++ {
				yy := y + k
				if yy < 0 {
					yy = 0
				}
				if yy >= h {
					yy = h - 1
				}
				p := src[yy*w+x]
				s = rgb{s.r + p.r, s.g + p.g, s.b + p.b}
				n++
			}
			dst[y*w+x] = rgb{s.r / float64(n), s.g / float64(n), s.b / float64(n)}
		}
	}
}

func drawTaskbarContents(img *image.RGBA, bar image.Rectangle, light bool) {
	fg := hex(255, 255, 255)
	if light {
		fg = hex(27, 27, 27)
	}
	h := bar.Dy()
	ic := h * 28 / 48
	gap := ic / 4
	// centred app icons
	colors := []rgb{hex(76, 194, 255), hex(128, 128, 128), hex(232, 69, 60), hex(45, 127, 249), hex(51, 180, 105), hex(245, 179, 1)}
	total := len(colors)*ic + (len(colors)-1)*gap
	x := bar.Min.X + (bar.Dx()-total)/2
	y := bar.Min.Y + (h-ic)/2
	for i, c := range colors {
		r := image.Rect(x, y, x+ic, y+ic)
		switch i {
		case 0: // start: four panes
			p := ic/2 - 1
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					fillRect(img, image.Rect(x+dx*(p+2), y+dy*(p+2), x+dx*(p+2)+p, y+dy*(p+2)+p), c, 1)
				}
			}
		case 1:
			roundRect(img, r, 6, c, 0.35)
		default:
			roundRect(img, r, 6, c, 1)
		}
		x += ic + gap
	}
	// tray: chevron box, two glyphs, clock lines
	tx := bar.Max.X - h/4
	clockW := ic * 2
	fillRect(img, image.Rect(tx-clockW, y+ic/4, tx, y+ic/4+3), fg, 0.8)
	fillRect(img, image.Rect(tx-clockW+6, y+ic/4+8, tx, y+ic/4+11), fg, 0.6)
	tx -= clockW + gap*2
	for i := 0; i < 2; i++ {
		roundRect(img, image.Rect(tx-ic/2, y+ic/4, tx, y+ic/4+ic/2), 3, fg, 0.7)
		tx -= ic/2 + gap
	}
	roundRect(img, image.Rect(tx-ic*2/3, y+ic/6, tx, y+ic*5/6), 5, hex(128, 128, 128), 0.35)
	cx, cy := tx-ic/3, y+ic/2+2
	for d := 0; d < 4; d++ { // ^
		set(img, cx-d, cy-4+d+1, fg)
		set(img, cx+d, cy-4+d+1, fg)
		set(img, cx-d, cy-4+d, fg)
		set(img, cx+d, cy-4+d, fg)
	}
}

func drawFlyout(img *image.RGBA, r image.Rectangle, in Input) {
	a := 0.92
	if in.FlyoutActive && !in.Disabled {
		a = float64(in.Flyout.Opacity) / 100
	}
	base := hex(36, 36, 36)
	fg := hex(255, 255, 255)
	if in.Light {
		base, fg = hex(243, 243, 243), hex(27, 27, 27)
	}
	bl := blur(img, r, 9)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cov := roundCoverage(float64(x)+0.5, float64(y)+0.5, r, 8)
			if cov == 0 {
				continue
			}
			c := at(bl, x, y).mix(base, a)
			set(img, x, y, at(img, x, y).mix(c, cov))
		}
	}
	// 1px border
	edge := hex(255, 255, 255)
	if in.Light {
		edge = hex(0, 0, 0)
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			outer := roundCoverage(float64(x)+0.5, float64(y)+0.5, r, 8)
			inner := roundCoverage(float64(x)+0.5, float64(y)+0.5, image.Rect(r.Min.X+1, r.Min.Y+1, r.Max.X-1, r.Max.Y-1), 7)
			if outer-inner > 0 {
				set(img, x, y, at(img, x, y).mix(edge, 0.14*(outer-inner)))
			}
		}
	}
	// six icon tiles
	pad := 10
	cw := (r.Dx() - pad*2 - 2*8) / 3
	ch := (r.Dy() - pad*2 - 8) / 2
	cols := []rgb{hex(76, 194, 255), hex(108, 203, 95), hex(255, 185, 0), hex(180, 140, 255), hex(255, 120, 120), hex(128, 128, 128)}
	for i := 0; i < 6; i++ {
		cx := r.Min.X + pad + (i%3)*(cw+8)
		cy := r.Min.Y + pad + (i/3)*(ch+8)
		roundRect(img, image.Rect(cx, cy, cx+cw, cy+ch), 6, fg, 0.12)
		inset := cw / 4
		roundRect(img, image.Rect(cx+inset, cy+ch/4, cx+cw-inset, cy+ch-ch/4), 4, cols[i], 0.9)
	}
}

// Color returns the average colour of region r, used by tests.
func Color(img *image.RGBA, r image.Rectangle) color.RGBA {
	var sr, sg, sb, n float64
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			c := at(img, x, y)
			sr, sg, sb, n = sr+c.r, sg+c.g, sb+c.b, n+1
		}
	}
	return color.RGBA{clamp8(sr / n), clamp8(sg / n), clamp8(sb / n), 255}
}
