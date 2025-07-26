package component_sidebar

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/routes/runtime_chi"
	"github.com/hjwalt/platform/runway/runtime"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/page/page_error_500"
	"github.com/hjwalt/platform/web/view/view_sidebar_item"
	"github.com/hjwalt/platform/web/view/view_sidebar_item_header"
	"github.com/hjwalt/platform/web/view/view_sidebar_item_list"
)

const (
	path = "/component/sidebar"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "instance.html"))

type Model struct {
}

type model struct {
}

func sidebar() []mvc.Component[web.Context] {
	sidebarTop := view_sidebar_item_list.Component{
		Model: view_sidebar_item_list.Model{},
		Components: []mvc.Component[web.Context]{
			view_sidebar_item.Model{Icon: "bi bi-book", Label: "Dashboard", Link: "/", Active: true},
			view_sidebar_item.Model{Icon: "receipt_long", Label: "Billing", Link: "/billing", Active: false},
		},
	}
	sidebarButton := view_sidebar_item_list.Component{
		Model: view_sidebar_item_list.Model{},
		Components: []mvc.Component[web.Context]{
			view_sidebar_item.Model{Label: "Documentation", Link: "https://getbootstrap.com/docs/5.3/getting-started/introduction/", Active: false},
		},
	}

	return []mvc.Component[web.Context]{
		sidebarTop,
		view_sidebar_item_header.Model{Label: "Documentations"},
		sidebarButton,
	}
}

func controller(c web.Context, w http.ResponseWriter, r *http.Request) (mvc.View[web.Context], error) {
	return mvc.ComponentSlice[web.Context, model]{
		Template:   Html,
		Model:      model{},
		Components: sidebar(),
	}, nil
}

func Get() runtime.Configuration[*runtime_chi.Runtime[web.Context]] {
	return runtime_chi.WithController(path, http.MethodGet, controller, page_error_500.Controller)
}
