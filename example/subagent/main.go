package main

import (
	"log/slog"

	agent_skill "github.com/hjwalt/platform/agent/skill"
)

func main() {
	properties, err := agent_skill.ReadProperties("./docs/agents/researcher-agent")
	slog.Info("parsed", "properties", properties, "error", err)
}
