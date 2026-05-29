package fork_tool

import (
	"context"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/format"
)

type Configuration struct {
	AgentName    string
	SystemPrompt string
	AllowedTools []string
}

type Request struct {
	Prompt string `json:"prompt" jsonschema:"prompt for this agent"`
}

type Tool struct {
	AgentName    string
	SystemPrompt string
	AllowedTools []string
	Producer     flow.Producer[agent.Message]
}

func (t *Tool) Send(ctx context.Context, parent agent.AgentContext, toolCall agent.ToolCall, params Request) error {
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

func (t *Tool) Name() string {
	return t.AgentName
}

func (t *Tool) Description() string {
	return "Achieve specific goals based on the prompt for this agent. \n\n The agent's capability is defined as follows: \n\n" + t.SystemPrompt
}

func (t *Tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *Tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *Tool) DescribeRequest(request Request) string {
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

func (t *Tool) Auto() bool {
	return true
}

func Create(config Configuration, producer flow.Producer[agent.Message]) agent.AsyncTool[Request] {
	return &Tool{
		AgentName:    config.AgentName,
		SystemPrompt: config.SystemPrompt,
		AllowedTools: config.AllowedTools,
		Producer:     producer,
	}
}

func AddToContainer(container agent.ToolContainer, config Configuration, producer flow.Producer[agent.Message]) {
	container.AddAsync(tool_string_wrapper.StringWrapAsync(Create(config, producer)))
}
