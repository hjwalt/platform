package web_fetch_tool

import (
	"context"
	"testing"

	"github.com/hjwalt/platform/agent"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

// stubLanguageModel implements agent.LanguageModel for testing.
type stubLanguageModel struct {
	chatResult []agent.Message
	chatError  error
}

func (s *stubLanguageModel) Chat(_ context.Context, _ []agent.Message, _ []string) ([]agent.Message, error) {
	return s.chatResult, s.chatError
}

func (s *stubLanguageModel) Start() error {
	return nil
}

func (s *stubLanguageModel) Stop() {}

func TestWebFetchName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	assert.Equal("web_fetch", tool.Name())
}

func TestWebFetchDescription(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	assert.NotEmpty(tool.Description())
}

func TestWebFetchAuto(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	assert.False(tool.Auto())
}

func TestWebFetchRequestSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	schema := tool.RequestSchema()
	assert.NotNil(schema)
}

func TestWebFetchResultSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	schema := tool.ResultSchema()
	assert.NotNil(schema)
}

func TestWebFetchRequestFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	assert.NotNil(tool.RequestFormat())
}

func TestWebFetchResultFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	assert.NotNil(tool.ResultFormat())
}

func TestWebFetchDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	desc := tool.DescribeRequest(Request{Link: "https://example.com"})

	assert.Contains(desc, "https://example.com")
	assert.Contains(desc, "fetch")
}

func TestWebFetchDescribeResult(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	desc := tool.DescribeResult(Response{
		Html:   "<html><body>test</body></html>",
		Parsed: "test content",
	})

	assert.Equal("test content", desc)
}

func TestWebFetchDescribeResultEmpty(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	desc := tool.DescribeResult(Response{Html: "", Parsed: ""})

	assert.Empty(desc)
}

func TestWebFetchCreateReturnsSyncTool(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{}, &stubLanguageModel{})

	var _ agent.SyncTool[Request, Response] = tool
	assert.NotNil(tool)
}

func TestAddToContainerRegistersWebFetchTool(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	AddToContainer(container, Configuration{}, &stubLanguageModel{})

	assert.True(container.Exists(agent.ToolCall{Name: "web_fetch"}))
}
