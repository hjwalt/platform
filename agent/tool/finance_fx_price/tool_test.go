package finance_fx_price_tool

import (
	"testing"

	"github.com/hjwalt/platform/agent"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

func TestFxPriceName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.Equal("finance_fx_price", tool.Name())
}

func TestFxPriceDescription(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.NotEmpty(tool.Description())
}

func TestFxPriceAuto(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.False(tool.Auto())
}

func TestFxPriceRequestSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	schema := tool.RequestSchema()
	assert.NotNil(schema)
}

func TestFxPriceResultSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	schema := tool.ResultSchema()
	assert.NotNil(schema)
}

func TestFxPriceRequestFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.NotNil(tool.RequestFormat())
}

func TestFxPriceResultFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.NotNil(tool.ResultFormat())
}

func TestFxPriceDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	desc := tool.DescribeRequest(Request{Base: "EUR", Quote: "USD"})

	assert.Contains(desc, "EUR")
	assert.Contains(desc, "USD")
	assert.Contains(desc, "FX")
}

func TestFxPriceDescribeResult(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	desc := tool.DescribeResult(Response{Base: "EUR", Quote: "USD", Value: 1.1234})

	assert.Contains(desc, "EUR")
	assert.Contains(desc, "USD")
	assert.Contains(desc, "FX")
	assert.Contains(desc, "1.12")
}

func TestFxPriceDescribeResultWithZeroValue(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	desc := tool.DescribeResult(Response{Base: "GBP", Quote: "JPY", Value: 0})

	assert.Contains(desc, "GBP")
	assert.Contains(desc, "JPY")
	assert.Contains(desc, "0.00")
}

func TestFxPriceCreateReturnsSyncTool(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	// Verify it implements SyncTool interface
	var _ agent.SyncTool[Request, Response] = tool
	assert.NotNil(tool)
}

func TestAddToContainerRegistersFxPriceTool(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	AddToContainer(container, Configuration{Secret: "test-secret"})

	assert.True(container.Exists(agent.ToolCall{Name: "finance_fx_price"}))
}
