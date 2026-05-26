package agent

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/format"
	"github.com/openai/openai-go/v3"
)

type ToolDefinition[REQ any, RES any] interface {
	Name() string
	Description() string
	RequestFormat() format.Format[REQ]
	RequestSchema() *jsonschema.Schema
	DescribeRequest(REQ) string
	ResultFormat() format.Format[RES]
	ResultSchema() *jsonschema.Schema
	DescribeResult(RES) string
	Auto() bool
}

type SyncTool[REQ any, RES any] interface {
	ToolDefinition[REQ, RES]
	Apply(context.Context, REQ) (RES, error)
}

type AsyncTool[REQ any, RES any] interface {
	ToolDefinition[REQ, RES]
	Send(context.Context, REQ) error
}

type SyncToolWrapper SyncTool[string, string]

type AsyncToolWrapper AsyncTool[string, string]

type ToolContainer interface {
	// Register
	AddSync(SyncToolWrapper)
	GetSync() map[string]SyncToolWrapper

	// Execution behaviour
	Execute(context.Context, ToolCall) (string, error)
	DescribeRequest(ToolCall) (string, error)
	Exists(ToolCall) bool
	Auto(ToolCall) bool

	// OpenAi
	OpenAiParams() []openai.ChatCompletionToolUnionParam
}
