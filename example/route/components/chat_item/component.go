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

var titles = map[agent.MessageType]string{
	agent.MessageType_Agent:       "Response",
	agent.MessageType_User:        "User",
	agent.MessageType_ToolRequest: "Tool Request",
	agent.MessageType_ToolResult:  "Tool Result",
	agent.MessageType_ToolExecute: "Tool Execute",
	agent.MessageType_Error:       "Error",
	agent.MessageType_System:      "System",
}

type Model struct {
	Title         string
	Message       string
	IsToolRequest bool
}

func View(c example.Context, message agent.Message) render.View {
	switch message.Type {
	case agent.MessageType_ToolRequest:
		{
			return render.Component(
				Html,
				Model{
					Title:         titles[message.Type],
					Message:       message.Message,
					IsToolRequest: true,
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
					Title:         titles[message.Type],
					Message:       message.Message,
					IsToolRequest: false,
				},
				make(map[string]render.View),
				[]render.View{},
			)
		}
	}
}
