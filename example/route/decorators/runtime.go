package decorators

import (
	"net/http"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/flow"
)

type RuntimeDecorator struct {
	Chat                 agent.LanguageModel
	RagStore             rag.Store
	AgentMessageProducer flow.Producer[agent.Message]
}

func (d *RuntimeDecorator) Decorate(c example.Context, w http.ResponseWriter, r *http.Request) (example.Context, error) {
	return example.Context{
		Context:              c.Context,
		Chat:                 d.Chat,
		RagStore:             d.RagStore,
		AgentMessageProducer: d.AgentMessageProducer,
	}, nil
}
