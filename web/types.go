package web

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow"
)

type Context struct {
	context.Context
	Chat                 agent.LanguageModel
	AgentHarnessStore    flow.Store[harness.ExecutionState]
	AgentMessageProducer flow.Producer[agent.Message]
}
