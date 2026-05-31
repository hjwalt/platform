package llm

import (
	"context"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
)

func createDeepSeek(config ModelConfig, tools agent.ToolContainer) agent.LanguageModel {
	return &deepSeekModel{
		Model:  config.Model,
		Secret: config.Secret,
		Tools:  tools,
	}
}

type deepSeekModel struct {
	Model  string
	Secret string
	Tools  agent.ToolContainer
	client *deepseek.Client
}

func (r *deepSeekModel) Start() error {
	client := deepseek.NewClient(r.Secret)
	r.client = client

	return nil
}

func (r *deepSeekModel) Stop() {
}

func (r *deepSeekModel) Chat(ctx context.Context, messages []agent.Message, allowedTools []string) ([]agent.Message, error) {
	modelMessage := make([]deepseek.ChatCompletionMessage, 0)
	for _, message := range messages {
		switch message.Type {
		case agent.MessageType_System:
			{
				modelMessage = append(modelMessage, DeepSeekSystemMessage(message))
			}
		case agent.MessageType_User:
			{
				modelMessage = append(modelMessage, DeepSeekUserMessage(message))
			}
		case agent.MessageType_ToolRequest:
			{
				modelMessage = append(modelMessage, DeepSeekToolRequestMessage(message))
			}
		case agent.MessageType_ToolResult:
			{
				modelMessage = append(modelMessage, DeepSeekToolResultMessage(message))
			}
		}
	}

	params := &deepseek.ChatCompletionRequest{
		Messages: modelMessage,
		Model:    r.Model,
		Tools:    r.Tools.DeepSeekParams(allowedTools),
	}

	completion, err := r.client.CreateChatCompletion(ctx, params)
	if err != nil {
		return []agent.Message{agent.NewMessage(
			messages[0].Context,
			agent.MessageType_Error,
			err.Error(),
			"",
			agent.ToolCall{},
		)}, err
	}

	outputMessages := make([]agent.Message, 0)
	for _, choice := range completion.Choices {
		switch choice.FinishReason {
		case "stop":
			{
				outputMessages = append(outputMessages, agent.NewMessage(
					messages[0].Context,
					agent.MessageType_Agent,
					choice.Message.Content,
					choice.Message.ReasoningContent,
					agent.ToolCall{},
				))
			}
		case "tool_calls":
			{
				for _, toolCall := range choice.Message.ToolCalls {
					toolData := agent.ToolCall{
						Id:        toolCall.ID,
						Name:      toolCall.Function.Name,
						Arguments: toolCall.Function.Arguments,
					}

					if toolRequestMessage, messageErr := r.Tools.DescribeRequest(toolData); messageErr == nil {
						outputMessages = append(outputMessages, agent.NewMessage(
							messages[0].Context,
							agent.MessageType_ToolRequest,
							toolRequestMessage,
							choice.Message.ReasoningContent,
							toolData,
						))
					} else {
						// TODO: do something with error
					}
				}
			}
		}
	}

	return outputMessages, nil
}

func DeepSeekToolFromJsonSchema(name string, description string, toolSchema *jsonschema.Schema) deepseek.Tool {
	unmarshalled, _ := format.Convert(toolSchema, schemaFormat, deepseekFormat)
	return deepseek.Tool{
		Type: "function",
		Function: deepseek.Function{
			Name:        name,
			Description: description,
			Parameters:  &unmarshalled,
		},
	}
}

func DeepSeekSystemMessage(message agent.Message) deepseek.ChatCompletionMessage {
	return deepseek.ChatCompletionMessage{
		Role:             deepseek.ChatMessageRoleSystem,
		Content:          message.Message,
		ReasoningContent: message.ReasoningContent,
	}
}

func DeepSeekUserMessage(message agent.Message) deepseek.ChatCompletionMessage {
	return deepseek.ChatCompletionMessage{
		Role:             deepseek.ChatMessageRoleUser,
		Content:          message.Message,
		ReasoningContent: message.ReasoningContent,
	}
}

func DeepSeekToolRequestMessage(message agent.Message) deepseek.ChatCompletionMessage {
	return deepseek.ChatCompletionMessage{
		Role:             deepseek.ChatMessageRoleAssistant,
		Content:          message.Message + " with tool call id " + message.Tool.Id,
		ReasoningContent: message.ReasoningContent,
		ToolCalls: []deepseek.ToolCall{
			{
				ID:   message.Tool.Id,
				Type: "function",
				Function: deepseek.ToolCallFunction{
					Name:      message.Tool.Name,
					Arguments: message.Tool.Arguments,
				},
			},
		},
	}
}

func DeepSeekToolResultMessage(message agent.Message) deepseek.ChatCompletionMessage {
	return deepseek.ChatCompletionMessage{
		Role:             deepseek.ChatMessageRoleTool,
		Content:          message.Message,
		ReasoningContent: message.ReasoningContent,
		ToolCallID:       message.Tool.Id,
	}
}

var (
	deepseekFormat = format.Json[deepseek.FunctionParameters]()
)
