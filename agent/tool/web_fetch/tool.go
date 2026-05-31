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
	Name = "web_fetch"
)

type Configuration struct {
}

type Request struct {
	Link string `json:"link" jsonschema:"Link or URL to fetch"`
}

type Response struct {
	Html   string `json:"html" jsonschema:"html output"`
	Parsed string `json:"parsed" jsonschema:"parsed content using language model"`
}

type tool struct {
	model agent.LanguageModel
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

	htmlResult := string(body)

	result, err := t.model.Chat(
		context.Background(),
		[]agent.Message{
			{
				Type:    agent.MessageType_User,
				Message: "parse the following html into its content: \n\n" + htmlResult,
			},
		},
		[]string{},
	)

	parsed := ""
	if err != nil {
		parsed = "failed to parse response HTML using model due to error " + err.Error()
	}

	if len(result) == 0 {
		parsed = "failed to parse response HTML due to model returning empty results"
	}

	if result[0].Type != agent.MessageType_Agent {
		parsed = "failed to parse response HTML due to model returning invalid response"
	}

	parsed = result[0].Message

	return Response{
		Html:   htmlResult,
		Parsed: parsed,
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
	return response.Parsed
}

func (t *tool) Auto() bool {
	return false
}

func Create(config Configuration, model agent.LanguageModel) agent.SyncTool[Request, Response] {
	return &tool{
		model: model,
	}
}

func AddToMcp(server *mcp.Server, model agent.LanguageModel) {
	tool_mcp.AddToMcp(server, Create(Configuration{}, model))
}

func AddToContainer(container agent.ToolContainer, config Configuration, model agent.LanguageModel) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config, model)))
}
