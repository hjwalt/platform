package view_sidebar_item_header

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/web"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Label string
}

func (c Model) Render(ctx web.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return mvc.ComponentRender[web.Context, Model](ctx, w, r, Html, c)
}
