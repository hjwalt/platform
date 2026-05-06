package decorator

import (
	"net/http"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/web"
)

type RuntimeDecorator struct {
	Chat                 agent.LanguageModel
	AgentMessageProducer flow.Producer[agent.Message]
	AgentHarnessStore    flow.Store[harness.ExecutionState]
	Tool                 map[string]agent.Tool
}

func (d *RuntimeDecorator) Decorate(c web.Context, w http.ResponseWriter, r *http.Request) (web.Context, error) {
	return web.Context{
		Context:              c.Context,
		Chat:                 d.Chat,
		AgentMessageProducer: d.AgentMessageProducer,
		AgentHarnessStore:    d.AgentHarnessStore,
		Tool:                 d.Tool,
	}, nil
}
