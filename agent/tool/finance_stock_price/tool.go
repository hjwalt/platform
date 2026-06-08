package finance_stock_price_tool

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	tool_mcp "github.com/hjwalt/platform/agent/util/mcp"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/format"
	"github.com/massive-com/client-go/v3/rest"
	"github.com/massive-com/client-go/v3/rest/gen"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name            = "finance_stock_price"
	defaultCurrency = "USD"
)

type Configuration struct {
	Secret string
}

type Request struct {
	Symbol string `json:"symbol" jsonschema:"stock ticker symbol, for example AAPL"`
}

type Response struct {
	Symbol   string  `json:"symbol" jsonschema:"stock ticker symbol"`
	Currency string  `json:"currency" jsonschema:"quote currency of the requested symbol"`
	Value    float32 `json:"value" jsonschema:"stock price"`
}

type tool struct {
	client *rest.Client
}

func (t *tool) Apply(ctx context.Context, params Request) (Response, error) {
	apiParams := &gen.GetStocksSMAParams{
		Timespan:   rest.Ptr(gen.GetStocksSMAParamsTimespanMinute),
		Adjusted:   rest.Ptr(true),
		Window:     rest.Ptr(50),
		SeriesType: rest.Ptr(gen.GetStocksSMAParamsSeriesTypeClose),
		Order:      rest.Ptr(gen.GetStocksSMAParamsOrderDesc),
		Limit:      rest.Ptr(10),
	}

	res, err := t.client.GetStocksSMAWithResponse(ctx, params.Symbol, apiParams)
	if err != nil {
		return Response{Symbol: params.Symbol, Currency: defaultCurrency}, err
	}

	if res == nil || res.HTTPResponse == nil {
		return Response{Symbol: params.Symbol, Currency: defaultCurrency}, errors.New("failed, missing http response")
	}

	if res.HTTPResponse.StatusCode != 200 {
		return Response{Symbol: params.Symbol, Currency: defaultCurrency}, errors.New("failed, http response status is " + res.HTTPResponse.Status)
	}

	if res.JSON200 == nil || res.JSON200.Results.Values == nil {
		return Response{Symbol: params.Symbol, Currency: defaultCurrency}, errors.New("failed, no result values")
	}

	value, ok := resolveValue(res)
	if !ok {
		return Response{Symbol: params.Symbol, Currency: defaultCurrency}, errors.New("failed, no result values")
	}

	currency := resolveCurrencyFromReferenceTicker(ctx, t.client, params.Symbol)

	slog.Info("stock quote", "request", params, "value", value, "currency", currency)

	return Response{
		Symbol:   params.Symbol,
		Currency: currency,
		Value:    float32(value),
	}, nil
}

func resolveValue(res *gen.GetStocksSMAResponse) (float64, bool) {
	if res == nil || res.JSON200 == nil || res.JSON200.Results.Values == nil {
		return 0, false
	}

	values := *res.JSON200.Results.Values
	if len(values) == 0 {
		return 0, false
	}

	if values[0].Value == nil {
		return 0, false
	}

	return float64(*values[0].Value), true
}

func resolveCurrencyFromReferenceTicker(ctx context.Context, client *rest.Client, symbol string) string {
	res, err := client.GetTickerWithResponse(ctx, symbol, nil)
	if err != nil || res == nil || res.HTTPResponse == nil || res.HTTPResponse.StatusCode != 200 || res.JSON200 == nil || res.JSON200.Results == nil {
		return defaultCurrency
	}
	if res.JSON200.Results.CurrencyName == "" {
		return defaultCurrency
	}
	return strings.ToUpper(res.JSON200.Results.CurrencyName)
}

func (t *tool) Name() string {
	return Name
}

func (t *tool) Description() string {
	return "get stock quote price"
}

func (t *tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *tool) DescribeRequest(request Request) string {
	outputBuilder := strings.Builder{}

	outputBuilder.WriteString("getting stock quote for ")
	outputBuilder.WriteString(request.Symbol)

	return outputBuilder.String()
}

func (t *tool) ResultFormat() format.Format[Response] {
	return format.Json[Response]()
}

func (t *tool) ResultSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Response](opts)
	return toolSchema
}

func (t *tool) DescribeResult(response Response) string {
	outputBuilder := strings.Builder{}

	outputBuilder.WriteString("stock value for ")
	outputBuilder.WriteString(response.Symbol)
	outputBuilder.WriteString(" value ")
	outputBuilder.WriteString(fmt.Sprintf("%.2f", response.Value))
	outputBuilder.WriteString(" ")
	outputBuilder.WriteString(response.Currency)

	return outputBuilder.String()
}

func (t *tool) Auto() bool {
	return false
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	client := rest.NewWithOptions(
		config.Secret,
		rest.WithTrace(false),
		rest.WithPagination(true),
	)

	return &tool{
		client: client,
	}
}

func AddToMcp(server *mcp.Server) {
	tool_mcp.AddToMcp(server, Create(Configuration{
		Secret: environment.GetString("MASSIVE_TOKEN", ""),
	}))
}

func AddToContainer(container agent.ToolContainer, config Configuration) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config)))
}
