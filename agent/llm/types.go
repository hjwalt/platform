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
