package example

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/flow"
)

type Context struct {
	context.Context
	Chat                 agent.LanguageModel
	RagStore             rag.Store
	AgentMessageProducer flow.Producer[agent.Message]
}
