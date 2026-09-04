package llm

import (
	"context"
	"testing"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/type/optional"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockToolContainer is a minimal agent.ToolContainer that never executes.
type mockToolContainer struct{}

func (m *mockToolContainer) AddSync(agent.SyncToolWrapper)   {}
func (m *mockToolContainer) AddAsync(agent.AsyncToolWrapper) {}

func (m *mockToolContainer) Execute(ctx context.Context, in agent.Message, call agent.ToolCall) (optional.Optional[string], error) {
	return optional.Empty[string](), nil
}

func (m *mockToolContainer) DescribeRequest(call agent.ToolCall) (string, error) {
	return "", nil
}

func (m *mockToolContainer) Exists(call agent.ToolCall) bool { return false }
func (m *mockToolContainer) Auto(call agent.ToolCall) bool   { return false }

func (m *mockToolContainer) OpenAiParamsFiltered(allowed []string) []openai.ChatCompletionToolUnionParam {
	return nil
}

func (m *mockToolContainer) DeepSeekParams(allowed []string) []deepseek.Tool {
	return nil
}

func TestNewOpenAiReturnsConfiguredModel(t *testing.T) {
	assert := assert.New(t)

	tools := &mockToolContainer{}
	model := New(ModelConfig{
		Type:     OpenAi,
		Model:    "gpt-5",
		Endpoint: "https://api.openai.com/v1",
		Secret:   "sk-test",
	}, tools)

	require.NotNil(t, model)
	openAi, ok := model.(*openAiModel)
	require.True(t, ok)
	assert.Equal("gpt-5", openAi.Model)
	assert.Equal("https://api.openai.com/v1", openAi.Endpoint)
	assert.Equal("sk-test", openAi.Secret)
	assert.Equal(tools, openAi.Tools)
}

func TestNewDeepSeekReturnsConfiguredModel(t *testing.T) {
	assert := assert.New(t)

	tools := &mockToolContainer{}
	model := New(ModelConfig{
		Type:     DeepSeek,
		Model:    "deepseek-chat",
		Endpoint: "https://api.deepseek.com",
		Secret:   "ds-test",
	}, tools)

	require.NotNil(t, model)
	deepSeek, ok := model.(*deepSeekModel)
	require.True(t, ok)
	assert.Equal("deepseek-chat", deepSeek.Model)
	assert.Equal("ds-test", deepSeek.Secret)
	assert.Equal(tools, deepSeek.Tools)
}

func TestNewUnknownTypeDefaultsToOpenAi(t *testing.T) {
	assert := assert.New(t)

	model := New(ModelConfig{Type: ModelType(99), Model: "m", Secret: "s"}, &mockToolContainer{})

	require.NotNil(t, model)
	_, ok := model.(*openAiModel)
	assert.True(ok, "unknown model type should fall back to openAi")
}

func TestNewIsLanguageModelAndRuntime(t *testing.T) {
	assert := assert.New(t)

	var languageModel agent.LanguageModel = New(ModelConfig{Type: DeepSeek, Secret: "s"}, &mockToolContainer{})
	assert.NotNil(languageModel)

	var rt interface {
		Start() error
		Stop()
	} = languageModel
	assert.NotNil(rt)
}

func TestOpenAiModelStartStopNoNetwork(t *testing.T) {
	model := &openAiModel{
		Model:    "gpt-5",
		Endpoint: "http://127.0.0.1:9",
		Secret:   "sk-test",
		Tools:    &mockToolContainer{},
	}

	require.NoError(t, model.Start())
	model.Stop()
	model.Stop() // idempotent
}

func TestDeepSeekModelStartStopNoNetwork(t *testing.T) {
	model := &deepSeekModel{
		Model:  "deepseek-chat",
		Secret: "ds-test",
		Tools:  &mockToolContainer{},
	}

	require.NoError(t, model.Start())
	model.Stop()
	model.Stop() // idempotent
}

func TestNewEmbeddingOpenAiType(t *testing.T) {
	assert := assert.New(t)

	embedding := NewEmbedding(EmbeddingConfig{
		Type:     OpenAi,
		Model:    "text-embedding-3-small",
		Endpoint: "https://api.openai.com/v1",
		Secret:   "sk-test",
	})

	require.NotNil(t, embedding)
	_, ok := embedding.(*openAiEmbedding)
	assert.True(ok)
}

func TestNewEmbeddingDefaultTypeIsOpenAi(t *testing.T) {
	assert := assert.New(t)

	embedding := NewEmbedding(EmbeddingConfig{Type: ModelType(77), Model: "m"})

	require.NotNil(t, embedding)
	var agentEmbedding agent.Embedding = embedding
	assert.NotNil(agentEmbedding)
}
