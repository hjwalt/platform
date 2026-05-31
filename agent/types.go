package agent

import (
	"context"

	"github.com/google/uuid"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/runtime"
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
	MessageType_Start       MessageType = "FORK"
)

type Message struct {
	Id               string
	Context          string
	Type             MessageType
	Message          string
	ReasoningContent string
	Tool             ToolCall
	Agent            AgentContext
}

type ToolCall struct {
	Id        string
	Name      string
	Arguments string
}

type AgentContext struct {
	ParentContext string
	SystemMessage string
	AllowedTools  []string
}

func NewMessage(
	context string,
	messageType MessageType,
	message string,
	reasoning string,
	tool ToolCall,
) Message {
	return Message{
		Id:               uuid.New().String(),
		Context:          context,
		Type:             messageType,
		Message:          message,
		ReasoningContent: reasoning,
		Tool:             tool,
		Agent:            AgentContext{},
	}
}

func Start(
	context string,
	message string,
	tool ToolCall,
	parent AgentContext,
) Message {
	return Message{
		Id:               uuid.New().String(),
		Context:          context,
		Type:             MessageType_Start,
		Message:          message,
		ReasoningContent: "",
		Tool:             tool,
		Agent:            parent,
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

type LanguageModel interface {
	runtime.Runtime
	Chat(context.Context, []Message, []string) ([]Message, error)
}

var (
	ToolDataFormat = format.Json[ToolCall]()
)
