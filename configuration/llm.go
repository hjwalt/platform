package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	tool_container "github.com/hjwalt/platform/agent/util/container"
	file_store "github.com/hjwalt/platform/state/file"
)

func RegisterParserModel(holder Context, conf Configuration) {
	model := llm.New(conf.Model.Parser, tool_container.New())
	holder.Add(model)

	holder.SetParserModel(model)
}

func RegisterAgentModel(holder Context, conf Configuration) {
	model := llm.New(conf.Model.Agent, holder.GetToolContainer())
	holder.Add(model)

	holder.SetAgentModel(model)
}

func RegisterAgentHarnessStore(holder Context, conf Configuration) {
	store := file_store.New(conf.Flow.Store)
	holder.Add(store)
	holder.SetAgentHarnessStore(store)
}
