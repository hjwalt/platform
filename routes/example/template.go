package example

import (
	"embed"
	"html/template"
)

//go:embed *
var files embed.FS

func Page(file string) *template.Template {
	return template.Must(
		template.New("page.html").ParseFS(files, "layout/page.html", file),
	)
}

func Component(file string) *template.Template {
	return template.Must(
		template.ParseFS(files, file),
	)
}
