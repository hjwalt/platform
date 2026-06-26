package skill_tool

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
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
	skills map[string]agent_skill.Skill
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
	skill, found := t.lookupSkill(name)
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

// lookupSkill performs a case-insensitive exact match against the registry.
func (t *tool) lookupSkill(name string) (agent_skill.Skill, bool) {
	lower := strings.ToLower(name)
	for k, v := range t.skills {
		if strings.ToLower(k) == lower {
			return v, true
		}
	}
	return agent_skill.Skill{}, false
}

func (t *tool) Name() string {
	return Name
}

func (t *tool) Description() string {
	return "Load an on-demand skill (markdown playbook) into context. Only load a skill right before you need its guidance."
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
func Create(skills map[string]agent_skill.Skill) agent.SyncTool[Request, Response] {
	if skills == nil {
		skills = make(map[string]agent_skill.Skill)
	}
	return &tool{
		skills: skills,
	}
}

// AddToContainer registers the skill tool as a sync tool in the container.
func AddToContainer(container agent.ToolContainer, skills map[string]agent_skill.Skill) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(skills)))
}
