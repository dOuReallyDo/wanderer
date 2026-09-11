package server

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var staticContent embed.FS

// staticFS è il filesystem embedded per la web UI.
var staticFS fs.FS

func init() {
	staticFS, _ = fs.Sub(staticContent, "static")
}