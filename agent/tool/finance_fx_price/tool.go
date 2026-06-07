package finance_fx_price_tool

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
	Name = "finance_fx_price"
)

type Configuration struct {
	Secret string
}

type Request struct {
	Base  string `json:"base" jsonschema:"first currency listed, reference point for the quote"`
	Quote string `json:"quote" jsonschema:"second currency listed. indicates value required to equal one unit of base quote"`
}

type Response struct {
	Base  string  `json:"base" jsonschema:"first currency listed, reference point for the trade"`
	Quote string  `json:"quote" jsonschema:"second currency listed. indicates value required to equal one unit of base quote"`
	Value float32 `json:"ask" jsonschema:"conversion price"`
}

type tool struct {
	client *rest.Client
}

func (t *tool) Apply(ctx context.Context, params Request) (Response, error) {

	apiParams := &gen.GetForexSMAParams{
		Timespan: rest.Ptr(gen.GetForexSMAParamsTimespanMinute),
		Adjusted: rest.Ptr(true),
		Window:   rest.Ptr(50),
		Order:    rest.Ptr(gen.GetForexSMAParamsOrder("desc")),
		Limit:    rest.Ptr(10),
	}

	res, err := t.client.GetForexSMAWithResponse(
		ctx,
		"C:"+params.Base+params.Quote,
		apiParams,
	)
	if res.HTTPResponse.StatusCode != 200 {
		return Response{
			Base:  params.Base,
			Quote: params.Quote,
		}, errors.New("failed, http response status is " + res.HTTPResponse.Status)
	}

	slog.Info("fx quote", "request", params, "response", res.JSON200)

	resValues := res.JSON200.Results.Values
	if resValues == nil {
		return Response{
			Base:  params.Base,
			Quote: params.Quote,
		}, errors.New("failed, no result values")

	}

	return Response{
		Base:  params.Base,
		Quote: params.Quote,
		Value: *(*resValues)[0].Value,
	}, err
}

func (t *tool) Name() string {
	return Name
}

func (t *tool) Description() string {
	return "get FX or currency conversion average price"
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

	outputBuilder.WriteString("getting FX quote for ")
	outputBuilder.WriteString(request.Base)
	outputBuilder.WriteString(request.Quote)

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

	outputBuilder.WriteString("FX conversion value for ")
	outputBuilder.WriteString(response.Base)
	outputBuilder.WriteString(response.Quote)
	outputBuilder.WriteString(" value ")
	outputBuilder.WriteString(fmt.Sprintf("%.2f", response.Value))

	return outputBuilder.String()
}

func (t *tool) Auto() bool {
	return false
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	client := rest.NewWithOptions(
		config.Secret,
		rest.WithTrace(false),     // set true for full request/response logging
		rest.WithPagination(true), // enables automatic pagination via iterator
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
