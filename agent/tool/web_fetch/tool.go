package web_fetch_tool

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	tool_mcp "github.com/hjwalt/platform/agent/util/mcp"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/format"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name = "web fetch"
)

type Configuration struct {
}

type Request struct {
	Link string `json:"link" jsonschema:"Link or URL to fetch"`
}

type Response struct {
	Html string `json:"html" jsonschema:"html output"`
}

type tool struct {
}

func (t *tool) Apply(ctx context.Context, params Request) (Response, error) {
	response, err := http.Get(params.Link)
	if err != nil {
		return Response{}, err
	}

	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return Response{}, err
	}

	return Response{
		Html: string(body),
	}, nil
}

func (t *tool) Name() string {
	return Name
}

func (t *tool) Description() string {
	return "Fetch content of link or URL"
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

	outputBuilder.WriteString("fetching URL `")
	outputBuilder.WriteString(request.Link)
	outputBuilder.WriteString("`")

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
	// TODO: parse this html to reduce the amount of token used
	return response.Html
}

func (t *tool) Auto() bool {
	return false
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	return &tool{}
}

func AddToMcp(server *mcp.Server) {
	tool_mcp.AddToMcp(server, Create(Configuration{}))
}

func AddToContainer(container agent.ToolContainer, config Configuration) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config)))
}
