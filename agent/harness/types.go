package harness

import (
	"strings"

	"github.com/hjwalt/platform/agent"
)

type ToolState int

const (
	ToolState_Requested ToolState = iota
	ToolState_Executed
	ToolState_Rejected
	ToolState_Failed
)

type ExecutionState struct {
	Context        string
	Parent         agent.AgentContext
	ParentToolCall agent.ToolCall
	Messages       []agent.Message
	ToolStates     map[string]ToolState
	Next           agent.Result
	LoadedSkills   map[string]bool
}

func (st ExecutionState) SetContext(ctx string) ExecutionState {
	st.Context = ctx
	return st
}

func (st ExecutionState) SetAgentContext(parent agent.AgentContext, tool agent.ToolCall) ExecutionState {
	st.Parent = parent
	st.ParentToolCall = tool
	return st
}

func (st ExecutionState) SetNext(result agent.Result) ExecutionState {
	st.Next = result
	return st
}

func (st ExecutionState) UpdateToolState(toolName string, state ToolState) ExecutionState {
	if st.ToolStates == nil {
		st.ToolStates = make(map[string]ToolState)
	}
	st.ToolStates[toolName] = state
	return st
}

func (st ExecutionState) AppendMessage(message agent.Message) ExecutionState {
	if st.Messages == nil {
		st.Messages = make([]agent.Message, 0)
	}
	st.Messages = append(st.Messages, message)
	return st
}

func (st ExecutionState) SkillLoaded(skill string) bool {
	if st.LoadedSkills == nil {
		return false
	}
	loaded, present := st.LoadedSkills[strings.ToLower(skill)]
	return present && loaded
}

func (st ExecutionState) AppendSkillLoaded(skill string) ExecutionState {
	if st.LoadedSkills == nil {
		st.LoadedSkills = make(map[string]bool, 0)
	}
	st.LoadedSkills[strings.ToLower(skill)] = true
	return st
}
