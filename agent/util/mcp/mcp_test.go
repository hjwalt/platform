package tool_mcp

import (
	"context"
	"errors"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mcpReq struct {
	Query string `json:"query"`
}

type mcpRes struct {
	Answer string `json:"answer"`
}

type mockMcpTool struct {
	name          string
	description   string
	requestSchema *jsonschema.Schema
	resultSchema  *jsonschema.Schema
	applyFn       func(ctx context.Context, req mcpReq) (mcpRes, error)
	applyCalls    int
	lastCtx       context.Context
	lastReq       mcpReq
}

func (m *mockMcpTool) Name() string                         { return m.name }
func (m *mockMcpTool) Description() string                  { return m.description }
func (m *mockMcpTool) Auto() bool                           { return false }
func (m *mockMcpTool) RequestSchema() *jsonschema.Schema    { return m.requestSchema }
func (m *mockMcpTool) ResultSchema() *jsonschema.Schema     { return m.resultSchema }
func (m *mockMcpTool) RequestFormat() format.Format[mcpReq] { return format.Json[mcpReq]() }
func (m *mockMcpTool) ResultFormat() format.Format[mcpRes]  { return format.Json[mcpRes]() }
func (m *mockMcpTool) DescribeRequest(req mcpReq) string    { return req.Query }
func (m *mockMcpTool) DescribeResult(res mcpRes) string     { return res.Answer }

var _ agent.SyncTool[mcpReq, mcpRes] = (*mockMcpTool)(nil)

func (m *mockMcpTool) Apply(ctx context.Context, req mcpReq) (mcpRes, error) {
	m.applyCalls++
	m.lastCtx = ctx
	m.lastReq = req
	if m.applyFn != nil {
		return m.applyFn(ctx, req)
	}
	return mcpRes{Answer: "default"}, nil
}

func TestMcpBehaviourInvokesToolWithParams(t *testing.T) {
	assert := assert.New(t)

	tool := &mockMcpTool{name: "echo"}
	handler := mcpBehaviour[mcpReq, mcpRes](tool)

	ctx := context.Background()
	request := &mcp.CallToolRequest{}
	params := mcpReq{Query: "what is go"}

	result, output, err := handler(ctx, request, params)

	assert.NoError(err)
	assert.Nil(result, "mcpBehaviour never populates a CallToolResult directly")
	assert.Equal(mcpRes{Answer: "default"}, output)
	assert.Equal(1, tool.applyCalls)
	assert.Equal(ctx, tool.lastCtx)
	assert.Equal(params, tool.lastReq)
}

func TestMcpBehaviourSatisfiesToolHandlerFor(t *testing.T) {
	var handler mcp.ToolHandlerFor[mcpReq, mcpRes] = mcpBehaviour[mcpReq, mcpRes](&mockMcpTool{name: "t"})
	require.NotNil(t, handler)
}

func TestMcpBehaviourPropagatesToolError(t *testing.T) {
	assert := assert.New(t)

	sentinel := errors.New("apply exploded")
	tool := &mockMcpTool{
		name: "echo",
		applyFn: func(ctx context.Context, req mcpReq) (mcpRes, error) {
			return mcpRes{}, sentinel
		},
	}
	handler := mcpBehaviour[mcpReq, mcpRes](tool)

	result, output, err := handler(context.Background(), &mcp.CallToolRequest{}, mcpReq{Query: "q"})

	assert.ErrorIs(err, sentinel)
	assert.Nil(result)
	assert.Empty(output)
	assert.Equal(1, tool.applyCalls)
}

func TestMcpBehaviourReceivesContextValues(t *testing.T) {
	assert := assert.New(t)

	type ctxKey struct{}
	tool := &mockMcpTool{name: "echo"}
	handler := mcpBehaviour[mcpReq, mcpRes](tool)

	ctx := context.WithValue(context.Background(), ctxKey{}, "v1")

	_, _, err := handler(ctx, &mcp.CallToolRequest{}, mcpReq{})

	assert.NoError(err)
	assert.Equal("v1", tool.lastCtx.Value(ctxKey{}))
}

func newSchema(t *testing.T) (*jsonschema.Schema, *jsonschema.Schema) {
	t.Helper()
	input, err := jsonschema.For[mcpReq](&jsonschema.ForOptions{})
	require.NoError(t, err)
	output, err := jsonschema.For[mcpRes](&jsonschema.ForOptions{})
	require.NoError(t, err)
	return input, output
}

func TestAddToMcpRegistersToolWithoutError(t *testing.T) {
	assert := assert.New(t)

	inputSchema, outputSchema := newSchema(t)
	server := mcp.NewServer(&mcp.Implementation{Name: "platform-test", Version: "1.0.0"}, nil)
	require.NotNil(t, server)

	// Registration is not directly inspectable (server state is unexported and
	// no in-memory transport is exposed), so this is a smoke test: AddToMcp
	// must build the tool metadata and register it without panicking.
	require.NotPanics(t, func() {
		AddToMcp(server, &mockMcpTool{
			name:          "echo",
			description:   "echoes back",
			requestSchema: inputSchema,
			resultSchema:  outputSchema,
		})
	})
	assert.True(true)
}

func TestAddToMcpWithoutSchemasPanics(t *testing.T) {
	assert := assert.New(t)

	server := mcp.NewServer(&mcp.Implementation{Name: "platform-test", Version: "1.0.0"}, nil)
	require.NotNil(t, server)

	// AddToMcp forwards the tool's (typed-nil when absent) RequestSchema into
	// the SDK's generic AddTool, which resolves the provided schema eagerly
	// and panics on nil. Tools exposed through AddToMcp must supply schemas.
	assert.Panics(func() {
		AddToMcp(server, &mockMcpTool{name: "echo", description: "echoes back"})
	})
}
