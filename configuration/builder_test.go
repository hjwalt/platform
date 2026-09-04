package configuration

import (
	"context"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/hjwalt/platform/state"
	"github.com/hjwalt/platform/type/optional"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"
)

// mocks used to exercise holder wiring without any real runtime/io

type mockRuntime struct{}

func (r *mockRuntime) Start() error { return nil }
func (r *mockRuntime) Stop()        {}

type mockLanguageModel struct{}

func (m *mockLanguageModel) Start() error { return nil }
func (m *mockLanguageModel) Stop()        {}
func (m *mockLanguageModel) Chat(ctx context.Context, msgs []agent.Message, tools []string) ([]agent.Message, error) {
	return nil, nil
}

type mockToolContainer struct{}

func (c *mockToolContainer) AddSync(tool agent.SyncToolWrapper) {}
func (c *mockToolContainer) AddAsync(tool agent.AsyncToolWrapper) {
}
func (c *mockToolContainer) Execute(ctx context.Context, msg agent.Message, call agent.ToolCall) (optional.Optional[string], error) {
	return optional.Empty[string](), nil
}
func (c *mockToolContainer) DescribeRequest(call agent.ToolCall) (string, error) {
	return "", nil
}
func (c *mockToolContainer) Exists(call agent.ToolCall) bool { return false }
func (c *mockToolContainer) Auto(call agent.ToolCall) bool   { return false }
func (c *mockToolContainer) OpenAiParamsFiltered(allowed []string) []openai.ChatCompletionToolUnionParam {
	return nil
}
func (c *mockToolContainer) DeepSeekParams(allowed []string) []deepseek.Tool {
	return nil
}

type mockSkillContainer struct{}

func (c *mockSkillContainer) Add(in agent.Instruction) {}
func (c *mockSkillContainer) Get(name string) (agent.Instruction, bool) {
	return agent.Instruction{}, false
}
func (c *mockSkillContainer) Assistant(ctx string) agent.Message { return agent.Message{} }

type mockKafkaProducer struct{}

func (p *mockKafkaProducer) Produce(ctx context.Context, msgs []message.Message[kafka.KafkaMetadata]) error {
	return nil
}
func (p *mockKafkaProducer) Start() error { return nil }
func (p *mockKafkaProducer) Stop()        {}

type mockAgentMessageProducer struct{}

func (p *mockAgentMessageProducer) ProduceMessage(ctx context.Context, msgs []flow.Message[agent.Message]) error {
	return nil
}
func (p *mockAgentMessageProducer) Produce(ctx context.Context, values []agent.Message) error {
	return nil
}
func (p *mockAgentMessageProducer) Start() error { return nil }
func (p *mockAgentMessageProducer) Stop()        {}

type mockStateStore struct{}

func (s *mockStateStore) Read(ctx context.Context, id string) (state.State, error) {
	return state.State{Id: id}, nil
}
func (s *mockStateStore) Write(ctx context.Context, st state.State) error { return nil }
func (s *mockStateStore) Delete(ctx context.Context, id string) error     { return nil }
func (s *mockStateStore) Keys(ctx context.Context) ([]string, error)      { return nil, nil }
func (s *mockStateStore) Start() error                                    { return nil }
func (s *mockStateStore) Stop()                                           {}

func TestContextAddAppendsRuntimes(t *testing.T) {
	assert := assert.New(t)

	ctx := ContextBuilder()
	runtimeOne := &mockRuntime{}
	runtimeTwo := &mockRuntime{}

	ctx.Add(runtimeOne)
	ctx.Add(runtimeTwo, runtimeOne)

	holder, ok := ctx.(*holder)
	assert.True(ok)
	assert.Len(holder.Runtimes, 3)
	assert.Same(runtimeOne, holder.Runtimes[0])
	assert.Same(runtimeTwo, holder.Runtimes[1])
	assert.Same(runtimeOne, holder.Runtimes[2])
}

