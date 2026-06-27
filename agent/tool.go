package agent

import (
	"context"

	"github.com/cohesion-org/deepseek-go"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/type/optional"
	"github.com/openai/openai-go/v3"
)

type BasicToolDefinition[REQ any] interface {
	Name() string
	Description() string
	RequestFormat() format.Format[REQ]
	RequestSchema() *jsonschema.Schema
	DescribeRequest(REQ) string
	Auto() bool
}

type SyncTool[REQ any, RES any] interface {
	BasicToolDefinition[REQ]
	Apply(context.Context, REQ) (RES, error)
	ResultFormat() format.Format[RES]
	ResultSchema() *jsonschema.Schema
	DescribeResult(RES) string
}

type AsyncTool[REQ any] interface {
	BasicToolDefinition[REQ]
	Send(context.Context, AgentContext, ToolCall, REQ) error
}

type SyncToolWrapper SyncTool[string, string]

type AsyncToolWrapper AsyncTool[string]

type ToolContainer interface {
	// Register
	AddSync(SyncToolWrapper)
	AddAsync(AsyncToolWrapper)

	// Execution behaviour
	Execute(context.Context, Message, ToolCall) (optional.Optional[string], error)
	DescribeRequest(ToolCall) (string, error)
	Exists(ToolCall) bool
	Auto(ToolCall) bool

	// ToolParams
	OpenAiParamsFiltered([]string) []openai.ChatCompletionToolUnionParam
	DeepSeekParams([]string) []deepseek.Tool
}

// common structure for AGENTS.md and SKILLS.md
type Instruction struct {
	Name          string
	Description   string
	License       optional.Optional[string]
	Compatibility optional.Optional[string]
	AllowedTools  []string
	Metadata      map[string]string
	Body          string
}

type SkillContainer interface {
	Add(Instruction)
	Get(string) (Instruction, bool)
}
