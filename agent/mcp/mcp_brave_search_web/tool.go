package mcp_brave_search_web

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	brave "github.com/hjwalt/platform/agent/tool/brave_search"
	"github.com/hjwalt/platform/agent/tool/brave_search_web"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/format"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var global *Tool

func defaultTool() {
	global = &Tool{
		Brave: &brave.BraveClient{
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
	Brave  *brave.BraveClient
	ApiKey string
}

func (t *Tool) Behaviour(ctx context.Context, req *mcp.CallToolRequest, params *Request) (*mcp.CallToolResult, *Response, error) {
	success, err := brave_search_web.WebSearch(
		context.Background(),
		t.Brave,
		[]brave.Param{
			brave_search_web.WithTerm(params.Term),
		},
		[]brave.Header{
			brave_search_web.WithSubscriptionToken(t.ApiKey),
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
	success, err := brave_search_web.WebSearch(
		context.Background(),
		t.Brave,
		[]brave.Param{
			brave_search_web.WithTerm(params.Term),
		},
		[]brave.Header{
			brave_search_web.WithSubscriptionToken(t.ApiKey),
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

func (t *Tool) Execute(input string) (string, error) {
	request, requestParseErr := RequestFormat.Unmarshal([]byte(input))
	if requestParseErr != nil {
		return "", requestParseErr
	}

	response, internalErr := t.internal(&request)
	if internalErr != nil {
		return "", internalErr
	}

	output, marshalErr := ResponseFormat.Marshal(*response)
	if marshalErr != nil {
		return "", marshalErr
	}

	return string(output), nil
}

func Add(server *mcp.Server) {
	defaultTool()

	opts := &jsonschema.ForOptions{}

	in, _ := jsonschema.For[Request](opts)
	out, _ := jsonschema.For[Response](opts)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         "Search",
			Title:        "Search",
			Description:  "Search the internet with the terms to gather more information. Use the URL in the search results to fetch the page for more information.",
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
