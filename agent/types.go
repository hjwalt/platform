package agent

import "github.com/hjwalt/platform/format"

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

type ToolData struct {
	Id string
}

var (
	ToolDataFormat = format.Json[ToolData]()
)
