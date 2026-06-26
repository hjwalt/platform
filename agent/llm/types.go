package llm

import "github.com/hjwalt/platform/agent"

type ModelType int

const (
	OpenAi ModelType = iota
	DeepSeek
)

type ModelConfig struct {
	Type     ModelType
	Model    string
	Endpoint string
	Secret   string
}

func New(config ModelConfig, tools agent.ToolContainer) agent.LanguageModel {
	switch config.Type {
	case OpenAi:
		return createOpenAi(config, tools)
	case DeepSeek:
		return createDeepSeek(config, tools)
	default:
		return createOpenAi(config, tools)
	}
}

type EmbeddingConfig struct {
	Type       ModelType
	Model      string
	Endpoint   string
	Secret     string
	Dimensions int // 0 = use model default
}

func NewEmbedding(config EmbeddingConfig) agent.Embedding {
	switch config.Type {
	case OpenAi:
		return createOpenAiEmbedding(config)
	default:
		return createOpenAiEmbedding(config)
	}
}
