package component_sidebar_button_list

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/routes/example"
	"github.com/hjwalt/platform/routes/mvc"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
}

type Component struct {
	Model      Model
	Components []mvc.Component[example.Context]
}

func (c Component) Render(ctx example.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return mvc.ComponentSliceRender[example.Context, Model](ctx, w, r, Html, c.Model, c.Components)
}
