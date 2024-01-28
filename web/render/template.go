package render

import (
	"embed"
	"html/template"
	"path"
)

func Embedded(files embed.FS, file string) *template.Template {
	return template.Must(template.ParseFS(files, file))
}

func Template(base string, file string) *template.Template {
	return template.Must(template.ParseFiles(path.Join(base, file)))
}