func TestContextGettersPanicBeforeSet(t *testing.T) {
	assert := assert.New(t)

	ctx := ContextBuilder()

	cases := []struct {
		name string
		get  func()
	}{
		{"GetToolContainer", func() { ctx.GetToolContainer() }},
		{"GetSkillContainer", func() { ctx.GetSkillContainer() }},
		{"GetAgentModel", func() { ctx.GetAgentModel() }},
		{"GetParserModel", func() { ctx.GetParserModel() }},
		{"GetKafkaProducer", func() { ctx.GetKafkaProducer() }},
		{"GetAgentMessageProducer", func() { ctx.GetAgentMessageProducer() }},
		{"GetAgentHarnessStore", func() { ctx.GetAgentHarnessStore() }},
		{"GetMemoryStore", func() { ctx.GetMemoryStore() }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.PanicsWithError("missing data required", tc.get)
		})
	}
}

func TestContextSetGetRoundTrip(t *testing.T) {
	assert := assert.New(t)

	ctx := ContextBuilder()

	toolContainer := &mockToolContainer{}
	skillContainer := &mockSkillContainer{}
	agentModel := &mockLanguageModel{}
	parserModel := &mockLanguageModel{}
	kafkaProducer := &mockKafkaProducer{}
	agentMessageProducer := &mockAgentMessageProducer{}
	agentHarnessStore := &mockStateStore{}
	memoryStore := &mockStateStore{}

	ctx.SetToolContainer(toolContainer)
	ctx.SetSkillContainer(skillContainer)
	ctx.SetAgentModel(agentModel)
	ctx.SetParserModel(parserModel)
	ctx.SetKafkaProducer(kafkaProducer)
	ctx.SetAgentMessageProducer(agentMessageProducer)
	ctx.SetAgentHarnessStore(agentHarnessStore)
	ctx.SetMemoryStore(memoryStore)

	assert.Same(toolContainer, ctx.GetToolContainer())
	assert.Same(skillContainer, ctx.GetSkillContainer())
	assert.Same(agentModel, ctx.GetAgentModel())
	assert.Same(parserModel, ctx.GetParserModel())
	assert.Same(kafkaProducer, ctx.GetKafkaProducer())
	assert.Same(agentMessageProducer, ctx.GetAgentMessageProducer())
	assert.Same(agentHarnessStore, ctx.GetAgentHarnessStore())
	assert.Same(memoryStore, ctx.GetMemoryStore())
}

func TestContextTypedNilValues(t *testing.T) {
	assert := assert.New(t)

	// optional.Of wraps the (typed nil) value and marks it present, so the
	// getters must return without panicking.

	ctx := ContextBuilder()
	ctx.SetKafkaProducer((*mockKafkaProducer)(nil))
	assert.NotPanics(func() { ctx.GetKafkaProducer() })
	assert.Equal((*mockKafkaProducer)(nil), ctx.GetKafkaProducer())

	ctx.SetAgentMessageProducer((*mockAgentMessageProducer)(nil))
	assert.NotPanics(func() { ctx.GetAgentMessageProducer() })
	assert.Equal((*mockAgentMessageProducer)(nil), ctx.GetAgentMessageProducer())

	ctx.SetAgentHarnessStore((*mockStateStore)(nil))
	assert.NotPanics(func() { ctx.GetAgentHarnessStore() })
	assert.Equal((*mockStateStore)(nil), ctx.GetAgentHarnessStore())

	ctx.SetMemoryStore((*mockStateStore)(nil))
	assert.NotPanics(func() { ctx.GetMemoryStore() })
	assert.Equal((*mockStateStore)(nil), ctx.GetMemoryStore())

	// a plain (untyped) nil is also stored as present and returned as nil
	plainNilCtx := ContextBuilder()
	plainNilCtx.SetKafkaProducer(nil)
	assert.NotPanics(func() { plainNilCtx.GetKafkaProducer() })
	assert.Nil(plainNilCtx.GetKafkaProducer())
}
