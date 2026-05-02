package agent

import (
	"context"

	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/runtime"
	"github.com/openai/openai-go/v3"
)

type MessageType string

const (
	MessageType_System      MessageType = "SYSTEM"
	MessageType_User        MessageType = "USER"
	MessageType_ToolRequest MessageType = "TOOL_REQUEST"
	MessageType_ToolResult  MessageType = "TOOL_RESULT"
	MessageType_Agent       MessageType = "AGENT"
	MessageType_Error       MessageType = "ERROR"
)

type Message struct {
	Context string
	Type    MessageType
	Message string
	Tool    ToolCall
}

type Tool interface {
	Name() string
	Description() string
	Schema() openai.ChatCompletionToolUnionParam
	Execute(string) (string, error)
}

type ToolCall struct {
	Id        string
	Name      string
	Arguments string
}

type LanguageModel interface {
	runtime.Runtime
	Chat(context.Context, []Message) ([]Message, error)
}

var (
	ToolDataFormat = format.Json[ToolCall]()
)
