package configuration

import (
	"errors"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/state"
	"github.com/hjwalt/platform/type/optional"
)

// This typed structure is better for configuration wiring because it prevents random runtime failures

type Context interface {
	Add(runtimes ...runtime.Runtime)
	SetToolContainer(agent.ToolContainer)
	GetToolContainer() agent.ToolContainer
	SetSkillContainer(agent.SkillContainer)
	GetSkillContainer() agent.SkillContainer
	SetAgentModel(agent.LanguageModel)
	GetAgentModel() agent.LanguageModel
	SetParserModel(agent.LanguageModel)
	GetParserModel() agent.LanguageModel
	SetKafkaProducer(message.Producer[kafka.KafkaMetadata])
	GetKafkaProducer() message.Producer[kafka.KafkaMetadata]
	SetAgentMessageProducer(flow.Producer[agent.Message])
	GetAgentMessageProducer() flow.Producer[agent.Message]
	SetAgentHarnessStore(state.Store)
	GetAgentHarnessStore() state.Store
	SetMemoryStore(state.Store)
	GetMemoryStore() state.Store
	Block()
}

func ContextBuilder() Context {
	return &holder{
		Runtimes:             make([]runtime.Runtime, 0),
		ToolContainer:        optional.Empty[agent.ToolContainer](),
		SkillContainer:       optional.Empty[agent.SkillContainer](),
		AgentModel:           optional.Empty[agent.LanguageModel](),
		ParserModel:          optional.Empty[agent.LanguageModel](),
		KafkaProducer:        optional.Empty[message.Producer[kafka.KafkaMetadata]](),
		AgentMessageProducer: optional.Empty[flow.Producer[agent.Message]](),
		AgentHarnessStore:    optional.Empty[state.Store](),
		MemoryStore:          optional.Empty[state.Store](),
	}
}

type holder struct {
	Runtimes             []runtime.Runtime
	ToolContainer        optional.Optional[agent.ToolContainer]
	SkillContainer       optional.Optional[agent.SkillContainer]
	AgentModel           optional.Optional[agent.LanguageModel]
	ParserModel          optional.Optional[agent.LanguageModel]
	KafkaProducer        optional.Optional[message.Producer[kafka.KafkaMetadata]]
	AgentMessageProducer optional.Optional[flow.Producer[agent.Message]]
	AgentHarnessStore    optional.Optional[state.Store]
	MemoryStore          optional.Optional[state.Store]
}

func (r *holder) Add(runtimes ...runtime.Runtime) {
	r.Runtimes = append(r.Runtimes, runtimes...)
}

func (r *holder) SetToolContainer(value agent.ToolContainer) {
	r.ToolContainer = optional.Of(value)
}

func (r *holder) GetToolContainer() agent.ToolContainer {
	if !r.ToolContainer.IsPresent() {
		r.Missing()
	}
	return r.ToolContainer.Get()
}

func (r *holder) SetSkillContainer(value agent.SkillContainer) {
	r.SkillContainer = optional.Of(value)
}

func (r *holder) GetSkillContainer() agent.SkillContainer {
	if !r.SkillContainer.IsPresent() {
		r.Missing()
	}
	return r.SkillContainer.Get()
}

func (r *holder) SetAgentModel(value agent.LanguageModel) {
	r.AgentModel = optional.Of(value)
}

func (r *holder) GetAgentModel() agent.LanguageModel {
	if !r.AgentModel.IsPresent() {
		r.Missing()
	}
	return r.AgentModel.Get()
}

func (r *holder) SetParserModel(value agent.LanguageModel) {
	r.ParserModel = optional.Of(value)
}

func (r *holder) GetParserModel() agent.LanguageModel {
	if !r.ParserModel.IsPresent() {
		r.Missing()
	}
	return r.ParserModel.Get()
}

func (r *holder) SetKafkaProducer(value message.Producer[kafka.KafkaMetadata]) {
	r.KafkaProducer = optional.Of(value)
}

func (r *holder) GetKafkaProducer() message.Producer[kafka.KafkaMetadata] {
	if !r.KafkaProducer.IsPresent() {
		r.Missing()
	}
	return r.KafkaProducer.Get()
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

func (r *holder) SetMemoryStore(value state.Store) {
	r.MemoryStore = optional.Of(value)
}

func (r *holder) GetMemoryStore() state.Store {
	if !r.MemoryStore.IsPresent() {
		r.Missing()
	}
	return r.MemoryStore.Get()
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
