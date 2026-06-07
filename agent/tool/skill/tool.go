package skill_tool

import (
	"context"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/format"
)

type Configuration struct {
	Name         string
	Description  string
	Skill        string
	AllowedTools []string
}

type Request struct {
	Prompt string `json:"prompt" jsonschema:"prompt for this agent"`
}

type tool struct {
	AgentName        string
	AgentDescription string
	SystemPrompt     string
	AllowedTools     []string
	Producer         flow.Producer[agent.Message]
}

func (t *tool) Send(ctx context.Context, parent agent.AgentContext, toolCall agent.ToolCall, params Request) error {
	t.Producer.Produce(ctx, []agent.Message{agent.Start(
		toolCall.Id,
		params.Prompt,
		toolCall,
		agent.AgentContext{
			ParentContext: parent.ParentContext,
			SystemMessage: t.SystemPrompt,
			AllowedTools:  t.AllowedTools,
		},
	)})
	return nil
}

func (t *tool) Name() string {
	return t.AgentName
}

func (t *tool) Description() string {
	return t.AgentDescription
}

func (t *tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *tool) DescribeRequest(request Request) string {
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString("Running agent: ")
	outputBuilder.WriteString("\n")
	outputBuilder.WriteString(t.AgentName)
	outputBuilder.WriteString("\n\n")
	outputBuilder.WriteString("With prompt: ")
	outputBuilder.WriteString("\n")
	outputBuilder.WriteString(request.Prompt)
	return outputBuilder.String()
}

func (t *tool) Auto() bool {
	return true
}

func Create(config Configuration, producer flow.Producer[agent.Message]) agent.AsyncTool[Request] {
	return &tool{
		AgentName:        config.Name,
		AgentDescription: config.Description,
		SystemPrompt:     config.Skill,
		AllowedTools:     config.AllowedTools,
		Producer:         producer,
	}
}

func FromSkill(config agent_skill.Skill, producer flow.Producer[agent.Message]) agent.AsyncTool[Request] {
	return &tool{
		AgentName:        config.Name,
		AgentDescription: config.Description,
		SystemPrompt:     config.Body,
		AllowedTools:     config.AllowedTools,
		Producer:         producer,
	}
}

func AddToContainer(container agent.ToolContainer, config Configuration, producer flow.Producer[agent.Message]) {
	container.AddAsync(tool_string_wrapper.StringWrapAsync(Create(config, producer)))
}

func AddSkillToContainer(container agent.ToolContainer, config agent_skill.Skill, producer flow.Producer[agent.Message]) {
	container.AddAsync(tool_string_wrapper.StringWrapAsync(FromSkill(config, producer)))
}
