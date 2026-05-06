package component_sidebar_button_list

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
}

func View(model Model, elements []render.View) render.View {
	return render.Component(
		Html,
		model,
		make(map[string]render.View),
		elements,
	)
}
