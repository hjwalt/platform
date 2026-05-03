package configuration

import (
	"errors"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/type/optional"
)

type Context interface {
	Add(runtimes ...runtime.Runtime)
	AddTool(tools ...agent.Tool)
	GetTool() map[string]agent.Tool
	SetLanguageModel(agent.LanguageModel)
	GetLanguageModel() agent.LanguageModel
	SetRagStore(rag.Store)
	GetRagStore() rag.Store
	SetAgentMessageProducer(flow.Producer[agent.Message])
	GetAgentMessageProducer() flow.Producer[agent.Message]
	Block()
}

func ContextBuilder() Context {
	return &holder{
		Runtimes:             make([]runtime.Runtime, 0),
		Tools:                make(map[string]agent.Tool, 0),
		RagModel:             optional.Empty[agent.LanguageModel](),
		RagStore:             optional.Empty[rag.Store](),
		AgentMessageProducer: optional.Empty[flow.Producer[agent.Message]](),
	}
}

type holder struct {
	Runtimes             []runtime.Runtime
	Tools                map[string]agent.Tool
	RagModel             optional.Optional[agent.LanguageModel]
	RagStore             optional.Optional[rag.Store]
	AgentMessageProducer optional.Optional[flow.Producer[agent.Message]]
}

func (r *holder) Add(runtimes ...runtime.Runtime) {
	r.Runtimes = append(r.Runtimes, runtimes...)
}

func (r *holder) AddTool(tools ...agent.Tool) {
	for _, tool := range tools {
		r.Tools[tool.Name()] = tool
	}
}

func (r *holder) GetTool() map[string]agent.Tool {
	return r.Tools
}

func (r *holder) SetLanguageModel(value agent.LanguageModel) {
	r.RagModel = optional.Of(value)
}

func (r *holder) GetLanguageModel() agent.LanguageModel {
	if !r.RagModel.IsPresent() {
		r.Missing()
	}
	return r.RagModel.Get()
}

func (r *holder) SetRagStore(value rag.Store) {
	r.RagStore = optional.Of(value)
}

func (r *holder) GetRagStore() rag.Store {
	if !r.RagStore.IsPresent() {
		r.Missing()
	}
	return r.RagStore.Get()
}

func (r *holder) SetAgentMessageProducer(value flow.Producer[agent.Message]) {
	r.AgentMessageProducer = optional.Of(value)
}

func (r *holder) GetAgentMessageProducer() flow.Producer[agent.Message] {
	if !r.AgentMessageProducer.IsPresent() {
		r.Missing()
	}
	return r.AgentMessageProducer.Get()
}

func (r *holder) Missing() {
	panic(errors.New("missing data required"))
}

func (r *holder) Block() {
	startErr := runtime.Start(
		r.Runtimes,
		time.Second,
	)

	if startErr != nil {
		panic(startErr)
	}

	slog.Info("started")

	runtime.Wait()

	slog.Info("stopped")
}
