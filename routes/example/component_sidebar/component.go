package component_sidebar

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
	Icon   string
	Label  string
	Active bool
}

type Component struct {
	Model  Model
	Top    mvc.Component[example.Context]
	Button mvc.Component[example.Context]
}

func (c Component) Render(ctx example.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	items := map[string]mvc.Component[example.Context]{}

	items["top"] = c.Top
	items["button"] = c.Button

	return mvc.ComponentMapRender(ctx, w, r, Html, c.Model, items)
}
