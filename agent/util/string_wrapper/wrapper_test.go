package tool_string_wrapper_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Test types and mocks
// ---------------------------------------------------------------------------

type wrapperReq struct {
	Term string `json:"term"`
}

type wrapperRes struct {
	Summary string `json:"summary"`
}

type mockSyncTool struct {
	name          string
	description   string
	auto          bool
	requestSchema *jsonschema.Schema
	resultSchema  *jsonschema.Schema
	applyFn       func(ctx context.Context, req wrapperReq) (wrapperRes, error)
	describeReqFn func(req wrapperReq) string
	describeResFn func(res wrapperRes) string
	applyCalls    int
	lastReq       wrapperReq
}

func (m *mockSyncTool) Name() string                             { return m.name }
func (m *mockSyncTool) Description() string                      { return m.description }
func (m *mockSyncTool) Auto() bool                               { return m.auto }
func (m *mockSyncTool) RequestSchema() *jsonschema.Schema        { return m.requestSchema }
func (m *mockSyncTool) ResultSchema() *jsonschema.Schema         { return m.resultSchema }
func (m *mockSyncTool) RequestFormat() format.Format[wrapperReq] { return format.Json[wrapperReq]() }
func (m *mockSyncTool) ResultFormat() format.Format[wrapperRes]  { return format.Json[wrapperRes]() }

func (m *mockSyncTool) Apply(ctx context.Context, req wrapperReq) (wrapperRes, error) {
	m.applyCalls++
	m.lastReq = req
	if m.applyFn != nil {
		return m.applyFn(ctx, req)
	}
	return wrapperRes{Summary: "default"}, nil
}

func (m *mockSyncTool) DescribeRequest(req wrapperReq) string {
	if m.describeReqFn != nil {
		return m.describeReqFn(req)
	}
	return "req:" + req.Term
}

func (m *mockSyncTool) DescribeResult(res wrapperRes) string {
	if m.describeResFn != nil {
		return m.describeResFn(res)
	}
	return "res:" + res.Summary
}

type mockAsyncTool struct {
	name          string
	description   string
	auto          bool
	requestSchema *jsonschema.Schema
	sendFn        func(ctx context.Context, parent agent.AgentContext, call agent.ToolCall, req wrapperReq) error
	describeReqFn func(req wrapperReq) string
	sendCalls     int
	lastParent    agent.AgentContext
	lastCall      agent.ToolCall
	lastReq       wrapperReq
}

func (m *mockAsyncTool) Name() string                      { return m.name }
func (m *mockAsyncTool) Description() string               { return m.description }
func (m *mockAsyncTool) Auto() bool                        { return m.auto }
func (m *mockAsyncTool) RequestSchema() *jsonschema.Schema { return m.requestSchema }
func (m *mockAsyncTool) RequestFormat() format.Format[wrapperReq] {
	return format.Json[wrapperReq]()
}

func (m *mockAsyncTool) Send(ctx context.Context, parent agent.AgentContext, call agent.ToolCall, req wrapperReq) error {
	m.sendCalls++
	m.lastParent = parent
	m.lastCall = call
	m.lastReq = req
	if m.sendFn != nil {
		return m.sendFn(ctx, parent, call, req)
	}
	return nil
}

func (m *mockAsyncTool) DescribeRequest(req wrapperReq) string {
	if m.describeReqFn != nil {
		return m.describeReqFn(req)
	}
	return "async-req:" + req.Term
}

func newSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	schema, err := jsonschema.For[wrapperReq](&jsonschema.ForOptions{})
	require.NoError(t, err)
	return schema
}

// ---------------------------------------------------------------------------
// StringWrapSync
// ---------------------------------------------------------------------------

func TestStringWrapSyncSatisfiesSyncToolWrapperInterface(t *testing.T) {
	assert := assert.New(t)

	var wrapper agent.SyncToolWrapper = tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](&mockSyncTool{name: "t"})
	require.NotNil(t, wrapper)
	assert.Equal("t", wrapper.Name())
}

