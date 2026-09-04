package llm

import (
	"testing"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeepSeekSystemMessage(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("ctx-1", agent.MessageType_System, "be helpful", "rc-1", agent.ToolCall{})
	out := DeepSeekSystemMessage(msg)

	assert.Equal(deepseek.ChatMessageRoleSystem, out.Role)
	assert.Equal("be helpful", out.Content)
	assert.Equal("rc-1", out.ReasoningContent)
}

func TestDeepSeekAssistantMessage(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("ctx-1", agent.MessageType_Assistant, "ok", "rc-2", agent.ToolCall{})
	out := DeepSeekAssistantMessage(msg)

	assert.Equal(deepseek.ChatMessageRoleAssistant, out.Role)
	assert.Equal("ok", out.Content)
	assert.Equal("rc-2", out.ReasoningContent)
}

func TestDeepSeekUserMessage(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("ctx-1", agent.MessageType_User, "hi", "", agent.ToolCall{})
	out := DeepSeekUserMessage(msg)

	assert.Equal(deepseek.ChatMessageRoleUser, out.Role)
	assert.Equal("hi", out.Content)
	assert.Empty(out.ReasoningContent)
}

func TestDeepSeekToolRequestMessage(t *testing.T) {
	assert := assert.New(t)

	tool := agent.ToolCall{Id: "call_001", Name: "web_search", Arguments: `{"term":"deepseek"}`}
	msg := agent.NewMessage("ctx-1", agent.MessageType_ToolRequest, "handle this", "rc-3", tool)
	out := DeepSeekToolRequestMessage(msg)

	assert.Equal(deepseek.ChatMessageRoleAssistant, out.Role)
	assert.Equal("handle this with tool call id call_001", out.Content)
	assert.Equal("rc-3", out.ReasoningContent)

	require.Len(t, out.ToolCalls, 1)
	assert.Equal("call_001", out.ToolCalls[0].ID)
	assert.Equal("function", out.ToolCalls[0].Type)
	assert.Equal("web_search", out.ToolCalls[0].Function.Name)
	assert.Equal(`{"term":"deepseek"}`, out.ToolCalls[0].Function.Arguments)
}

func TestDeepSeekToolResultMessage(t *testing.T) {
	assert := assert.New(t)

	tool := agent.ToolCall{Id: "call_002", Name: "web_search", Arguments: "{}"}
	msg := agent.NewMessage("ctx-1", agent.MessageType_ToolResult, "the results", "rc-4", tool)
	out := DeepSeekToolResultMessage(msg)

	assert.Equal(deepseek.ChatMessageRoleTool, out.Role)
	assert.Equal("the results", out.Content)
	assert.Equal("rc-4", out.ReasoningContent)
	assert.Equal("call_002", out.ToolCallID)
}

func TestDeepSeekToolFromJsonSchemaSetsFields(t *testing.T) {
	assert := assert.New(t)

	schema, err := jsonschema.For[schemaStruct](&jsonschema.ForOptions{})
	require.NoError(t, err)

	out := DeepSeekToolFromJsonSchema("web_search", "search the web", schema)

	assert.Equal("function", out.Type)
	assert.Equal("web_search", out.Function.Name)
	assert.Equal("search the web", out.Function.Description)
	require.NotNil(t, out.Function.Parameters)
}

func TestDeepSeekToolFromJsonSchemaAcceptsNilSchema(t *testing.T) {
	assert := assert.New(t)

	out := DeepSeekToolFromJsonSchema("no_schema", "no params", nil)

	assert.Equal("function", out.Type)
	assert.Equal("no_schema", out.Function.Name)
	assert.Equal("no params", out.Function.Description)
	// nil schema converts to empty bytes; unmarshal yields the zero-valued struct
	assert.NotNil(out.Function.Parameters)
	assert.Empty(out.Function.Parameters.Properties)
}
