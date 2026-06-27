package harness_container

import (
	"strings"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/type/optional"
)

func NewSkillContainer() agent.SkillContainer {
	return &skillContainer{
		skills: make(map[string]agent.Instruction),
	}
}

type skillContainer struct {
	skills map[string]agent.Instruction
}

func (r *skillContainer) Add(in agent.Instruction) {
	license := in.License
	if license == nil {
		license = optional.Empty[string]()
	}
	compatibility := in.Compatibility
	if compatibility == nil {
		compatibility = optional.Empty[string]()
	}
	r.skills[strings.ToLower(in.Name)] = agent.Instruction{
		Name:          in.Name,
		Description:   in.Description,
		License:       license,
		Compatibility: compatibility,
		AllowedTools:  in.AllowedTools,
		Metadata:      in.Metadata,
		Body:          in.Body,
	}
}

func (r *skillContainer) Get(name string) (agent.Instruction, bool) {
	s, ok := r.skills[strings.ToLower(name)]
	if !ok {
		return agent.Instruction{}, false
	}
	return agent.Instruction{
		Name:          s.Name,
		Description:   s.Description,
		License:       s.License,
		Compatibility: s.Compatibility,
		AllowedTools:  s.AllowedTools,
		Metadata:      s.Metadata,
		Body:          s.Body,
	}, true
}

func (r *skillContainer) Assistant(ctx string) agent.Message {
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString("available skills:")
	outputBuilder.WriteString("\n")
	outputBuilder.WriteString("\n")
	for _, instruction := range r.skills {

		outputBuilder.WriteString("name:\n")
		outputBuilder.WriteString(instruction.Name)
		outputBuilder.WriteString("\n")
		outputBuilder.WriteString("description:\n")
		outputBuilder.WriteString(instruction.Description)
		outputBuilder.WriteString("\n\n")
	}
	return agent.NewMessage(
		ctx,
		agent.MessageType_Assistant,
		outputBuilder.String(),
		"",
		agent.ToolCall{},
	)
}
