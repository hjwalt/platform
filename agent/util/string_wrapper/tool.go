package tool_string_wrapper

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	agent_tool "github.com/hjwalt/platform/agent/tool"
	"github.com/hjwalt/platform/format"
	"github.com/openai/openai-go/v3"
)

type syncToolWrapper[REQ any, RES any] struct {
	delegate agent_tool.Sync[REQ, RES]
}

func (t *syncToolWrapper[REQ, RES]) Apply(stringRequest string) (string, error) {
	request, requestParseErr := format.Convert(stringRequest, t.RequestFormat(), t.delegate.RequestFormat())
	if requestParseErr != nil {
		return "", requestParseErr
	}

	response, responseErr := t.delegate.Apply(request)
	if responseErr != nil {
		return "", responseErr
	}

	return format.Convert(response, t.delegate.ResultFormat(), t.ResultFormat())
}

func (t *syncToolWrapper[REQ, RES]) Name() string {
	return t.delegate.Name()
}

func (t *syncToolWrapper[REQ, RES]) Description() string {
	return t.delegate.Description()
}

func (t *syncToolWrapper[REQ, RES]) RequestFormat() format.Format[string] {
	return format.String()
}

func (t *syncToolWrapper[REQ, RES]) RequestSchema() *jsonschema.Schema {
	return t.delegate.RequestSchema()
}

func (t *syncToolWrapper[REQ, RES]) DescribeRequest(stringRequest string) string {
	request, requestParseErr := format.Convert(stringRequest, t.RequestFormat(), t.delegate.RequestFormat())
	if requestParseErr != nil {
		return "failed to parse request from string to type"
	}
	return t.delegate.DescribeRequest(request)
}

func (t *syncToolWrapper[REQ, RES]) ResultFormat() format.Format[string] {
	return format.String()
}

func (t *syncToolWrapper[REQ, RES]) ResultSchema() *jsonschema.Schema {
	return t.delegate.ResultSchema()
}

func (t *syncToolWrapper[REQ, RES]) DescribeResult(stringResponse string) string {
	response, responseParseErr := format.Convert(stringResponse, t.ResultFormat(), t.delegate.ResultFormat())
	if responseParseErr != nil {
		return "failed to parse result from string to type"
	}
	return t.delegate.DescribeResult(response)
}

func (t *syncToolWrapper[REQ, RES]) Auto() bool {
	return t.delegate.Auto()
}

func StringWrapSync[REQ any, RES any](delegate agent_tool.Sync[REQ, RES]) agent_tool.SyncWrapper {
	return &syncToolWrapper[REQ, RES]{
		delegate: delegate,
	}
}

// TODO: remove the need of these
type ToolWrapper struct {
	tool agent_tool.SyncWrapper
}

func (t *ToolWrapper) Name() string {
	return t.tool.Name()
}

func (t *ToolWrapper) Description() string {
	return t.tool.Description()
}

func (t *ToolWrapper) Schema() openai.ChatCompletionToolUnionParam {
	return llm.FromJsonSchema(t.Name(), t.Description(), t.tool.RequestSchema())
}

func (t *ToolWrapper) Execute(input string) (string, error) {
	response, internalErr := t.tool.Apply(input)
	if internalErr != nil {
		return "", internalErr
	}
	return t.tool.DescribeResult(response), nil
}

func (t *ToolWrapper) Request(input string) (string, error) {
	return t.tool.DescribeRequest(input), nil
}

func (t *ToolWrapper) Auto() bool {
	return false
}

func WrapSync[REQ any, RES any](tool agent_tool.Sync[REQ, RES]) agent.Tool {
	return &ToolWrapper{
		tool: StringWrapSync[REQ, RES](tool),
	}
}

func WrapWrapper(tool agent_tool.SyncWrapper) agent.Tool {
	return &ToolWrapper{
		tool: tool,
	}
}
