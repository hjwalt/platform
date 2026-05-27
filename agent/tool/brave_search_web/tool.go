package brave_search_web_tool

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/util/brave_search"
	tool_mcp "github.com/hjwalt/platform/agent/util/mcp"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/reflect"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Configuration struct {
	BaseUrl string
	Secret  string
}

type Request struct {
	Term string `json:"term" jsonschema:"search query term"`
}

type Response struct {
	Results []SearchResult `json:"results" jsonschema:"search results"`
}

type SearchResult struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Description   string   `json:"description"`
	Language      string   `json:"language"`
	ContentType   string   `json:"content_type"`
	ExtraSnippets []string `json:"extra_snippets"`
}

type Tool struct {
	Brave  *brave_search.BraveClient
	ApiKey string
}

func (t *Tool) Apply(ctx context.Context, params Request) (Response, error) {
	success, err := brave_search.WebSearch(
		context.Background(),
		t.Brave,
		[]brave_search.Param{
			brave_search.WithTerm(params.Term),
		},
		[]brave_search.Header{
			brave_search.WithSubscriptionToken(t.ApiKey),
		},
	)

	if err != nil {
		return Response{Results: make([]SearchResult, 0)}, err
	}

	results := make([]SearchResult, len(success.Web.Results))
	for i, braveResults := range success.Web.Results {
		results[i] = SearchResult{
			Title:         braveResults.Title,
			URL:           braveResults.URL,
			Description:   braveResults.Description,
			Language:      braveResults.Language,
			ContentType:   braveResults.ContentType,
			ExtraSnippets: braveResults.ExtraSnippets,
		}
	}

	slog.Info("search", "request", params, "response", len(results))

	return Response{Results: results}, err
}

func (t *Tool) Name() string {
	return "web search"
}

func (t *Tool) Description() string {
	return "Search the internet with the terms to gather more information. Use the URL in the search results to fetch the page for more information."
}

func (t *Tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *Tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *Tool) DescribeRequest(request Request) string {
	outputBuilder := strings.Builder{}

	outputBuilder.WriteString("search the web with term ")
	outputBuilder.WriteString(request.Term)

	return outputBuilder.String()
}

func (t *Tool) ResultFormat() format.Format[Response] {
	return format.Json[Response]()
}

func (t *Tool) ResultSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Response](opts)
	return toolSchema
}

func (t *Tool) DescribeResult(response Response) string {
	outputBuilder := strings.Builder{}

	for i, result := range response.Results {
		if i > 0 {
			outputBuilder.WriteString("---------------------")
			outputBuilder.WriteString("\n\n")
		}
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("result ")
		outputBuilder.WriteString(reflect.GetString(i + 1))
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("title: ")
		outputBuilder.WriteString(result.Title)
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("url: ")
		outputBuilder.WriteString(result.URL)
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("description: ")
		outputBuilder.WriteString(result.Description)
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("language: ")
		outputBuilder.WriteString(result.Language)
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("content_type: ")
		outputBuilder.WriteString(result.ContentType)
		outputBuilder.WriteString("\n\n")
		outputBuilder.WriteString("extra_snippets: ")
		outputBuilder.WriteString("\n\n")
		for _, snippet := range result.ExtraSnippets {
			outputBuilder.WriteString("- ")
			outputBuilder.WriteString(snippet)
			outputBuilder.WriteString("\n\n")
		}
	}

	return outputBuilder.String()
}

func (t *Tool) Auto() bool {
	return true
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	return &Tool{
		Brave: &brave_search.BraveClient{
			Client:  http.DefaultClient,
			BaseUrl: config.BaseUrl,
		},
		ApiKey: config.Secret,
	}
}

func AddToMcp(server *mcp.Server) {
	tool_mcp.AddToMcp(server, Create(Configuration{
		BaseUrl: "https://api.search.brave.com/res/v1/",
		Secret:  environment.GetString("BRAVE_TOKEN", ""),
	}))
}

func AddToContainer(container agent.ToolContainer, config Configuration) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config)))
}
