package configuration

import (
	"log/slog"

	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	finance_stock_price_tool "github.com/hjwalt/platform/agent/tool/finance_stock_price"
	linux_shell_tool "github.com/hjwalt/platform/agent/tool/linux_shell"
	memory_tool "github.com/hjwalt/platform/agent/tool/memory"
	skill_tool "github.com/hjwalt/platform/agent/tool/skill"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	web_search_tool "github.com/hjwalt/platform/agent/tool/web_search"
	tool_container "github.com/hjwalt/platform/agent/util/container"
)

func RegisterTools(holder Context, conf Configuration) {
	container := tool_container.New()

	// register tools
	web_search_tool.AddToContainer(container, conf.Tool.WebSearch)
	linux_shell_tool.AddToContainer(container, conf.Tool.Shell)
	web_fetch_tool.AddToContainer(container, conf.Tool.WebFetch, holder.GetParserModel())
	finance_fx_price_tool.AddToContainer(container, conf.Tool.FxPrice)
	finance_stock_price_tool.AddToContainer(container, conf.Tool.StockPrice)
	for _, memoryConfig := range conf.Tool.Memory {
		memoryErr := memory_tool.AddToContainer(container, memoryConfig)
		if memoryErr != nil {
			slog.Error("failed to register memory tool set", "config", memoryConfig, "error", memoryErr)
		}
	}

	// register skills
	TryRegisterSkill(holder, conf, container, "./skills/researcher-agent")

	holder.SetToolContainer(container)
}

func TryRegisterSkill(holder Context, conf Configuration, container agent.ToolContainer, path string) {
	properties, err := agent_skill.ReadProperties(path)
	if err != nil {
		slog.Error("failed to register skill", "path", path, "error", err)
		return
	}
	slog.Info("registered skill", "path", path, "name", properties)
	skill_tool.AddSkillToContainer(container, *properties, holder.GetAgentMessageProducer())
}
