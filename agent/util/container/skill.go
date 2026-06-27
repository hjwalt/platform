package harness_container

import (
	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	"github.com/hjwalt/platform/type/optional"
)

func NewSkillContainer() agent.SkillContainer {
	return &skillContainer{
		skills: make(map[string]agent_skill.Skill),
	}
}

type skillContainer struct {
	skills map[string]agent_skill.Skill
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
	r.skills[in.Name] = agent_skill.Skill{
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
	s, ok := r.skills[name]
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
