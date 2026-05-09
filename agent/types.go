package agent

import (
	"context"

	"github.com/google/uuid"
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
	MessageType_ToolExecute MessageType = "TOOL_EXECUTE"
	MessageType_Agent       MessageType = "AGENT"
	MessageType_Error       MessageType = "ERROR"
)

type Message struct {
	Id      string
	Context string
	Type    MessageType
	Message string
	Tool    ToolCall
}

func NewMessage(
	context string,
	messageType MessageType,
	message string,
	tool ToolCall,
) Message {
	return Message{
		Id:      uuid.New().String(),
		Context: context,
		Type:    messageType,
		Message: message,
		Tool:    tool,
	}
}

type Result struct {
	Id       string
	Messages []Message
}

func NewResult(messages []Message) Result {
	return Result{
		Id:       uuid.New().String(),
		Messages: messages,
	}
}

func SingleResult(message Message) Result {
	return Result{
		Id:       uuid.New().String(),
		Messages: []Message{message},
	}
}

func EmptyResult() Result {
	return Result{
		Id:       uuid.New().String(),
		Messages: make([]Message, 0),
	}
}

type Tool interface {
	Name() string
	Description() string
	Schema() openai.ChatCompletionToolUnionParam
	Execute(string) (string, error)
	Request(string) (string, error)
	Auto() bool
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
