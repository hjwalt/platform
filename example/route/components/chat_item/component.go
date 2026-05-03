package chat_item

import (
	"embed"
	"html/template"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Message string
}

func View(c example.Context, message agent.Message) render.View {
	switch message.Type {
	case agent.MessageType_ToolRequest:
		{
			if tool, exists := c.Tool[message.Tool.Name]; exists {
				if toolView, err := tool.RequestView(message); err == nil {
					return toolView
				} else {
					return render.Component(
						Html,
						Model{
							Message: "failed to render tool request",
						},
						make(map[string]render.View),
						[]render.View{},
					)
				}
			}
			return render.Component(
				Html,
				Model{
					Message: "tool missing to render",
				},
				make(map[string]render.View),
				[]render.View{},
			)
		}
	default:
		{
			return render.Component(
				Html,
				Model{
					Message: message.Message,
				},
				make(map[string]render.View),
				[]render.View{},
			)
		}
	}
}
