package agent_test

import (
	"testing"

	"github.com/hjwalt/platform/agent"
	"github.com/stretchr/testify/assert"
)

func TestMessageTypeConstants(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(agent.MessageType("SYSTEM"), agent.MessageType_System)
	assert.Equal(agent.MessageType("USER"), agent.MessageType_User)
	assert.Equal(agent.MessageType("TOOL_REQUEST"), agent.MessageType_ToolRequest)
	assert.Equal(agent.MessageType("TOOL_RESULT"), agent.MessageType_ToolResult)
	assert.Equal(agent.MessageType("TOOL_EXECUTE"), agent.MessageType_ToolExecute)
	assert.Equal(agent.MessageType("AGENT"), agent.MessageType_Agent)
	assert.Equal(agent.MessageType("ERROR"), agent.MessageType_Error)
	assert.Equal(agent.MessageType("FORK"), agent.MessageType_Start)
	assert.Equal(agent.MessageType("ASSISTANT"), agent.MessageType_Assistant)
}

func TestNewMessageSetsAllFields(t *testing.T) {
	assert := assert.New(t)

	tool := agent.ToolCall{Id: "tool-1", Name: "web_search", Arguments: `{"term":"go"}`}
	m := agent.NewMessage("ctx-1", agent.MessageType_User, "hello", "reasoning", tool)

	assert.NotEmpty(m.Id)
	assert.Equal("ctx-1", m.Context)
	assert.Equal(agent.MessageType_User, m.Type)
	assert.Equal("hello", m.Message)
	assert.Equal("reasoning", m.ReasoningContent)
	assert.Equal(tool, m.Tool)
	assert.Equal(agent.AgentContext{}, m.Agent)
}

func TestNewMessageGeneratesUniqueIds(t *testing.T) {
	assert := assert.New(t)
	a := agent.NewMessage("c", agent.MessageType_User, "x", "", agent.ToolCall{})
	b := agent.NewMessage("c", agent.MessageType_User, "x", "", agent.ToolCall{})
	assert.NotEmpty(a.Id)
	assert.NotEqual(a.Id, b.Id)
}

func TestStartSetsTypeAndContext(t *testing.T) {
	assert := assert.New(t)

	parent := agent.AgentContext{
		ParentContext: "parent-1",
		SystemMessage: "you are helpful",
		AllowedTools:  []string{"web_search", "web_fetch"},
	}
	tool := agent.ToolCall{Id: "tool-9", Name: "memory", Arguments: "{}"}

	m := agent.Start("ctx-2", "start message", tool, parent)

	assert.NotEmpty(m.Id)
	assert.Equal("ctx-2", m.Context)
	assert.Equal(agent.MessageType_Start, m.Type)
	assert.Equal("start message", m.Message)
	assert.Empty(m.ReasoningContent)
	assert.Equal(tool, m.Tool)
	// AgentContext must be propagated verbatim
	assert.Equal(parent, m.Agent)
}

func TestNewResult(t *testing.T) {
	assert := assert.New(t)

	messages := []agent.Message{
		agent.NewMessage("c", agent.MessageType_User, "one", "", agent.ToolCall{}),
		agent.NewMessage("c", agent.MessageType_Agent, "two", "", agent.ToolCall{}),
	}

	r := agent.NewResult(messages)
	assert.NotEmpty(r.Id)
	assert.Equal(messages, r.Messages)
	assert.Len(r.Messages, 2)
}

func TestSingleResult(t *testing.T) {
	assert := assert.New(t)

	m := agent.NewMessage("c", agent.MessageType_System, "only", "", agent.ToolCall{})

	r := agent.SingleResult(m)
	assert.NotEmpty(r.Id)
	assert.Len(r.Messages, 1)
	assert.Equal(m, r.Messages[0])
}

func TestEmptyResult(t *testing.T) {
	assert := assert.New(t)

	r := agent.EmptyResult()
	assert.NotEmpty(r.Id)
	assert.Empty(r.Messages)
	assert.NotNil(r.Messages, "EmptyResult initialises a non-nil empty slice")
}

func TestResultsGenerateUniqueIds(t *testing.T) {
	assert := assert.New(t)
	assert.NotEqual(agent.EmptyResult().Id, agent.EmptyResult().Id)
	assert.NotEqual(agent.NewResult(nil).Id, agent.NewResult(nil).Id)
}
