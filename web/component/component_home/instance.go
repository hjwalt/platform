package component_home

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/page/page_error_500"
	"github.com/hjwalt/platform/web/view/view_protobuf_message"
	"github.com/hjwalt/routes/mvc"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/runtime"
)

const (
	path = "/component/home"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "instance.html"))

type Model struct {
}

func schemaView(t *model.ProtobufType) web.Component {
	return &view_protobuf_message.Model{
		Model: t.GetMessage(),
	}
}

func content(c web.Context) []web.Component {
	result := []web.Component{}

	for _, t := range c.Schema.GetTypes() {
		result = append(result, schemaView(t))
	}

	return result
}

func controller(c web.Context, w http.ResponseWriter, r *http.Request) (web.View, error) {
	return mvc.ComponentSlice[web.Context, Model]{
		Template:   Html,
		Model:      Model{},
		Components: content(c),
	}, nil
}

func Get() runtime.Configuration[*runtime_chi.Runtime[web.Context]] {
	return runtime_chi.WithController(path, http.MethodGet, controller, page_error_500.Controller)
}