func TestStringWrapSyncApplyRoundTripsStringToTypedToResult(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockSyncTool{
		name: "sync_tool",
		applyFn: func(ctx context.Context, req wrapperReq) (wrapperRes, error) {
			return wrapperRes{Summary: "out-" + req.Term}, nil
		},
	}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	result, err := wrapper.Apply(context.Background(), `{"term":"go"}`)

	assert.NoError(err)
	assert.Equal(`{"summary":"out-go"}`, result)
	assert.Equal(1, delegate.applyCalls)
	assert.Equal(wrapperReq{Term: "go"}, delegate.lastReq)

	// the result string converts back into the typed result format
	back, backErr := format.Convert(result, format.String(), format.Json[wrapperRes]())
	assert.NoError(backErr)
	assert.Equal(wrapperRes{Summary: "out-go"}, back)
}

func TestStringWrapSyncApplyParseError(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockSyncTool{name: "sync_tool"}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	_, err := wrapper.Apply(context.Background(), "not json")

	require.Error(t, err)
	assert.Equal(0, delegate.applyCalls, "delegate must not run when the request fails to parse")
}

func TestStringWrapSyncApplyDelegateError(t *testing.T) {
	assert := assert.New(t)

	sentinel := errors.New("apply boom")
	delegate := &mockSyncTool{name: "sync_tool", applyFn: func(ctx context.Context, req wrapperReq) (wrapperRes, error) { return wrapperRes{}, sentinel }}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	_, err := wrapper.Apply(context.Background(), `{"term":"go"}`)

	assert.ErrorIs(err, sentinel)
}

func TestStringWrapSyncDelegatesMetadata(t *testing.T) {
	assert := assert.New(t)

	schema := newSchema(t)
	delegate := &mockSyncTool{
		name:          "meta_tool",
		description:   "does metadata",
		auto:          true,
		requestSchema: schema,
		resultSchema:  schema,
	}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	assert.Equal("meta_tool", wrapper.Name())
	assert.Equal("does metadata", wrapper.Description())
	assert.True(wrapper.Auto())
	assert.Equal(schema, wrapper.RequestSchema())
	assert.Equal(schema, wrapper.ResultSchema())
	assert.Equal("req:go", wrapper.DescribeRequest(`{"term":"go"}`))
	assert.Equal("res:sum", wrapper.DescribeResult(`{"summary":"sum"}`))
}

func TestStringWrapSyncDescribeRequestParseFailureMessage(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockSyncTool{name: "sync_tool", describeReqFn: func(req wrapperReq) string { return "unexpected" }}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	assert.Equal("failed to parse request from string to type", wrapper.DescribeRequest("nope"))
}

func TestStringWrapSyncDescribeResultParseFailureMessage(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockSyncTool{name: "sync_tool", describeResFn: func(res wrapperRes) string { return "unexpected" }}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	assert.Equal("failed to parse result from string to type", wrapper.DescribeResult("nope"))
}

func TestStringWrapSyncDescribeRequestReceivesParsedRequest(t *testing.T) {
	assert := assert.New(t)

	var got wrapperReq
	delegate := &mockSyncTool{
		name:          "sync_tool",
		describeReqFn: func(req wrapperReq) string { got = req; return "described" },
	}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	assert.Equal("described", wrapper.DescribeRequest(`{"term":"search"}`))
	assert.Equal(wrapperReq{Term: "search"}, got)
}

func TestStringWrapSyncDescribeResultReceivesParsedResult(t *testing.T) {
	assert := assert.New(t)

	var got wrapperRes
	delegate := &mockSyncTool{
		name:          "sync_tool",
		describeResFn: func(res wrapperRes) string { got = res; return "described" },
	}
	wrapper := tool_string_wrapper.StringWrapSync[wrapperReq, wrapperRes](delegate)

	assert.Equal("described", wrapper.DescribeResult(`{"summary":"s"}`))
	assert.Equal(wrapperRes{Summary: "s"}, got)
}

