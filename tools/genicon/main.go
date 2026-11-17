// Command genicon renders assets/taskglass.ico (a translucent taskbar glyph).
// Run: go run ./tools/genicon
package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

const size = 32

func main() {
	px := make([]byte, size*size*4) // BGRA, top-down
	set := func(x, y int, r, g, b, a byte) {
		i := (y*size + x) * 4
		px[i], px[i+1], px[i+2], px[i+3] = b, g, r, a
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// rounded "screen"
			in := x >= 2 && x < 30 && y >= 4 && y < 28
			if !in {
				continue
			}
			// trim the four corners to round them off
			dx, dy := 0, 0
			if x < 5 {
				dx = 5 - x
			} else if x > 26 {
				dx = x - 26
			}
			if y < 7 {
				dy = 7 - y
			} else if y > 24 {
				dy = y - 24
			}
			if dx+dy > 5 {
				continue
			}
			// glassy gradient background
			t := byte(40 + y*4)
			set(x, y, 30, t, 120+byte(y*3), 235)
			// taskbar strip at the bottom: lighter, more transparent
			if y >= 22 {
				set(x, y, 235, 245, 255, 150)
			}
		}
	}
	// three "icons" on the strip
	for _, cx := range []int{8, 15, 22} {
		for y := 23; y < 27; y++ {
			for x := cx; x < cx+4; x++ {
				set(x, y, 255, 255, 255, 255)
			}
		}
	}

	var dib bytes.Buffer
	hdr := []any{uint32(40), int32(size), int32(size * 2), uint16(1), uint16(32),
		uint32(0), uint32(size * size * 4), int32(0), int32(0), uint32(0), uint32(0)}
	for _, v := range hdr {
		binary.Write(&dib, binary.LittleEndian, v)
	}
	for y := size - 1; y >= 0; y-- { // bottom-up rows
		dib.Write(px[y*size*4 : (y+1)*size*4])
	}
	dib.Write(make([]byte, size*4)) // AND mask (32 rows * 4 bytes)

	var ico bytes.Buffer
	for _, v := range []any{uint16(0), uint16(1), uint16(1)} {
		binary.Write(&ico, binary.LittleEndian, v)
	}
	ico.Write([]byte{size, size, 0, 0})
	for _, v := range []any{uint16(1), uint16(32), uint32(dib.Len()), uint32(22)} {
		binary.Write(&ico, binary.LittleEndian, v)
	}
	ico.Write(dib.Bytes())

	if err := os.MkdirAll("assets", 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("assets/taskglass.ico", ico.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
