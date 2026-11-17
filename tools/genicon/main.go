// Command genicon packs rendered PNGs into Windows .ico files.
//
// It reads <in>/<set>-<size>.png for the sets below and writes
// <out>/<name>.ico with every size it finds (PNG-compressed entries, which
// Windows Vista and later load natively). Run scripts/render-assets.ps1 to
// produce the PNGs and call this tool.
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var sets = []struct{ prefix, name string }{
	{"color", "taskglass.ico"},  // app icon / exe icon / coloured tray icon
	{"light", "tray_light.ico"}, // white glyph, for dark taskbars
	{"dark", "tray_dark.ico"},   // dark glyph, for light taskbars
}

var sizes = []int{16, 20, 24, 32, 40, 48, 64, 128, 256}

func main() {
	in := flag.String("in", "", "directory with <set>-<size>.png files")
	out := flag.String("out", "assets", "output directory")
	flag.Parse()
	if *in == "" {
		log.Fatal("-in is required")
	}
	for _, s := range sets {
		var imgs []entry
		for _, sz := range sizes {
			p := filepath.Join(*in, fmt.Sprintf("%s-%d.png", s.prefix, sz))
			data, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			cfg, err := png.DecodeConfig(bytes.NewReader(data))
			if err != nil {
				log.Fatalf("%s: %v", p, err)
			}
			if cfg.Width != sz || cfg.Height != sz {
				log.Fatalf("%s is %dx%d, want %dx%d", p, cfg.Width, cfg.Height, sz, sz)
			}
			imgs = append(imgs, entry{sz, data})
		}
		if len(imgs) == 0 {
			log.Fatalf("no PNGs for set %q in %s", s.prefix, *in)
		}
		if err := writeICO(filepath.Join(*out, s.name), imgs); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote %s (%d sizes)\n", s.name, len(imgs))
	}
}

type entry struct {
	size int
	png  []byte
}

func writeICO(path string, imgs []entry) error {
	sort.Slice(imgs, func(i, j int) bool { return imgs[i].size < imgs[j].size })
	var b bytes.Buffer
	w := func(v any) { _ = binary.Write(&b, binary.LittleEndian, v) }
	w(uint16(0))
	w(uint16(1))
	w(uint16(len(imgs)))
	offset := 6 + 16*len(imgs)
	for _, e := range imgs {
		dim := byte(e.size)
		if e.size >= 256 {
			dim = 0 // 0 means 256 in the ICO directory
		}
		b.Write([]byte{dim, dim, 0, 0})
		w(uint16(1))  // planes
		w(uint16(32)) // bit depth
		w(uint32(len(e.png)))
		w(uint32(offset))
		offset += len(e.png)
	}
	for _, e := range imgs {
		b.Write(e.png)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b.Bytes(), 0o644)
}
