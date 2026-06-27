package subagent_tool

import (
	"context"
	"strings"
	"testing"

	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/hjwalt/platform/flow"
	"github.com/stretchr/testify/assert"
)

// stubProducer is a simple stub implementation of flow.Producer[agent.Message].
type stubProducer struct {
	messages []agent.Message
}

func (s *stubProducer) Produce(_ context.Context, msgs []agent.Message) error {
	s.messages = append(s.messages, msgs...)
	return nil
}

func (s *stubProducer) ProduceMessage(_ context.Context, _ []flow.Message[agent.Message]) error {
	return nil
}

func (s *stubProducer) Start() error {
	return nil
}

func (s *stubProducer) Stop() {}

func TestSubagentName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "A test agent",
		Skill:       "## Playbook\nDo stuff",
	}, &stubProducer{})

	assert.Equal("test_agent", tool.Name())
}

func TestSubagentNameWithHyphens(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "my-custom-agent",
		Description: "An agent with hyphens",
		Skill:       "## Playbook",
	}, &stubProducer{})

	assert.Equal("my_custom_agent", tool.Name())
}

func TestSubagentDescription(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "A test agent description",
		Skill:       "## Playbook",
	}, &stubProducer{})

	assert.Equal("A test agent description", tool.Description())
}

func TestSubagentAuto(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "desc",
		Skill:       "## Playbook",
	}, &stubProducer{})

	assert.True(tool.Auto())
}

func TestSubagentRequestSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "desc",
		Skill:       "## Playbook",
	}, &stubProducer{})

	schema := tool.RequestSchema()
	assert.NotNil(schema)
}

func TestSubagentRequestFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "desc",
		Skill:       "## Playbook",
	}, &stubProducer{})

	assert.NotNil(tool.RequestFormat())
}

func TestSubagentDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "code-reviewer",
		Description: "Reviews code",
		Skill:       "## Code Review Playbook\n\n1. Check correctness",
	}, &stubProducer{})

	desc := tool.DescribeRequest(Request{Prompt: "review main.go"})

	assert.Contains(desc, "code_reviewer")
	assert.Contains(desc, "review main.go")
}

func TestSubagentDescribeRequestContainsAgentName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "deploy-agent",
		Description: "Deploys services",
		Skill:       "## Deploy Playbook",
	}, &stubProducer{})

	desc := tool.DescribeRequest(Request{Prompt: "deploy to production"})

	assert.Contains(desc, "deploy_agent")
	assert.Contains(desc, "deploy to production")
}

func TestSubagentAllowedTools(t *testing.T) {
	assert := assert.New(t)

	allowedTools := []string{"linux_shell", "web_fetch", "web_search"}
	tool := Create(Configuration{
		Name:         "test-agent",
		Description:  "desc",
		Skill:        "## Playbook",
		AllowedTools: allowedTools,
	}, &stubProducer{})

	// Verify the tool is created successfully
	assert.NotNil(tool)
	assert.Equal("test_agent", tool.Name())
}

func TestSubagentCreateReturnsAsyncTool(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "desc",
		Skill:       "## Playbook",
	}, &stubProducer{})

	var _ agent.AsyncTool[Request] = tool
	assert.NotNil(tool)
}

func TestSubagentFromSubagent(t *testing.T) {
	assert := assert.New(t)

	skill := agent_skill.Skill{
		Name:         "code-review",
		Description:  "Review code for bugs",
		Body:         "## Playbook\n\n1. Check correctness",
		AllowedTools: []string{"linux_shell", "web_fetch"},
	}

	tool := FromSubagent(skill, &stubProducer{})

	assert.Equal("code-review", tool.Name())
	assert.Equal("Review code for bugs", tool.Description())
	assert.NotNil(tool)
}

func TestSubagentFromSubagentNamePreserved(t *testing.T) {
	assert := assert.New(t)

	skill := agent_skill.Skill{
		Name:         "my-skill-name",
		Description:  "A skill with hyphens in name",
		Body:         "## Playbook",
		AllowedTools: nil,
	}

	// FromSubagent preserves the original name (no replace of hyphens)
	tool := FromSubagent(skill, &stubProducer{})

	assert.Equal("my-skill-name", tool.Name())
	assert.Equal("A skill with hyphens in name", tool.Description())
}

func TestSubagentSend(t *testing.T) {
	assert := assert.New(t)

	producer := &stubProducer{}
	tool := Create(Configuration{
		Name:        "test-agent",
		Description: "desc",
		Skill:       "## Playbook\n\nDo X then Y",
	}, producer)

	parentCtx := agent.AgentContext{
		ParentContext: "parent context id",
		SystemMessage: "test system message",
		AllowedTools:  []string{"tool_a"},
	}

	err := tool.Send(nil, parentCtx, agent.ToolCall{Id: "call-1"}, Request{Prompt: "do something"})

	assert.NoError(err)
	assert.Len(producer.messages, 1)
	assert.Equal(agent.MessageType_Start, producer.messages[0].Type)
	assert.Equal("call-1", producer.messages[0].Tool.Id)
	assert.Equal("do something", producer.messages[0].Message)
}

func TestSubagentSendPreservesContext(t *testing.T) {
	assert := assert.New(t)

	producer := &stubProducer{}
	tool := Create(Configuration{
		Name:         "test-agent",
		Description:  "desc",
		Skill:        "## System prompt body",
		AllowedTools: []string{"linux_shell", "web_search"},
	}, producer)

	parentCtx := agent.AgentContext{
		ParentContext: "parent-id",
		SystemMessage: "parent system",
		AllowedTools:  []string{"parent_tool"},
	}

	err := tool.Send(nil, parentCtx, agent.ToolCall{Id: "call-2"}, Request{Prompt: "run task"})

	assert.NoError(err)
	assert.Len(producer.messages, 1)

	msg := producer.messages[0]
	assert.Equal(agent.MessageType_Start, msg.Type)
	assert.Equal("call-2", msg.Tool.Id)

	// Verify the context carries the system prompt from the skill body
	assert.Equal("## System prompt body", msg.Agent.SystemMessage)
	// Verify allowed tools are passed through
	assert.Equal([]string{"linux_shell", "web_search"}, msg.Agent.AllowedTools)
}

func TestAddToContainerRegistersSubagentTool(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	AddToContainer(container, Configuration{
		Name:        "my-agent",
		Description: "My custom agent",
		Skill:       "## Playbook",
	}, &stubProducer{})

	assert.True(container.Exists(agent.ToolCall{Name: "my_agent"}))
}

func TestAddSubagentToContainer(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	skill := agent_skill.Skill{
		Name:         "deploy",
		Description:  "Deploy services",
		Body:         "## Deploy Playbook",
		AllowedTools: []string{"linux_shell"},
	}

	AddSubagentToContainer(container, skill, &stubProducer{})

	assert.True(container.Exists(agent.ToolCall{Name: "deploy"}))
}

func TestSubagentDescribeRequestMultiline(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		Name:        "analyzer",
		Description: "Analyzes things",
		Skill:       "## Analyzer Playbook",
	}, &stubProducer{})

	desc := tool.DescribeRequest(Request{Prompt: "line1\nline2\nline3"})

	assert.Contains(desc, "analyzer")
	assert.True(strings.Contains(desc, "line1"), "should contain multiline prompt")
}
