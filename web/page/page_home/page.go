package page_home

import (
	"embed"
	"net/http"

	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/component/component_sidebar"
	"github.com/hjwalt/platform/web/layout"
	"github.com/hjwalt/platform/web/page/page_error_500"
	"github.com/hjwalt/platform/web/render"
	"github.com/hjwalt/platform/web/route"
)

//go:embed *
var files embed.FS

var Html = render.Embedded(files, "page.html")

const (
	path = "/"
)

type model struct {
}

func handler(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			Html,
			model{},
			map[string]render.View{},
			[]render.View{},
		)), nil
}

func Add(builder route.Builder[web.Context]) {
	builder.Handle(path, http.MethodGet, render.Handle(handler, page_error_500.Error))
}
