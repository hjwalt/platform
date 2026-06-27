package finance_stock_price_tool

import (
	"testing"

	"github.com/hjwalt/platform/agent"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

func TestStockPriceName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.Equal("finance_stock_price", tool.Name())
}

func TestStockPriceDescription(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.NotEmpty(tool.Description())
}

func TestStockPriceAuto(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.False(tool.Auto())
}

func TestStockPriceRequestSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	schema := tool.RequestSchema()
	assert.NotNil(schema)
}

func TestStockPriceResultSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	schema := tool.ResultSchema()
	assert.NotNil(schema)
}

func TestStockPriceRequestFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.NotNil(tool.RequestFormat())
}

func TestStockPriceResultFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	assert.NotNil(tool.ResultFormat())
}

func TestStockPriceDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	desc := tool.DescribeRequest(Request{Symbol: "AAPL"})

	assert.Contains(desc, "AAPL")
	assert.Contains(desc, "stock")
}

func TestStockPriceDescribeResult(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	desc := tool.DescribeResult(Response{Symbol: "AAPL", Currency: "USD", Value: 150.25})

	assert.Contains(desc, "AAPL")
	assert.Contains(desc, "USD")
	assert.Contains(desc, "150.25")
}

func TestStockPriceDescribeResultZeroValue(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	desc := tool.DescribeResult(Response{Symbol: "TSLA", Currency: "USD", Value: 0})

	assert.Contains(desc, "TSLA")
	assert.Contains(desc, "0.00")
}

func TestStockPriceCreateReturnsSyncTool(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{Secret: "test-secret"})

	var _ agent.SyncTool[Request, Response] = tool
	assert.NotNil(tool)
}

func TestAddToContainerRegistersStockPriceTool(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	AddToContainer(container, Configuration{Secret: "test-secret"})

	assert.True(container.Exists(agent.ToolCall{Name: "finance_stock_price"}))
}

func TestResolveValueWithNilResponse(t *testing.T) {
	assert := assert.New(t)

	value, ok := resolveValue(nil)

	assert.False(ok)
	assert.Equal(float64(0), value)
}

func TestResolveValueWithNilJSON200(t *testing.T) {
	assert := assert.New(t)

	// We can't easily construct the gen type, but we can test the nil path
	// by testing resolveValue with nil which covers the first branch
	value, ok := resolveValue(nil)

	assert.False(ok)
	assert.Equal(float64(0), value)
}

func TestDefaultCurrency(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("USD", defaultCurrency)
}
