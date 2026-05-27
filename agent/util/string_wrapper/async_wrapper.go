package tool_string_wrapper

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
)

type asyncWrapper[REQ any] struct {
	delegate agent.AsyncTool[REQ]
}

func (t *asyncWrapper[REQ]) Send(ctx context.Context, parent agent.Parent, toolCall agent.ToolCall, stringRequest string) error {
	request, requestParseErr := format.Convert(stringRequest, t.RequestFormat(), t.delegate.RequestFormat())
	if requestParseErr != nil {
		return requestParseErr
	}

	return t.delegate.Send(ctx, parent, toolCall, request)
}

func (t *asyncWrapper[REQ]) Name() string {
	return t.delegate.Name()
}

func (t *asyncWrapper[REQ]) Description() string {
	return t.delegate.Description()
}

func (t *asyncWrapper[REQ]) RequestFormat() format.Format[string] {
	return format.String()
}

func (t *asyncWrapper[REQ]) RequestSchema() *jsonschema.Schema {
	return t.delegate.RequestSchema()
}

func (t *asyncWrapper[REQ]) DescribeRequest(stringRequest string) string {
	request, requestParseErr := format.Convert(stringRequest, t.RequestFormat(), t.delegate.RequestFormat())
	if requestParseErr != nil {
		return "failed to parse request from string to type"
	}
	return t.delegate.DescribeRequest(request)
}

func (t *asyncWrapper[REQ]) Auto() bool {
	return t.delegate.Auto()
}

func StringWrapAsync[REQ any](delegate agent.AsyncTool[REQ]) agent.AsyncToolWrapper {
	return &asyncWrapper[REQ]{
		delegate: delegate,
	}
}
