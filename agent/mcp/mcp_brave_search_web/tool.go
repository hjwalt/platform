package mcp_brave_search_web

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/agent/tool/brave_search"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/reflect"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3"
)

var global *Tool

func defaultTool() {
	global = &Tool{
		Brave: &brave_search.BraveClient{
			Client:  http.DefaultClient,
			BaseUrl: "https://api.search.brave.com/res/v1/",
		},
		ApiKey: environment.GetString("BRAVE_TOKEN", ""),
	}
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

func (t *Tool) Behaviour(ctx context.Context, req *mcp.CallToolRequest, params *Request) (*mcp.CallToolResult, *Response, error) {
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
		return nil, nil, err
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

	slog.Info("search", "request", req, "response", len(results))

	return nil, &Response{Results: results}, err
}

func (t *Tool) internal(params *Request) (*Response, error) {
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
		return nil, err
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

	return &Response{Results: results}, err
}

func (t *Tool) Name() string {
	return "web search"
}

func (t *Tool) Description() string {
	return "Search the internet with the terms to gather more information. Use the URL in the search results to fetch the page for more information."
}

func (t *Tool) Schema() openai.ChatCompletionToolUnionParam {
	return llm.OpenAiToolSchema[Request](t.Name(), t.Description())
}

func (t *Tool) Execute(input string) (string, error) {
	request, requestParseErr := RequestFormat.Unmarshal([]byte(input))
	if requestParseErr != nil {
		return "", requestParseErr
	}

	response, internalErr := t.internal(&request)
	if internalErr != nil {
		return "", internalErr
	}

	outputBuilder := strings.Builder{}

	for i, result := range response.Results {
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
		outputBuilder.WriteString("---------------------")
		outputBuilder.WriteString("\n\n")
	}

	return outputBuilder.String(), nil
}

func Add(server *mcp.Server) {
	defaultTool()

	opts := &jsonschema.ForOptions{}

	in, _ := jsonschema.For[Request](opts)
	out, _ := jsonschema.For[Response](opts)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         global.Name(),
			Title:        global.Name(),
			Description:  global.Description(),
			InputSchema:  in,
			OutputSchema: out,
		},
		global.Behaviour,
	)
}

func Instance() agent.Tool {
	defaultTool()

	return global
}

var (
	RequestFormat  = format.Json[Request]()
	ResponseFormat = format.Json[Response]()
)
