package configuration

import (
	brave_search_web_tool "github.com/hjwalt/platform/agent/tool/brave_search_web"
	shell_tool "github.com/hjwalt/platform/agent/tool/shell"
	agent_skill "github.com/hjwalt/platform/agent/tool/skill"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	tool_container "github.com/hjwalt/platform/agent/util/container"
)

func RegisterTools(holder Context, conf Configuration) {
	container := tool_container.New()

	brave_search_web_tool.AddToContainer(container, conf.Tool.BraveSearch)
	shell_tool.AddToContainer(container, conf.Tool.Shell)
	web_fetch_tool.AddToContainer(container, conf.Tool.WebFetch)
	agent_skill.AddToContainer(container, conf.Tool.ResearchAgent, holder.GetAgentMessageProducer())

	holder.SetToolContainer(container)
}
