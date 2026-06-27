package skill_tool

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/format"
)

const Name = "skill"

// Request is the input for the skill tool.
type Request struct {
	Name string `json:"name" jsonschema:"name of the skill to load"`
}

// Response is the output of the skill tool.
type Response struct {
	Name         string   `json:"name" jsonschema:"name of the skill"`
	Description  string   `json:"description" jsonschema:"description of the skill"`
	Body         string   `json:"body" jsonschema:"markdown playbook content loaded into context"`
	AllowedTools []string `json:"allowed_tools" jsonschema:"tools allowed by the skill"`
	Found        bool     `json:"found" jsonschema:"whether the skill was found in the registry"`
}

type tool struct {
	skills agent.SkillContainer
}

func (t *tool) Apply(ctx context.Context, params Request) (Response, error) {
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return Response{
			Name:  params.Name,
			Found: false,
		}, nil
	}

	// Case-insensitive lookup
	skill, found := t.skills.Get(name)
	if !found {
		return Response{
			Name:  params.Name,
			Found: false,
		}, nil
	}

	return Response{
		Name:         skill.Name,
		Description:  skill.Description,
		Body:         skill.Body,
		AllowedTools: skill.AllowedTools,
		Found:        true,
	}, nil
}

func (t *tool) Name() string {
	return Name
}

func (t *tool) Description() string {
	return "Load an on-demand skill (markdown playbook) into context. Only load a skill right before you need its guidance."
}

func (t *tool) RequestFormat() format.Format[Request] {
	return RequestFormat
}

func (t *tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *tool) DescribeRequest(request Request) string {
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString("loading skill `")
	outputBuilder.WriteString(request.Name)
	outputBuilder.WriteString("`")
	return outputBuilder.String()
}

func (t *tool) ResultFormat() format.Format[Response] {
	return format.Json[Response]()
}

func (t *tool) ResultSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Response](opts)
	return toolSchema
}

func (t *tool) DescribeResult(response Response) string {
	if !response.Found {
		return fmt.Sprintf("skill `%s` not found", response.Name)
	}

	outputBuilder := strings.Builder{}
	outputBuilder.WriteString("loaded skill `")
	outputBuilder.WriteString(response.Name)
	outputBuilder.WriteString("`")
	outputBuilder.WriteString("\n\n")
	outputBuilder.WriteString(response.Description)
	outputBuilder.WriteString("\n\n")
	outputBuilder.WriteString("---\n")
	outputBuilder.WriteString(response.Body)
	return outputBuilder.String()
}

func (t *tool) Auto() bool {
	return true
}

// Create constructs a new skill SyncTool from a skill registry.
// The registry maps skill names to their parsed Skill definitions.
func Create(skills agent.SkillContainer) agent.SyncTool[Request, Response] {
	if skills == nil {
		skills = harness_container.NewSkillContainer()
	}
	return &tool{
		skills: skills,
	}
}

// AddToContainer registers the skill tool as a sync tool in the container.
func AddToContainer(container agent.ToolContainer, skills agent.SkillContainer) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(skills)))
}

func Parse(args string) (Request, error) {
	return RequestFormat.Unmarshal([]byte(args))
}

var RequestFormat = format.Json[Request]()
