package harness

import (
	"context"
	"log/slog"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
)

type FlowMetadata struct {
}

func (r *FlowMetadata) Key(ctx context.Context, in agent.Message) (string, error) {
	id := in.Context
	if id == "" {
		id = "DEFAULT"
	}
	return id, nil
}

func (r *FlowMetadata) ResultMetadata(ctx context.Context, pref flow.Metadata, value agent.Result) flow.Metadata {
	return flow.Metadata{
		Id:       value.Id,
		Group:    pref.Group,
		Attempt:  0,
		Sequence: pref.Sequence + 1,
		Source:   "AGENT_HARNESS",
	}
}

func (r *FlowMetadata) MessageMetadata(ctx context.Context, pref flow.Metadata, value agent.Message) flow.Metadata {
	return flow.Metadata{
		Id:       value.Id,
		Group:    value.Context,
		Attempt:  0,
		Sequence: pref.Sequence + 1,
		Source:   "AGENT_HARNESS",
	}
}

type Flow struct {
	Tools agent.ToolContainer
	Model agent.LanguageModel
}

func (r *Flow) Update(inctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	ctx := logger.WithContext(inctx, "context", in.Context)

	// set context
	if st.Context == "" {
		st = st.SetContext(in.Context)
	}

	// reset next
	st = st.SetNext(agent.EmptyResult())
	st = st.AppendMessage(in)

	slog.InfoContext(ctx, "updating state", "message", in.Type)
	switch in.Type {
	case agent.MessageType_Start:
		{
			return r.startAgent(in, st)
		}
	case agent.MessageType_Agent:
		{
			return r.mergeMessage(in, st)
		}
	case agent.MessageType_User:
		{
			return r.modelExecute(ctx, in, st)
		}
	case agent.MessageType_ToolRequest:
		{
			return r.toolRequest(ctx, in, st)
		}
	case agent.MessageType_ToolResult:
		{
			return r.toolResult(ctx, in, st)
		}
	case agent.MessageType_ToolExecute:
		{
			return r.toolExecute(ctx, in, st)
		}
	default:
		{
			return either.Left[ExecutionState, agent.Message](st)
		}
	}
}

func (r *Flow) Next(ctx context.Context, in agent.Message, st ExecutionState) (optional.Optional[agent.Result], optional.Optional[agent.Message]) {
	slog.InfoContext(ctx, "next", "len", len(st.Next.Messages))
	return optional.Of(st.Next), optional.Empty[agent.Message]()
}

func (r *Flow) Explode(ctx context.Context, in agent.Result) (optional.Optional[[]agent.Message], optional.Optional[agent.Message]) {
	slog.InfoContext(ctx, "exploding", "len", len(in.Messages))
	return optional.Of(in.Messages), optional.Empty[agent.Message]()
}

func (r *Flow) startAgent(in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	st = st.AppendMessage(agent.NewMessage(
		in.Context,
		agent.MessageType_System,
		in.Agent.SystemMessage,
		in.ReasoningContent,
		in.Tool,
	))
	st = st.SetAgentContext(in.Agent, in.Tool)
	st = st.SetNext(agent.SingleResult(agent.NewMessage(
		in.Context,
		agent.MessageType_User,
		in.Message,
		in.ReasoningContent,
		agent.ToolCall{},
	)))
	return either.Left[ExecutionState, agent.Message](st)
}

func (r *Flow) mergeMessage(in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	if st.Parent.ParentContext != "" {
		st = st.SetNext(agent.SingleResult(agent.NewMessage(
			st.Parent.ParentContext,
			agent.MessageType_ToolResult,
			in.Message,
			in.ReasoningContent,
			st.ParentToolCall,
		)))
	}

	return either.Left[ExecutionState, agent.Message](st)
}

func (r *Flow) modelExecute(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	// TODO: tidy up messages based on tool execution result

	if result, err := r.Model.Chat(context.Background(), st.Messages, st.Parent.AllowedTools); err != nil {
		st = st.SetNext(agent.SingleResult(agent.NewMessage(
			in.Context,
			agent.MessageType_Error,
			err.Error(),
			in.ReasoningContent,
			in.Tool,
		)))
	} else {
		st = st.SetNext(agent.NewResult(result))
		slog.InfoContext(ctx, "next", "len", len(result))
	}

	return either.Left[ExecutionState, agent.Message](st)
}

func (r *Flow) toolRequest(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	toolData := in.Tool
	if exists := r.Tools.Exists(toolData); exists {
		st = st.UpdateToolState(in.Tool.Id, ToolState_Requested)
		if r.Tools.Auto(toolData) {
			st = st.SetNext(agent.SingleResult(agent.NewMessage(
				in.Context,
				agent.MessageType_ToolExecute,
				"execution approved to "+in.Message,
				in.ReasoningContent,
				in.Tool,
			)))
		}
	} else {
		st = st.UpdateToolState(in.Tool.Id, ToolState_Failed)
		st = st.SetNext(agent.SingleResult(agent.NewMessage(
			in.Context,
			agent.MessageType_Error,
			"tool "+toolData.Name+" does not exist",
			in.ReasoningContent,
			toolData,
		)))
	}

	return either.Left[ExecutionState, agent.Message](st)
}

func (r *Flow) toolExecute(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	toolData := in.Tool

	toolState, exist := st.ToolStates[toolData.Id]
	if !exist {
		toolState = ToolState_Requested
	}

	if toolState != ToolState_Requested {
		st = st.SetNext(agent.SingleResult(agent.NewMessage(
			in.Context,
			agent.MessageType_Error,
			"tool "+toolData.Name+" already executed",
			in.ReasoningContent,
			toolData,
		)))
	} else if result, toolError := r.Tools.Execute(ctx, in, toolData); toolError == nil {
		st = st.UpdateToolState(in.Tool.Id, ToolState_Executed)
		if result.IsPresent() {
			st = st.SetNext(agent.SingleResult(agent.NewMessage(
				in.Context,
				agent.MessageType_ToolResult,
				result.Get(),
				in.ReasoningContent,
				toolData,
			)))
		}
	} else {
		st = st.UpdateToolState(in.Tool.Id, ToolState_Failed)
		st = st.SetNext(agent.SingleResult(agent.NewMessage(
			in.Context,
			agent.MessageType_Error,
			toolError.Error(),
			in.ReasoningContent,
			toolData,
		)))
	}

	return either.Left[ExecutionState, agent.Message](st)
}

func (r *Flow) toolResult(ctx context.Context, in agent.Message, st ExecutionState) either.Either[ExecutionState, agent.Message] {
	st = st.UpdateToolState(in.Tool.Id, ToolState_Executed)
	return r.modelExecute(ctx, in, st)
}
