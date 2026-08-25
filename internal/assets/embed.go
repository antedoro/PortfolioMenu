package assets

import "embed"

//go:embed templates/*.html
var Templates embed.FS

//go:embed icon.png
var Icon []byte
