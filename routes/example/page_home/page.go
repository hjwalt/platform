package page_home

import (
	"net/http"

	"github.com/hjwalt/platform/routes/example"
	"github.com/hjwalt/platform/routes/example/component_sidebar"
	"github.com/hjwalt/platform/routes/example/component_sidebar_button"
	"github.com/hjwalt/platform/routes/example/component_sidebar_button_list"
	"github.com/hjwalt/platform/routes/example/component_sidebar_item"
	"github.com/hjwalt/platform/routes/example/component_sidebar_item_header"
	"github.com/hjwalt/platform/routes/example/component_sidebar_item_list"
	"github.com/hjwalt/platform/routes/example/page_error_500"
	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/routes/runtime_chi"
	"github.com/hjwalt/platform/runway/inverse"
	"github.com/hjwalt/platform/runway/runtime"
)

const (
	directory = "page_home"
	path      = "/"
)

var Html = example.Page(directory + "/page.html")

type model struct {
}

func sidebar() mvc.Component[example.Context] {
	sidebarTop := component_sidebar_item_list.Component{
		Model: component_sidebar_item_list.Model{},
		Components: []mvc.Component[example.Context]{
			component_sidebar_item.Model{Icon: "dashboard", Label: "Dashboard", Link: "/", Active: true},
			component_sidebar_item.Model{Icon: "table_view", Label: "Tables", Link: "/pages/tables.html", Active: false},
			component_sidebar_item.Model{Icon: "receipt_long", Label: "Billing", Link: "/billing", Active: false},
			component_sidebar_item.Model{Icon: "view_in_ar", Label: "Virtual Reality", Link: "/pages/virtual-reality.html", Active: false},
			component_sidebar_item.Model{Icon: "format_textdirection_r_to_l", Label: "RTL", Link: "/pages/rtl.html", Active: false},
			component_sidebar_item.Model{Icon: "notifications", Label: "Notifications", Link: "/pages/notifications.html", Active: false},
			component_sidebar_item_header.Model{Label: "Account pages"},
			component_sidebar_item.Model{Icon: "person", Label: "Profile", Link: "/pages/profile.html", Active: false},
			component_sidebar_item.Model{Icon: "login", Label: "Sign In", Link: "/pages/sign-in.html", Active: false},
			component_sidebar_item.Model{Icon: "assignment", Label: "Sign Up", Link: "/pages/sign-up.html", Active: false},
		},
	}
	sidebarButton := component_sidebar_button_list.Component{
		Model: component_sidebar_button_list.Model{},
		Components: []mvc.Component[example.Context]{
			component_sidebar_button.Model{Label: "Documentation", Link: "https://www.creative-tim.com/learning-lab/bootstrap/overview/material-dashboard?ref=sidebarfree", Outlined: true},
			component_sidebar_button.Model{Label: "Upgrade to pro", Link: "https://www.creative-tim.com/product/material-dashboard-pro?ref=sidebarfree", Outlined: false},
		},
	}

	return component_sidebar.Component{
		Model:  component_sidebar.Model{},
		Top:    sidebarTop,
		Button: sidebarButton,
	}
}

func controller(c example.Context, w http.ResponseWriter, r *http.Request) (mvc.View[example.Context], error) {
	return mvc.ComponentMap[example.Context, model]{
		Template: Html,
		Model:    model{},
		Components: map[string]mvc.Component[example.Context]{
			"sidebar": sidebar(),
		},
	}, nil
}

func Get() runtime.Configuration[*runtime_chi.Runtime[example.Context]] {
	return runtime_chi.WithController(path, http.MethodGet, controller, page_error_500.Controller)
}

func Add(ic inverse.Container) {
	runtime_chi.AddRoute[example.Context](
		ic,
		runtime_chi.AddController(path, http.MethodGet, controller, page_error_500.Controller),
	)
}
