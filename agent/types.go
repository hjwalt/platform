package agent

import (
	"context"

	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/runtime"
)

const (
	MessageType_User        = "USER"
	MessageType_ToolRequest = "TOOL_REQUEST"
	MessageType_ToolResult  = "TOOL_RESULT"
	MessageType_Agent       = "AGENT"
	MessageType_Error       = "ERROR"
)

type Message struct {
	Type    string
	Message string
	Raw     string
}

type Tool interface {
	Execute(string) (string, error)
}

type ToolData struct {
	Id   string
	Name string
}

type LanguageModel interface {
	runtime.Runtime
	Chat(context.Context, []Message) ([]Message, error)
}

var (
	ToolDataFormat = format.Json[ToolData]()
)
