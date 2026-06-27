package configuration

import (
	"log/slog"

	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	finance_stock_price_tool "github.com/hjwalt/platform/agent/tool/finance_stock_price"
	memory_tool "github.com/hjwalt/platform/agent/tool/memory"
	skill_tool "github.com/hjwalt/platform/agent/tool/skill"
	subagent_tool "github.com/hjwalt/platform/agent/tool/subagent"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	web_search_tool "github.com/hjwalt/platform/agent/tool/web_search"
	harness_container "github.com/hjwalt/platform/agent/util/container"
)

func RegisterTools(holder Context, conf Configuration) {
	container := harness_container.NewToolContainer()

	// register tools
	web_search_tool.AddToContainer(container, conf.Tool.WebSearch)
	// linux_shell_tool.AddToContainer(container, conf.Tool.Shell)
	web_fetch_tool.AddToContainer(container, conf.Tool.WebFetch, holder.GetParserModel())
	finance_fx_price_tool.AddToContainer(container, conf.Tool.FxPrice)
	finance_stock_price_tool.AddToContainer(container, conf.Tool.StockPrice)

	for _, memoryConfig := range conf.Tool.Memory {
		if validateErr := memory_tool.Validate(memoryConfig, holder.GetMemoryStore()); validateErr != nil {
			slog.Error("failed to register memory tool", "key", memoryConfig.Key, "error", validateErr)
			continue
		}
		memory_tool.AddToContainer(container, memoryConfig, holder.GetMemoryStore())
	}

	// register skill tool for on-demand context injection
	registerSkills(
		holder,
		conf,
		container,
		"./docs/skills/self-improving",
	)

	// register subagents (async delegation)
	registerSubAgents(
		holder,
		conf,
		container,
		"./docs/agents/researcher-agent",
	)

	holder.SetToolContainer(container)
}

// registerSkills reads skill properties from the given directories
// and returns a name-indexed registry for the skill tool.
// Directories that fail to load are logged and skipped.
func registerSkills(holder Context, conf Configuration, container agent.ToolContainer, dirs ...string) {
	registry := make(map[string]agent_skill.Skill)
	for _, dir := range dirs {
		props, err := agent_skill.ReadProperties(dir)
		if err != nil {
			slog.Warn("failed to load skill for registry", "dir", dir, "error", err)
			continue
		}
		registry[props.Name] = *props
		slog.Info("loaded skill into registry", "name", props.Name, "dir", dir)
	}
	skill_tool.AddToContainer(container, registry)
}

func registerSubAgents(holder Context, conf Configuration, container agent.ToolContainer, dirs ...string) {
	for _, path := range dirs {
		properties, err := agent_skill.ReadProperties(path)
		if err != nil {
			slog.Error("failed to register subagent", "path", path, "error", err)
			return
		}
		slog.Info("registered subagent", "path", path, "name", properties)
		subagent_tool.AddSubagentToContainer(container, *properties, holder.GetAgentMessageProducer())
	}
}
