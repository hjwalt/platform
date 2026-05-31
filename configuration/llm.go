package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	tool_container "github.com/hjwalt/platform/agent/util/container"
	file_store "github.com/hjwalt/platform/state/file"
)

func RegisterParserModel(holder Context, conf Configuration) {
	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    conf.Model.Parser.Model,
		Endpoint: conf.Model.Parser.Endpoint,
		Secret:   conf.Model.Parser.Secret,
		Tools:    tool_container.New(),
	})
	holder.Add(model)

	holder.SetParserModel(model)
}

func RegisterAgentModel(holder Context, conf Configuration) {
	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    conf.Model.Agent.Model,
		Endpoint: conf.Model.Agent.Endpoint,
		Secret:   conf.Model.Agent.Secret,
		Tools:    holder.GetToolContainer(),
	})
	holder.Add(model)

	holder.SetAgentModel(model)
}

func RegisterAgentHarnessStore(holder Context, conf Configuration) {
	store := file_store.New(conf.Flow.Store)
	holder.Add(store)
	holder.SetAgentHarnessStore(store)
}
