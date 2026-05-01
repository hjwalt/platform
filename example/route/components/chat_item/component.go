package chat_item

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Message string
}

func View(message agent.Message) render.View {
	return render.Component(
		Html,
		Model{
			Message: message.Message,
		},
		make(map[string]render.View),
		[]render.View{},
	)
}
