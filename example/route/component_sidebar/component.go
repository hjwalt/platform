package component_sidebar

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/example/route/component_sidebar_button"
	"github.com/hjwalt/platform/example/route/component_sidebar_button_list"
	"github.com/hjwalt/platform/example/route/component_sidebar_item"
	"github.com/hjwalt/platform/example/route/component_sidebar_item_header"
	"github.com/hjwalt/platform/example/route/component_sidebar_item_list"
	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Icon   string
	Label  string
	Active bool
}

func View() render.View {
	sidebarTop := component_sidebar_item_list.View(
		component_sidebar_item_list.Model{},
		[]render.View{
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "chart-bar", Label: "Dashboard", Link: "/", Active: false}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "table", Label: "Tables", Link: "/pages/tables.html", Active: false}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "receipt", Label: "Billing", Link: "/billing", Active: false}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "vr-cardboard", Label: "Virtual Reality", Link: "/pages/virtual-reality.html", Active: false}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "bell", Label: "Notifications", Link: "/pages/notifications.html", Active: false}),
			component_sidebar_item_header.View(component_sidebar_item_header.Model{Label: "Account pages"}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "person", Label: "Profile", Link: "/pages/profile.html", Active: false}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "right-to-bracket", Label: "Sign In", Link: "/pages/sign-in.html", Active: false}),
			component_sidebar_item.View(component_sidebar_item.Model{Icon: "user-plus", Label: "Sign Up", Link: "/pages/sign-up.html", Active: false}),
		},
	)

	sidebarButton := component_sidebar_button_list.View(
		component_sidebar_button_list.Model{},
		[]render.View{
			component_sidebar_button.View(component_sidebar_button.Model{Label: "Documentation", Link: "https://www.creative-tim.com/learning-lab/bootstrap/overview/material-dashboard?ref=sidebarfree", Outlined: true}),
		},
	)

	return render.Component(
		Html,
		Model{},
		map[string]render.View{
			"top":    sidebarTop,
			"button": sidebarButton,
		},
		[]render.View{},
	)
}
