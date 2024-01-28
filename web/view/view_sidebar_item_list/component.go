package view_sidebar_item_list

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/routes/mvc"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
}

type Component struct {
	Model      Model
	Components []mvc.Component[web.Context]
}

func (c Component) Render(ctx web.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return mvc.ComponentSliceRender[web.Context, Model](ctx, w, r, Html, c.Model, c.Components)
}
