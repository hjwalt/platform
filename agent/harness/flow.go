package harness

import (
	"context"
	"log/slog"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
)

type Flow struct {
	Tools map[string]agent.Tool
	Model agent.LanguageModel
}

func (r *Flow) Key(ctx context.Context, in agent.Message) (string, error) {
	id := in.Context
	if id == "" {
		id = "DEFAULT"
	}
	return id, nil
}

func (r *Flow) Update(inctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	ctx := logger.WithContext(inctx, "context", in.Context)

	// setting state defaults
	if st.Messages == nil {
		st.Messages = make([]agent.Message, 0)
	}
	if st.ToolStates == nil {
		st.ToolStates = make(map[string]ToolState)
	}

	slog.InfoContext(ctx, "updating state", "message", in.Type)
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

func (r *Flow) Next(ctx context.Context, in agent.Message, st ExecutionState) (optional.Optional[Result], optional.Optional[agent.Message]) {
	// setting state defaults
	if st.Next.Messages == nil {
		st.Next.Messages = make([]agent.Message, 0)
	}

	return optional.Of(st.Next), optional.Empty[agent.Message]()
}

func (r *Flow) Explode(ctx context.Context, in Result) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	slog.Info("exploding", "len", len(in.Messages))
	return optional.Of(in.Messages), optional.Empty[agent.Message]()
}

func (r *Flow) modelExecute(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
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

func (r *Flow) toolRequest(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	toolData := in.Tool
	if tool, exists := r.Tools[toolData.Name]; exists {
		newToolStates := st.ToolStates
		newToolStates[in.Tool.Id] = ToolState_Requested

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
		newToolStates[in.Tool.Id] = ToolState_Failed

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

func (r *Flow) toolExecute(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	toolData := in.Tool

	toolState, exist := st.ToolStates[toolData.Id]
	if !exist {
		toolState = ToolState_Requested
	}

	if toolState != ToolState_Requested {
		return either.Left[ExecutionState, agent.Message](ExecutionState{
			Messages:   append(st.Messages, in),
			ToolStates: st.ToolStates,
			Next: Result{Messages: []agent.Message{{
				Context: in.Context,
				Type:    agent.MessageType_Error,
				Message: "tool " + toolData.Name + " already executed",
				Tool:    toolData,
			}}},
		})
	}

	if tool, exists := r.Tools[toolData.Name]; exists {
		result, err := tool.Execute(toolData.Arguments)
		if err != nil {
			return r.updateError(ctx, in, err)
		}

		newToolStates := st.ToolStates
		newToolStates[in.Tool.Id] = ToolState_Executed

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
		newToolStates[in.Tool.Id] = ToolState_Failed

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

func (r *Flow) updateError(ctx context.Context, in agent.Message, err error) either.Either[ExecutionState, agent.Message] {
	return either.Right[ExecutionState, agent.Message](agent.Message{
		Context: in.Context,
		Type:    agent.MessageType_Error,
		Message: err.Error(),
		Tool:    in.Tool,
	})
}
