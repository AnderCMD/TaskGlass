// Package assets embeds static resources into the binary.
package assets

import _ "embed"

// Icon is the tray icon in .ico format. Regenerate with `go run ./tools/genicon`.
//
//go:embed taskglass.ico
var Icon []byte