// ---------------------------------------------------------------------------
// StringWrapAsync
// ---------------------------------------------------------------------------

func TestStringWrapAsyncSatisfiesAsyncToolWrapperInterface(t *testing.T) {
	assert := assert.New(t)

	var wrapper agent.AsyncToolWrapper = tool_string_wrapper.StringWrapAsync[wrapperReq](&mockAsyncTool{name: "a"})
	require.NotNil(t, wrapper)
	assert.Equal("a", wrapper.Name())
}

func TestStringWrapAsyncSendConvertsAndForwards(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockAsyncTool{name: "async_tool"}
	wrapper := tool_string_wrapper.StringWrapAsync[wrapperReq](delegate)

	parent := agent.AgentContext{ParentContext: "parent-ctx", SystemMessage: "sys"}
	call := agent.ToolCall{Id: "call_9", Name: "async_tool", Arguments: `{"term":"go"}`}

	err := wrapper.Send(context.Background(), parent, call, `{"term":"go"}`)

	assert.NoError(err)
	assert.Equal(1, delegate.sendCalls)
	assert.Equal(parent, delegate.lastParent)
	assert.Equal(call, delegate.lastCall)
	assert.Equal(wrapperReq{Term: "go"}, delegate.lastReq)
}

func TestStringWrapAsyncSendInvalidRequestDoesNotReachDelegate(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockAsyncTool{name: "async_tool"}
	wrapper := tool_string_wrapper.StringWrapAsync[wrapperReq](delegate)

	err := wrapper.Send(context.Background(), agent.AgentContext{}, agent.ToolCall{}, "broken")

	require.Error(t, err)
	assert.Equal(0, delegate.sendCalls)
}

func TestStringWrapAsyncSendDelegateError(t *testing.T) {
	assert := assert.New(t)

	sentinel := errors.New("send boom")
	delegate := &mockAsyncTool{
		name:   "async_tool",
		sendFn: func(ctx context.Context, p agent.AgentContext, c agent.ToolCall, r wrapperReq) error { return sentinel },
	}
	wrapper := tool_string_wrapper.StringWrapAsync[wrapperReq](delegate)

	err := wrapper.Send(context.Background(), agent.AgentContext{}, agent.ToolCall{}, `{"term":"go"}`)

	assert.ErrorIs(err, sentinel)
}

func TestStringWrapAsyncDelegatesMetadata(t *testing.T) {
	assert := assert.New(t)

	schema := newSchema(t)
	delegate := &mockAsyncTool{
		name:          "meta_async",
		description:   "async meta",
		auto:          false,
		requestSchema: schema,
	}
	wrapper := tool_string_wrapper.StringWrapAsync[wrapperReq](delegate)

	assert.Equal("meta_async", wrapper.Name())
	assert.Equal("async meta", wrapper.Description())
	assert.False(wrapper.Auto())
	assert.Equal(schema, wrapper.RequestSchema())
	assert.Equal("async-req:go", wrapper.DescribeRequest(`{"term":"go"}`))
}

func TestStringWrapAsyncDescribeRequestParseFailureMessage(t *testing.T) {
	assert := assert.New(t)

	delegate := &mockAsyncTool{name: "async_tool"}
	wrapper := tool_string_wrapper.StringWrapAsync[wrapperReq](delegate)

	assert.Equal("failed to parse request from string to type", wrapper.DescribeRequest("nope"))
}

func TestStringWrapAsyncDescribeRequestReceivesParsedRequest(t *testing.T) {
	assert := assert.New(t)

	var got wrapperReq
	delegate := &mockAsyncTool{
		name:          "async_tool",
		describeReqFn: func(req wrapperReq) string { got = req; return "described" },
	}
	wrapper := tool_string_wrapper.StringWrapAsync[wrapperReq](delegate)

	assert.Equal("described", wrapper.DescribeRequest(`{"term":"search"}`))
	assert.Equal(wrapperReq{Term: "search"}, got)
}
