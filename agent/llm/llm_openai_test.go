package llm

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// schemaStruct is a small struct used to derive real JSON schemas.
type schemaStruct struct {
	Term string `json:"term" jsonschema:"the search term"`
}

func TestOpenAiSystemMessage(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("ctx-1", agent.MessageType_System, "be helpful", "rc", agent.ToolCall{})
	union := OpenAiSystemMessage(msg)

	require.NotNil(t, union.OfSystem)
	assert.Equal("be helpful", union.OfSystem.Content.OfString.Value)
	assert.Nil(union.OfUser)
	assert.Nil(union.OfAssistant)
	assert.Nil(union.OfTool)
}

func TestOpenAiAssistantMessage(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("ctx-1", agent.MessageType_Assistant, "sure thing", "", agent.ToolCall{})
	union := OpenAiAssistantMessage(msg)

	require.NotNil(t, union.OfAssistant)
	assert.Equal("sure thing", union.OfAssistant.Content.OfString.Value)
	assert.Nil(union.OfSystem)
}

func TestOpenAiUserMessage(t *testing.T) {
	assert := assert.New(t)

	msg := agent.NewMessage("ctx-1", agent.MessageType_User, "hello", "", agent.ToolCall{})
	union := OpenAiUserMessage(msg)

	require.NotNil(t, union.OfUser)
	assert.Equal("hello", union.OfUser.Content.OfString.Value)
	assert.Nil(union.OfSystem)
	assert.Nil(union.OfAssistant)
}

func TestOpenAiToolRequestMessage(t *testing.T) {
	assert := assert.New(t)

	tool := agent.ToolCall{Id: "call_123", Name: "web_search", Arguments: `{"term":"go"}`}
	msg := agent.NewMessage("ctx-1", agent.MessageType_ToolRequest, "search the web", "", tool)
	union := OpenAiToolRequestMessage(msg)

	require.NotNil(t, union.OfAssistant)
	require.Len(t, union.OfAssistant.ToolCalls, 1)
	call := union.OfAssistant.ToolCalls[0]
	require.NotNil(t, call.OfFunction)
	assert.Equal("call_123", call.OfFunction.ID)
	assert.Equal("web_search", call.OfFunction.Function.Name)
	assert.Equal(`{"term":"go"}`, call.OfFunction.Function.Arguments)
}

func TestOpenAiToolResultMessage(t *testing.T) {
	assert := assert.New(t)

	tool := agent.ToolCall{Id: "call_456", Name: "web_search", Arguments: "{}"}
	msg := agent.NewMessage("ctx-1", agent.MessageType_ToolResult, "results here", "", tool)
	union := OpenAiToolResultMessage(msg)

	require.NotNil(t, union.OfTool)
	assert.Equal("results here", union.OfTool.Content.OfString.Value)
	assert.Equal("call_456", union.OfTool.ToolCallID)
}

func TestOpenAiFromJsonSchemaSetsMetadata(t *testing.T) {
	assert := assert.New(t)

	schema, err := jsonschema.For[schemaStruct](&jsonschema.ForOptions{})
	require.NoError(t, err)

	union := OpenAiFromJsonSchema("web_search", "search the web", schema)

	require.NotNil(t, union.OfFunction)
	assert.Equal("web_search", union.OfFunction.Function.Name)
	assert.Equal("search the web", union.OfFunction.Function.Description.Value)
	assert.NotEmpty(union.OfFunction.Function.Parameters,
		"parameters should contain the JSON-serialised schema")
}

func TestOpenAiFromJsonSchemaAcceptsNilSchema(t *testing.T) {
	assert := assert.New(t)

	union := OpenAiFromJsonSchema("empty_tool", "does nothing", nil)

	require.NotNil(t, union.OfFunction)
	assert.Equal("empty_tool", union.OfFunction.Function.Name)
	assert.Equal("does nothing", union.OfFunction.Function.Description.Value)
	// nil schema converts to empty bytes which unmarshal to an empty (nil) map
	assert.Nil(union.OfFunction.Function.Parameters)
}

func TestOpenAiToolSchemaBuildsFunctionTool(t *testing.T) {
	assert := assert.New(t)

	union := OpenAiToolSchema[schemaStruct]("web_fetch", "fetch a url")

	require.NotNil(t, union.OfFunction)
	assert.Equal("web_fetch", union.OfFunction.Function.Name)
	assert.Equal("fetch a url", union.OfFunction.Function.Description.Value)
	assert.NotEmpty(union.OfFunction.Function.Parameters)
}

func TestOpenAiSchemaResultIsUnionOfFunction(t *testing.T) {
	assert := assert.New(t)

	var _ openai.ChatCompletionToolUnionParam = OpenAiToolSchema[schemaStruct]("x", "y")
	schema, err := jsonschema.For[schemaStruct](&jsonschema.ForOptions{})
	require.NoError(t, err)

	union := OpenAiFromJsonSchema("x", "y", schema)
	require.NotNil(t, union.OfFunction)
	assert.Nil(union.OfCustom)
}
