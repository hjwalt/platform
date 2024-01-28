package page_billing

import (
	"embed"
	"net/http"

	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/example/route/component_sidebar"
	"github.com/hjwalt/platform/example/route/layout"
	"github.com/hjwalt/platform/example/route/page_error_500"
	"github.com/hjwalt/platform/web/render"
	"github.com/hjwalt/platform/web/route"
)

//go:embed *
var files embed.FS

var Html = render.Embedded(files, "page.html")

const (
	path = "/billing"
)

type model struct {
}

func handler(c example.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			Html,
			model{},
			make(map[string]render.View),
			[]render.View{},
		)), nil
}

func Add(builder route.Builder[example.Context]) {
	builder.Handle(path, http.MethodGet, render.Handle(handler, page_error_500.Error))
}
