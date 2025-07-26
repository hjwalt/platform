package component_sidebar_item

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
	Link   string
	Active bool
}

func (c Model) Render(ctx example.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return mvc.ComponentRender[example.Context, Model](ctx, w, r, Html, c)
}
