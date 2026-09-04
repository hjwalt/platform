package harness_container

import (
	"context"
	"errors"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mock tools
// ---------------------------------------------------------------------------

// mockSyncTool implements agent.SyncToolWrapper (SyncTool[string, string]).
type mockSyncTool struct {
	name           string
	description    string
	auto           bool
	schema         *jsonschema.Schema
	applyFn        func(ctx context.Context, request string) (string, error)
	describeReqFn  func(request string) string
	describeResFn  func(result string) string
	applyCalls     int
	lastApplyInput string
}

func (m *mockSyncTool) Name() string                         { return m.name }
func (m *mockSyncTool) Description() string                  { return m.description }
func (m *mockSyncTool) RequestFormat() format.Format[string] { return format.String() }
func (m *mockSyncTool) RequestSchema() *jsonschema.Schema    { return m.schema }
func (m *mockSyncTool) ResultFormat() format.Format[string]  { return format.String() }
func (m *mockSyncTool) ResultSchema() *jsonschema.Schema     { return nil }
func (m *mockSyncTool) Auto() bool                           { return m.auto }

func (m *mockSyncTool) DescribeRequest(request string) string {
	if m.describeReqFn != nil {
		return m.describeReqFn(request)
	}
	return "request: " + request
}

func (m *mockSyncTool) DescribeResult(result string) string {
	if m.describeResFn != nil {
		return m.describeResFn(result)
	}
	return "result: " + result
}

func (m *mockSyncTool) Apply(ctx context.Context, request string) (string, error) {
	m.applyCalls++
	m.lastApplyInput = request
	if m.applyFn != nil {
		return m.applyFn(ctx, request)
	}
	return "applied:" + request, nil
}

// mockAsyncTool implements agent.AsyncToolWrapper (AsyncTool[string]).
type mockAsyncTool struct {
	name        string
	description string
	auto        bool
	schema      *jsonschema.Schema
	sendFn      func(ctx context.Context, parent agent.AgentContext, call agent.ToolCall, request string) error
	sendCalls   int
	lastParent  agent.AgentContext
	lastCall    agent.ToolCall
	lastRequest string
}

func (m *mockAsyncTool) Name() string                         { return m.name }
func (m *mockAsyncTool) Description() string                  { return m.description }
func (m *mockAsyncTool) RequestFormat() format.Format[string] { return format.String() }
func (m *mockAsyncTool) RequestSchema() *jsonschema.Schema    { return m.schema }
func (m *mockAsyncTool) Auto() bool                           { return m.auto }

func (m *mockAsyncTool) DescribeRequest(request string) string {
	return "async request: " + request
}

func (m *mockAsyncTool) Send(ctx context.Context, parent agent.AgentContext, call agent.ToolCall, request string) error {
	m.sendCalls++
	m.lastParent = parent
	m.lastCall = call
	m.lastRequest = request
	if m.sendFn != nil {
		return m.sendFn(ctx, parent, call, request)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func newTestContainer() *toolContainer {
	return NewToolContainer().(*toolContainer)
}

func TestNewToolContainerIsNonNilAndEmpty(t *testing.T) {
	assert := assert.New(t)

	container := NewToolContainer()
	require.NotNil(t, container)
	assert.Empty(container.OpenAiParamsFiltered(nil))
	assert.Empty(container.DeepSeekParams(nil))
	assert.False(container.Exists(agent.ToolCall{Name: "nope"}))
}

func TestAddSyncRegistersTool(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	tool := &mockSyncTool{name: "sync_tool", description: "does sync", auto: true}
	container.AddSync(tool)

	assert.True(container.Exists(agent.ToolCall{Name: "sync_tool"}))
	assert.False(container.Exists(agent.ToolCall{Name: "other"}))
	assert.True(container.Auto(agent.ToolCall{Name: "sync_tool"}))
	assert.False(container.Auto(agent.ToolCall{Name: "other"}))
}

func TestAddAsyncRegistersTool(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	tool := &mockAsyncTool{name: "async_tool", description: "does async", auto: false}
	container.AddAsync(tool)

	assert.True(container.Exists(agent.ToolCall{Name: "async_tool"}))
	assert.False(container.Auto(agent.ToolCall{Name: "async_tool"}))
}

func TestExecuteSyncToolSuccess(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	tool := &mockSyncTool{
		name:          "sync_tool",
		applyFn:       func(ctx context.Context, request string) (string, error) { return "tool-out", nil },
		describeResFn: func(result string) string { return "described:" + result },
	}
	container.AddSync(tool)

	in := agent.NewMessage("ctx-1", agent.MessageType_ToolExecute, "go", "", agent.ToolCall{})
	call := agent.ToolCall{Name: "sync_tool", Arguments: `{"x":1}`}

	result, err := container.Execute(context.Background(), in, call)

	assert.NoError(err)
	assert.True(result.IsPresent())
	assert.Equal("described:tool-out", result.Get())
	assert.Equal(1, tool.applyCalls)
	assert.Equal(`{"x":1}`, tool.lastApplyInput)
}

func TestExecuteSyncToolError(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	sentinel := errors.New("apply exploded")
	tool := &mockSyncTool{name: "sync_tool", applyFn: func(ctx context.Context, request string) (string, error) { return "", sentinel }}
	container.AddSync(tool)

	_, err := container.Execute(context.Background(), agent.Message{}, agent.ToolCall{Name: "sync_tool"})

	assert.ErrorIs(err, sentinel)
}

func TestExecuteAsyncToolSuccessIsEmptyOptional(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	parent := agent.AgentContext{ParentContext: "ctx-1"}
	in := agent.Message{Context: "ctx-1", Type: agent.MessageType_ToolExecute, Message: "go"}
	tool := &mockAsyncTool{name: "async_tool"}
	container.AddAsync(tool)
	call := agent.ToolCall{Name: "async_tool", Id: "c1", Arguments: "abc"}

	result, err := container.Execute(context.Background(), in, call)

	assert.NoError(err)
	assert.False(result.IsPresent())
	assert.Equal(1, tool.sendCalls)
	// async tools receive a parent context derived from the incoming message
	assert.Equal(parent, tool.lastParent)
	assert.Equal(call, tool.lastCall)
	assert.Equal("abc", tool.lastRequest)
}

func TestExecuteAsyncToolError(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	sentinel := errors.New("send failed")
	tool := &mockAsyncTool{name: "async_tool", sendFn: func(ctx context.Context, p agent.AgentContext, c agent.ToolCall, r string) error { return sentinel }}
	container.AddAsync(tool)

	_, err := container.Execute(context.Background(), agent.Message{}, agent.ToolCall{Name: "async_tool"})

	assert.ErrorIs(err, sentinel)
}

func TestExecuteUnknownTool(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	_, err := container.Execute(context.Background(), agent.Message{}, agent.ToolCall{Name: "ghost"})

	require.Error(t, err)
	assert.Equal("tool ghost does not exist", err.Error())
}

func TestDescribeRequestSyncAndAsync(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	syncTool := &mockSyncTool{name: "sync_tool", describeReqFn: func(r string) string { return "sync-desc:" + r }}
	asyncTool := &mockAsyncTool{name: "async_tool"}
	container.AddSync(syncTool)
	container.AddAsync(asyncTool)

	desc, err := container.DescribeRequest(agent.ToolCall{Name: "sync_tool", Arguments: "a1"})
	assert.NoError(err)
	assert.Equal("sync-desc:a1", desc)

	desc, err = container.DescribeRequest(agent.ToolCall{Name: "async_tool", Arguments: "a2"})
	assert.NoError(err)
	assert.Equal("async request: a2", desc)

	_, err = container.DescribeRequest(agent.ToolCall{Name: "missing"})
	require.Error(t, err)
	assert.Equal("tool missing does not exist", err.Error())
}

func TestOpenAiParamsFilteredEmptyAllowedListsAllTools(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	schema, schemaErr := jsonschema.For[schemaTestType](&jsonschema.ForOptions{})
	require.NoError(t, schemaErr)

	container.AddSync(&mockSyncTool{name: "a", description: "desc-a", schema: schema})
	container.AddAsync(&mockAsyncTool{name: "b", description: "desc-b", schema: schema})

	params := container.OpenAiParamsFiltered(nil)

	require.Len(t, params, 2)
	assert.NotNil(t, params[0].OfFunction)
	assert.NotNil(t, params[1].OfFunction)
	assert.NotEmpty(t, params[0].OfFunction.Function.Parameters)
}

func TestOpenAiParamsFilteredFiltersByAllowedList(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	container.AddSync(&mockSyncTool{name: "a", description: "desc-a"})
	container.AddAsync(&mockAsyncTool{name: "b", description: "desc-b"})

	params := container.OpenAiParamsFiltered([]string{"b"})

	require.Len(t, params, 1)
	require.NotNil(t, params[0].OfFunction)
	assert.Equal("b", params[0].OfFunction.Function.Name)
	assert.Equal("desc-b", params[0].OfFunction.Function.Description.Value)
}

func TestOpenAiParamsFilteredSkipsUnknownNames(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	container.AddSync(&mockSyncTool{name: "a", description: "desc-a"})

	params := container.OpenAiParamsFiltered([]string{"a", "ghost"})

	require.Len(t, params, 1)
	require.NotNil(t, params[0].OfFunction)
	assert.Equal("a", params[0].OfFunction.Function.Name)
}

func TestOpenAiParamsFilteredWithNilSchema(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	container.AddSync(&mockSyncTool{name: "a", description: "desc-a"}) // schema nil

	params := container.OpenAiParamsFiltered(nil)

	require.Len(t, params, 1)
	require.NotNil(t, params[0].OfFunction)
	assert.Equal("a", params[0].OfFunction.Function.Name)
	// a nil *jsonschema.Schema converts through empty bytes to a nil map
	assert.Nil(params[0].OfFunction.Function.Parameters)
}

func TestDeepSeekParamsEmptyAllowedListsAllTools(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	container.AddSync(&mockSyncTool{name: "a", description: "desc-a"})
	container.AddAsync(&mockAsyncTool{name: "b", description: "desc-b"})

	params := container.DeepSeekParams(nil)

	require.Len(t, params, 2)
	assert.Equal("function", params[0].Type)
}

func TestDeepSeekParamsFiltersAndSkipsUnknown(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	container.AddSync(&mockSyncTool{name: "a", description: "desc-a"})
	container.AddAsync(&mockAsyncTool{name: "b", description: "desc-b"})

	params := container.DeepSeekParams([]string{"b", "ghost", "a"})

	require.Len(t, params, 2)
	names := []string{params[0].Function.Name, params[1].Function.Name}
	assert.Contains(names, "a")
	assert.Contains(names, "b")
	byName := map[string]string{params[0].Function.Name: params[0].Function.Description,
		params[1].Function.Name: params[1].Function.Description}
	assert.Equal("desc-a", byName["a"])
	assert.Equal("desc-b", byName["b"])
}

func TestAutoSyncAndAsyncRespected(t *testing.T) {
	assert := assert.New(t)
	container := newTestContainer()

	container.AddSync(&mockSyncTool{name: "auto_sync", auto: true})
	container.AddSync(&mockSyncTool{name: "manual_sync", auto: false})
	container.AddAsync(&mockAsyncTool{name: "auto_async", auto: true})

	assert.True(container.Auto(agent.ToolCall{Name: "auto_sync"}))
	assert.False(container.Auto(agent.ToolCall{Name: "manual_sync"}))
	assert.True(container.Auto(agent.ToolCall{Name: "auto_async"}))
	assert.False(container.Auto(agent.ToolCall{Name: "ghost"}))
}

type schemaTestType struct {
	Query string `json:"query"`
}
