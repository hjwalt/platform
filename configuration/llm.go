package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/state/state_file"
)

func RegisterOpenAi(holder Context, conf Configuration) {
	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    conf.OpenAi.Model,
		Endpoint: conf.OpenAi.Endpoint,
		Secret:   conf.OpenAi.Secret,
		Tools:    holder.GetToolContainer().AsToolMap(),
	})
	holder.Add(model)

	holder.SetLanguageModel(model)
}

func RegisterAgentHarnessStore(holder Context, conf Configuration) {
	store := &state_file.FileStore{
		Path: "/home/hjwalt/Projects/tmp/platform/",
	}
	holder.Add(store)
	holder.SetAgentHarnessStore(store)
}
