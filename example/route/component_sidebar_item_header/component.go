package component_sidebar_item_header

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Label string
}

func View(model Model) render.View {
	return render.Component(
		Html,
		model,
		make(map[string]render.View),
		[]render.View{},
	)
}
