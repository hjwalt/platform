package configuration

import (
	"errors"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/agent"
	agent_tool "github.com/hjwalt/platform/agent/tool"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/state"
	"github.com/hjwalt/platform/type/optional"
)

type Context interface {
	Add(runtimes ...runtime.Runtime)
	SetToolContainer(agent_tool.Container)
	GetToolContainer() agent_tool.Container
	SetLanguageModel(agent.LanguageModel)
	GetLanguageModel() agent.LanguageModel
	SetAgentMessageProducer(flow.Producer[agent.Message])
	GetAgentMessageProducer() flow.Producer[agent.Message]
	SetAgentHarnessStore(state.Store)
	GetAgentHarnessStore() state.Store
	Block()
}

func ContextBuilder() Context {
	return &holder{
		Runtimes:             make([]runtime.Runtime, 0),
		ToolContainer:        optional.Empty[agent_tool.Container](),
		RagModel:             optional.Empty[agent.LanguageModel](),
		AgentMessageProducer: optional.Empty[flow.Producer[agent.Message]](),
		AgentHarnessStore:    optional.Empty[state.Store](),
	}
}

type holder struct {
	Runtimes             []runtime.Runtime
	ToolContainer        optional.Optional[agent_tool.Container]
	RagModel             optional.Optional[agent.LanguageModel]
	AgentMessageProducer optional.Optional[flow.Producer[agent.Message]]
	AgentHarnessStore    optional.Optional[state.Store]
}

func (r *holder) Add(runtimes ...runtime.Runtime) {
	r.Runtimes = append(r.Runtimes, runtimes...)
}

func (r *holder) SetToolContainer(value agent_tool.Container) {
	r.ToolContainer = optional.Of(value)
}

func (r *holder) GetToolContainer() agent_tool.Container {
	if !r.ToolContainer.IsPresent() {
		r.Missing()
	}
	return r.ToolContainer.Get()
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

func (r *holder) SetAgentMessageProducer(value flow.Producer[agent.Message]) {
	r.AgentMessageProducer = optional.Of(value)
}

func (r *holder) GetAgentMessageProducer() flow.Producer[agent.Message] {
	if !r.AgentMessageProducer.IsPresent() {
		r.Missing()
	}
	return r.AgentMessageProducer.Get()
}

func (r *holder) SetAgentHarnessStore(value state.Store) {
	r.AgentHarnessStore = optional.Of(value)
}

func (r *holder) GetAgentHarnessStore() state.Store {
	if !r.AgentHarnessStore.IsPresent() {
		r.Missing()
	}
	return r.AgentHarnessStore.Get()
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
		runtime.Stop()
		panic(startErr)
	}

	slog.Info("started")

	runtime.Wait()

	slog.Info("stopped")
}
