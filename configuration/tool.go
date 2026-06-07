package configuration

import (
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	linux_shell_tool "github.com/hjwalt/platform/agent/tool/linux_shell"
	skill_tool "github.com/hjwalt/platform/agent/tool/skill"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	web_search_tool "github.com/hjwalt/platform/agent/tool/web_search"
	tool_container "github.com/hjwalt/platform/agent/util/container"
)

func RegisterTools(holder Context, conf Configuration) {
	container := tool_container.New()

	web_search_tool.AddToContainer(container, conf.Tool.WebSearch)
	linux_shell_tool.AddToContainer(container, conf.Tool.Shell)
	web_fetch_tool.AddToContainer(container, conf.Tool.WebFetch, holder.GetParserModel())
	skill_tool.AddToContainer(container, conf.Tool.ResearchAgent, holder.GetAgentMessageProducer())
	finance_fx_price_tool.AddToContainer(container, conf.Tool.FxPrice)

	holder.SetToolContainer(container)
}
