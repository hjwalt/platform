package component_navbar

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/page/page_error_500"
	"github.com/hjwalt/routes/mvc"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/runtime"
)

const (
	path = "/component/navbar"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "instance.html"))

type Model struct {
}

type model struct {
}

func controller(c web.Context, w http.ResponseWriter, r *http.Request) (mvc.View[web.Context], error) {
	return mvc.ComponentSlice[web.Context, model]{
		Template:   Html,
		Model:      model{},
		Components: []mvc.Component[web.Context]{},
	}, nil
}

func Get() runtime.Configuration[*runtime_chi.Runtime[web.Context]] {
	return runtime_chi.WithController(path, http.MethodGet, controller, page_error_500.Controller)
}
