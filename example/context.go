package example

import (
	"context"

	"github.com/hjwalt/platform/agent/llm"
)

type Context struct {
	context.Context
	Chat llm.LanguageModel
}
