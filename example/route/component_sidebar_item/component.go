package component_sidebar_item

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Icon   string
	Label  string
	Link   string
	Active bool
}

func View(model Model) render.View {
	return render.Component(
		Html,
		model,
		make(map[string]render.View),
		[]render.View{},
	)
}
