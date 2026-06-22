package configuration

import (
	"log/slog"

	"github.com/hjwalt/platform/agent"
	agent_skill "github.com/hjwalt/platform/agent/skill"
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	finance_stock_price_tool "github.com/hjwalt/platform/agent/tool/finance_stock_price"
	linux_shell_tool "github.com/hjwalt/platform/agent/tool/linux_shell"
	memory_tool "github.com/hjwalt/platform/agent/tool/memory"
	subagent_tool "github.com/hjwalt/platform/agent/tool/subagent"
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
		if validateErr := memory_tool.Validate(memoryConfig, holder.GetMemoryStore()); validateErr != nil {
			slog.Error("failed to register memory tool", "key", memoryConfig.Key, "error", validateErr)
			continue
		}
		memory_tool.AddToContainer(container, memoryConfig, holder.GetMemoryStore())
	}

	// register skills
	TryRegisterSubAgent(holder, conf, container, "./docs/agents/researcher-agent")

	holder.SetToolContainer(container)
}

func TryRegisterSubAgent(holder Context, conf Configuration, container agent.ToolContainer, path string) {
	properties, err := agent_skill.ReadProperties(path)
	if err != nil {
		slog.Error("failed to register subagent", "path", path, "error", err)
		return
	}
	slog.Info("registered subagent", "path", path, "name", properties)
	subagent_tool.AddSubagentToContainer(container, *properties, holder.GetAgentMessageProducer())
}
