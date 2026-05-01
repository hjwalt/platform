package chat_list

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var html = template.Must(template.ParseFS(files, "component.html"))

func View(elements []render.View) render.View {
	return render.Component(
		html,
		"",
		make(map[string]render.View),
		elements,
	)
}
