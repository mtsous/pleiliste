package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

func HandleStatic(mux *http.ServeMux) {
	sub, _ := fs.Sub(staticFS, "static")

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
}
