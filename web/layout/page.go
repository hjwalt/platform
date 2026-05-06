package layout

import (
	"embed"

	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

func Page(content render.View) render.View {
	return render.Component(
		render.Embedded(files, "page.html"),
		"",
		map[string]render.View{
			"content": content,
		},
		[]render.View{},
	)
}

func Dashboard(sidebar render.View, content render.View) render.View {
	return render.Component(
		render.Embedded(files, "dashboard.html"),
		"",
		map[string]render.View{
			"sidebar": sidebar,
			"content": content,
		},
		[]render.View{},
	)
}
