// Package assets embeds static resources into the binary.
package assets

import _ "embed"

// Color is the full-colour logo icon.
//
//go:embed taskglass.ico
var Color []byte

// GlyphLight is the white tray glyph, readable on dark taskbars.
//
//go:embed tray_light.ico
var GlyphLight []byte

// GlyphDark is the dark tray glyph, readable on light taskbars.
//
//go:embed tray_dark.ico
var GlyphDark []byte
