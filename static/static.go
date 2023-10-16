package static

import (
	"embed"
)

//go:embed css/*
var CssFiles embed.FS

//go:embed js/*
var JsFiles embed.FS
