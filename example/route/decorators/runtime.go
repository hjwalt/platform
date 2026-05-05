package decorators

import (
	"net/http"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/flow"
)

type RuntimeDecorator struct {
	Chat                 agent.LanguageModel
	AgentMessageProducer flow.Producer[agent.Message]
	AgentHarnessStore    flow.Store[harness.ExecutionState]
	Tool                 map[string]agent.Tool
}

func (d *RuntimeDecorator) Decorate(c example.Context, w http.ResponseWriter, r *http.Request) (example.Context, error) {
	return example.Context{
		Context:              c.Context,
		Chat:                 d.Chat,
		AgentMessageProducer: d.AgentMessageProducer,
		AgentHarnessStore:    d.AgentHarnessStore,
		Tool:                 d.Tool,
	}, nil
}
