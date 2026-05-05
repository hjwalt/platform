package harness

import (
	"context"
	"log/slog"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
)

type OpenAiFlow struct {
	Tools map[string]agent.Tool
	Model agent.LanguageModel
}

func (r *OpenAiFlow) Key(ctx context.Context, in agent.Message) (string, error) {
	id := in.Context
	if id == "" {
		id = "DEFAULT"
	}
	return id, nil
}

func (r *OpenAiFlow) Update(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	// setting defaults
	if st.ToolStates == nil {
		st.ToolStates = make(map[string]ToolState)
	}

	slog.Info("updating state", "message", in.Type)
	switch in.Type {
	case agent.MessageType_User, agent.MessageType_ToolResult:
		{
			return r.modelExecute(ctx, in, st)
		}
	case agent.MessageType_ToolRequest:
		{
			return r.toolRequest(ctx, in, st)
		}
	case agent.MessageType_ToolExecute:
		{
			return r.toolExecute(ctx, in, st)
		}
	default:
		{
			return either.Left[ExecutionState, agent.Message](ExecutionState{
				Messages:   append(st.Messages, in),
				ToolStates: st.ToolStates,
				Next:       Result{Messages: []agent.Message{}},
			})
		}
	}
}

func (r *OpenAiFlow) Next(ctx context.Context, in agent.Message, st ExecutionState) (optional.Optional[Result], optional.Optional[agent.Message]) {
	return optional.Of(st.Next), optional.Empty[agent.Message]()
}

func (r *OpenAiFlow) Explode(ctx context.Context, in Result) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	slog.Info("exploding", "len", len(in.Messages))
	return optional.Of(in.Messages), optional.Empty[agent.Message]()
}

func (r *OpenAiFlow) modelExecute(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	allmessages := make([]agent.Message, 0)
	allmessages = append(allmessages, st.Messages...)
	allmessages = append(allmessages, in)

	result, err := r.Model.Chat(context.Background(), allmessages)
	if err != nil {
		return r.updateError(ctx, in, err)
	}

	return either.Left[ExecutionState, agent.Message](ExecutionState{
		Messages:   append(st.Messages, in),
		ToolStates: st.ToolStates,
		Next:       Result{Messages: result},
	})
}

func (r *OpenAiFlow) toolRequest(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	toolData := in.Tool
	if tool, exists := r.Tools[toolData.Name]; exists {
		newToolStates := st.ToolStates
		newToolStates[in.Tool.Name] = ToolState_Requested

		if tool.Auto() {
			return either.Left[ExecutionState, agent.Message](ExecutionState{
				Messages:   append(st.Messages, in),
				ToolStates: newToolStates,
				Next: Result{Messages: []agent.Message{{
					Context: in.Context,
					Type:    agent.MessageType_ToolExecute,
					Message: "execution approved to " + in.Message,
					Tool:    in.Tool,
				}}},
			})
		} else {
			return either.Left[ExecutionState, agent.Message](ExecutionState{
				Messages:   append(st.Messages, in),
				ToolStates: newToolStates,
				Next:       Result{Messages: []agent.Message{}},
			})
		}
	} else {
		newToolStates := st.ToolStates
		newToolStates[in.Tool.Name] = ToolState_Failed

		return either.Left[ExecutionState, agent.Message](ExecutionState{
			Messages:   append(st.Messages, in),
			ToolStates: newToolStates,
			Next: Result{Messages: []agent.Message{{
				Context: in.Context,
				Type:    agent.MessageType_Error,
				Message: "tool " + toolData.Name + " does not exist",
				Tool:    toolData,
			}}},
		})
	}
}

func (r *OpenAiFlow) toolExecute(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	toolData := in.Tool
	if tool, exists := r.Tools[toolData.Name]; exists {
		result, err := tool.Execute(toolData.Arguments)
		if err != nil {
			return r.updateError(ctx, in, err)
		}

		newToolStates := st.ToolStates
		newToolStates[in.Tool.Name] = ToolState_Executed

		return either.Left[ExecutionState, agent.Message](ExecutionState{
			Messages:   append(st.Messages, in),
			ToolStates: newToolStates,
			Next: Result{Messages: []agent.Message{{
				Context: in.Context,
				Type:    agent.MessageType_ToolResult,
				Message: result,
				Tool:    toolData,
			}}},
		})
	} else {
		newToolStates := st.ToolStates
		newToolStates[in.Tool.Name] = ToolState_Failed

		return either.Left[ExecutionState, agent.Message](ExecutionState{
			Messages:   append(st.Messages, in),
			ToolStates: newToolStates,
			Next: Result{Messages: []agent.Message{{
				Context: in.Context,
				Type:    agent.MessageType_Error,
				Message: "tool " + toolData.Name + " does not exist",
				Tool:    toolData,
			}}},
		})
	}
}

func (r *OpenAiFlow) updateError(ctx context.Context, in agent.Message, err error) either.Either[ExecutionState, agent.Message] {
	return either.Right[ExecutionState, agent.Message](agent.Message{
		Context: in.Context,
		Type:    agent.MessageType_Error,
		Message: err.Error(),
		Tool:    in.Tool,
	})
}
