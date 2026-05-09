package harness

import (
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
	Messages   []agent.Message
	ToolStates map[string]ToolState
	Next       agent.Result
}
