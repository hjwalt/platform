package configuration

import (
	"log/slog"

	agent_skill "github.com/hjwalt/platform/agent/skill"
	harness_container "github.com/hjwalt/platform/agent/util/container"
)

func RegisterSkills(holder Context, conf Configuration) {
	container := harness_container.NewSkillContainer()
	holder.SetSkillContainer(container)

	// TODO: load directly from directory(ies)
	skillPaths := []string{
		"./docs/skills/self-improving",
	}
	for _, dir := range skillPaths {
		props, err := agent_skill.ReadProperties(dir)
		if err != nil {
			slog.Warn("failed to load skill for registry", "dir", dir, "error", err)
			continue
		}
		container.Add(props)
		slog.Info("loaded skill into registry", "name", props.Name, "dir", dir)
	}
}
