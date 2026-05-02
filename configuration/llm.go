package configuration

import (
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/runtime"
	"github.com/openai/openai-go/v3"
)

func RegisterOpenAi(holder runtime.Holder, conf Configuration, tools []agent.Tool) agent.LanguageModel {
	schemas := make([]openai.ChatCompletionToolUnionParam, 0)
	for _, tool := range tools {
		schemas = append(schemas, tool.Schema())
	}
	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    conf.OpenAi.Model,
		Endpoint: conf.OpenAi.Endpoint,
		Secret:   conf.OpenAi.Secret,
		Tools:    schemas,
	})
	holder.Add(model)
	return model
}

func RegisterInMemoryRagMemory(holder runtime.Holder, conf Configuration) rag.Store {
	store := rag.Memory()
	holder.Add(store)
	return store
}

func RegisterRagModel(holder runtime.Holder, conf Configuration, model agent.LanguageModel, store rag.Store) agent.LanguageModel {
	ragModel := rag.Rag(model, store)
	holder.Add(ragModel)
	return ragModel
}
