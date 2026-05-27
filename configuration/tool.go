package configuration

import (
	brave_search_web_tool "github.com/hjwalt/platform/agent/tool/brave_search_web"
	fork_tool "github.com/hjwalt/platform/agent/tool/fork"
	shell_tool "github.com/hjwalt/platform/agent/tool/shell"
	tool_container "github.com/hjwalt/platform/agent/util/container"
)

func RegisterTools(holder Context, conf Configuration) {
	container := tool_container.New()

	brave_search_web_tool.AddToContainer(container, conf.Tool.BraveSearch)
	shell_tool.AddToContainer(container, conf.Tool.Shell)
	fork_tool.AddToContainer(container, conf.Tool.ResearchAgent, holder.GetAgentMessageProducer())

	holder.SetToolContainer(container)
}
