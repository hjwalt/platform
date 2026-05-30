package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	file_store "github.com/hjwalt/platform/state/file"
)

func RegisterOpenAi(holder Context, conf Configuration) {
	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    conf.OpenAi.Model,
		Endpoint: conf.OpenAi.Endpoint,
		Secret:   conf.OpenAi.Secret,
		Tools:    holder.GetToolContainer(),
	})
	holder.Add(model)

	holder.SetLanguageModel(model)
}

func RegisterAgentHarnessStore(holder Context, conf Configuration) {
	store := file_store.New(conf.Flow.Store)
	holder.Add(store)
	holder.SetAgentHarnessStore(store)
}
