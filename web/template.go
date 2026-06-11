package web

import (
	"embed"
	"html/template"
)

//go:embed template/*.html
var templateFS embed.FS

func HandleTmpl() *template.Template {
	return template.Must(template.ParseFS(
		templateFS,
		"template/*.html",
	))
}
