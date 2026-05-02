package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/openai/openai-go/v3"
)

func RegisterOpenAi(holder Context, conf Configuration) {
	schemas := make([]openai.ChatCompletionToolUnionParam, 0)
	for _, tool := range holder.GetTool() {
		schemas = append(schemas, tool.Schema())
	}
	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    conf.OpenAi.Model,
		Endpoint: conf.OpenAi.Endpoint,
		Secret:   conf.OpenAi.Secret,
		Tools:    schemas,
	})
	holder.Add(model)

	holder.SetLanguageModel(model)
}

func RegisterInMemoryRagMemory(holder Context, conf Configuration) {
	store := rag.Memory()
	holder.Add(store)
	holder.SetRagStore(store)
}
