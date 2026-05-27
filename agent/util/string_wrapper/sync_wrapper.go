package tool_string_wrapper

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
)

type syncToolWrapper[REQ any, RES any] struct {
	delegate agent.SyncTool[REQ, RES]
}

func (t *syncToolWrapper[REQ, RES]) Apply(ctx context.Context, stringRequest string) (string, error) {
	request, requestParseErr := format.Convert(stringRequest, t.RequestFormat(), t.delegate.RequestFormat())
	if requestParseErr != nil {
		return "", requestParseErr
	}

	response, responseErr := t.delegate.Apply(ctx, request)
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

func StringWrapSync[REQ any, RES any](delegate agent.SyncTool[REQ, RES]) agent.SyncToolWrapper {
	return &syncToolWrapper[REQ, RES]{
		delegate: delegate,
	}
}
