package llm

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/runtime"
)

type LanguageModel interface {
	runtime.Runtime
	Chat(context.Context, []agent.Message) ([]agent.Message, error)
}
