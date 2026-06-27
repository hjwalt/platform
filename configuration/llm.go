package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	file_store "github.com/hjwalt/platform/state/file"
)

func RegisterParserModel(holder Context, conf Configuration) {
	model := llm.New(conf.Model.Configurations[conf.Model.Parser], harness_container.NewToolContainer())
	holder.Add(model)

	holder.SetParserModel(model)
}

func RegisterAgentModel(holder Context, conf Configuration) {
	model := llm.New(conf.Model.Configurations[conf.Model.Agent], holder.GetToolContainer())
	holder.Add(model)

	holder.SetAgentModel(model)
}

func RegisterAgentHarnessStore(holder Context, conf Configuration) {
	store := file_store.New(conf.Store.Agent)
	holder.Add(store)
	holder.SetAgentHarnessStore(store)
}

func RegisterMemoryStore(holder Context, conf Configuration) {
	store := file_store.New(conf.Store.Memory)
	holder.Add(store)
	holder.SetMemoryStore(store)
}
