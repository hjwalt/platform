package page_billing

import (
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/routes/example"
	"github.com/hjwalt/platform/routes/example/page_error_500"
	"github.com/hjwalt/platform/routes/runtime_chi"
	"github.com/hjwalt/platform/runway/inverse"
	"github.com/hjwalt/platform/runway/runtime"
)

const (
	directory = "page_billing"
	path      = "/billing"
)

var html = example.Page(directory + "/page.html")

type model struct {
}

func page(c example.Context, w http.ResponseWriter, r *http.Request) (*template.Template, model, error) {
	return html, model{}, nil
}

func Get() runtime.Configuration[*runtime_chi.Runtime[example.Context]] {
	return runtime_chi.WithPage(path, http.MethodGet, page, page_error_500.Error)
}

func Add(ic inverse.Container) {
	runtime_chi.AddRoute[example.Context](
		ic,
		runtime_chi.AddPage(path, http.MethodGet, page, page_error_500.Error),
	)
}
