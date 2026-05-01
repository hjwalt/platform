package harness

import (
	"context"
	"log/slog"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/type/optional"
	"github.com/kultivator-consulting/goharmony"
)

type OpenAiFlow[C context.Context] struct {
	Tools map[string]agent.Tool
	Model agent.LanguageModel
}

func (r *OpenAiFlow[C]) Handle(ctx C, in agent.Message) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	slog.Info("new message", "message", in)
	switch in.Type {
	case agent.MessageType_User, agent.MessageType_ToolResult:
		{
			return r.runModel(ctx, in)
		}
	case agent.MessageType_ToolRequest:
		{
			return r.toolRequest(ctx, in)
		}
	case agent.MessageType_Error:
		{
			slog.Info("error", "message", in)
			return optional.Empty[[]agent.Message](), optional.Empty[agent.Message]()
		}
	case agent.MessageType_Agent:
		{
			slog.Info("agent", "message", in.Message)
			parser := goharmony.NewParser()

			messages, _ := parser.ParseResponse(in.Message)
			for _, msg := range messages {
				slog.Info("agent", "channel", msg.Channel, "content", msg.Content)
			}

			return optional.Empty[[]agent.Message](), optional.Empty[agent.Message]()
		}
	}
	return optional.Empty[[]agent.Message](), optional.Empty[agent.Message]()
}

func (r *OpenAiFlow[C]) runModel(ctx C, in agent.Message) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	result, err := r.Model.Chat(context.Background(), []agent.Message{in})
	if err != nil {
		return optional.Empty[[]agent.Message](), optional.Of(agent.Message{
			Type:    agent.MessageType_Error,
			Message: err.Error(),
			Raw:     "",
		})
	}
	return optional.Of(result), optional.Empty[agent.Message]()
}

func (r *OpenAiFlow[C]) toolRequest(ctx C, in agent.Message) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	messages := make([]agent.Message, 0)

	rawMessage, unmarshallErr := llm.OpenAiMessageFormat.Unmarshal([]byte(in.Raw))
	if unmarshallErr != nil {
		panic(unmarshallErr)
	}
	toolCalls := rawMessage.ToolCalls
	for _, toolCall := range toolCalls {
		toolData := agent.ToolData{
			Id:   toolCall.ID,
			Name: toolCall.Function.Name,
		}

		if tool, exists := r.Tools[toolData.Name]; exists {
			result, err := tool.Execute(toolCall.Function.Arguments)
			if err != nil {
				return optional.Empty[[]agent.Message](), optional.Of(agent.Message{
					Type:    agent.MessageType_Error,
					Message: err.Error(),
					Raw:     "",
				})
			}

			toolDataRaw, toolMarshallErr := agent.ToolDataFormat.Marshal(toolData)
			if toolMarshallErr != nil {
				return optional.Empty[[]agent.Message](), optional.Of(agent.Message{
					Type:    agent.MessageType_Error,
					Message: toolMarshallErr.Error(),
					Raw:     "",
				})
			}

			messages = append(messages, agent.Message{
				Type:    agent.MessageType_ToolResult,
				Message: result,
				Raw:     string(toolDataRaw),
			})
		}
	}

	return optional.Of(messages), optional.Empty[agent.Message]()
}
