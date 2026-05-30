package component_chat_item

import (
	"embed"
	"html/template"
	"strings"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
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
	Title                string
	IsToolRequest        bool
	IsAgentRequest       bool
	IsToolApprovalNeeded bool
	Message              agent.Message
}

func ViewWithState(message agent.Message, state harness.ExecutionState) render.View {
	switch message.Type {
	case agent.MessageType_ToolRequest:
		{
			askForExecution := state.ToolStates[message.Tool.Id] == harness.ToolState_Requested
			showAgentLink := strings.Contains(message.Tool.Name, "agent")
			return render.Component(
				Html,
				Model{
					Title:                titles[message.Type],
					Message:              message,
					IsToolRequest:        askForExecution || showAgentLink,
					IsToolApprovalNeeded: askForExecution,
					IsAgentRequest:       showAgentLink,
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
					Title:                titles[message.Type],
					Message:              message,
					IsToolRequest:        false,
					IsToolApprovalNeeded: false,
					IsAgentRequest:       false,
				},
				make(map[string]render.View),
				[]render.View{},
			)
		}
	}
}

func View(message agent.Message) render.View {
	switch message.Type {
	case agent.MessageType_ToolRequest:
		{

			return render.Component(
				Html,
				Model{
					Title:          titles[message.Type],
					Message:        message,
					IsToolRequest:  true,
					IsAgentRequest: strings.Contains(message.Tool.Name, "agent"),
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
					Title:          titles[message.Type],
					Message:        message,
					IsToolRequest:  false,
					IsAgentRequest: false,
				},
				make(map[string]render.View),
				[]render.View{},
			)
		}
	}
}
