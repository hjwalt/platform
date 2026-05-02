package harness

import (
	"context"
	"errors"
	"log/slog"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/type/optional"
	"github.com/kultivator-consulting/goharmony"
)

type OpenAiFlow struct {
	Store rag.Store
	Tools map[string]agent.Tool
	Model agent.LanguageModel
}

func (r *OpenAiFlow) Handle(ctx context.Context, in agent.Message) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
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
			if storeErr := r.Store.Add(in.Context, []agent.Message{in}); storeErr != nil {
				slog.Error("failed to store error message", "error", storeErr)
			}
			return optional.Empty[[]agent.Message](), optional.Empty[agent.Message]()
		}
	case agent.MessageType_Agent:
		{
			if storeErr := r.Store.Add(in.Context, []agent.Message{in}); storeErr != nil {
				return r.writeError(ctx, in, storeErr)
			}

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

func (r *OpenAiFlow) runModel(ctx context.Context, in agent.Message) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	id := in.Context
	if id == "" {
		id = "DEFAULT"
	}

	storedMessages, getErr := r.Store.GetAll(id)
	if getErr != nil {
		return r.writeError(ctx, in, getErr)
	}

	if storeErr := r.Store.Add(id, []agent.Message{in}); storeErr != nil {
		return r.writeError(ctx, in, storeErr)
	}

	allmessages := make([]agent.Message, 0)
	allmessages = append(allmessages, storedMessages...)
	allmessages = append(allmessages, in)

	result, err := r.Model.Chat(context.Background(), allmessages)
	if err != nil {
		return r.writeError(ctx, in, err)
	}

	return optional.Of(result), optional.Empty[agent.Message]()
}

func (r *OpenAiFlow) toolRequest(ctx context.Context, in agent.Message) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	messages := make([]agent.Message, 0)

	toolData := in.Tool
	if tool, exists := r.Tools[toolData.Name]; exists {
		result, err := tool.Execute(toolData.Arguments)
		if err != nil {
			return r.writeError(ctx, in, err)
		}

		messages = append(messages, agent.Message{
			Context: in.Context,
			Type:    agent.MessageType_ToolResult,
			Message: result,
			Tool:    toolData,
		})
		return optional.Of(messages), optional.Empty[agent.Message]()
	}

	return r.writeError(ctx, in, errors.New("tool "+toolData.Name+" does not exist"))
}

func (r *OpenAiFlow) writeError(ctx context.Context, in agent.Message, err error) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	return optional.Empty[[]agent.Message](), optional.Of(agent.Message{
		Context: in.Context,
		Type:    agent.MessageType_Error,
		Message: err.Error(),
		Tool:    in.Tool,
	})
}
