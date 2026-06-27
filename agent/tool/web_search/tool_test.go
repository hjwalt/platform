package web_search_tool

import (
	"testing"

	"github.com/hjwalt/platform/agent"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

func TestWebSearchName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	assert.Equal("web_search", tool.Name())
}

func TestWebSearchDescription(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	assert.NotEmpty(tool.Description())
}

func TestWebSearchAuto(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	assert.True(tool.Auto())
}

func TestWebSearchRequestSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	schema := tool.RequestSchema()
	assert.NotNil(schema)
}

func TestWebSearchResultSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	schema := tool.ResultSchema()
	assert.NotNil(schema)
}

func TestWebSearchRequestFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	assert.NotNil(tool.RequestFormat())
}

func TestWebSearchResultFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	assert.NotNil(tool.ResultFormat())
}

func TestWebSearchDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	desc := tool.DescribeRequest(Request{Term: "golang testing"})

	assert.Contains(desc, "golang testing")
	assert.Contains(desc, "search")
}

func TestWebSearchDescribeResultSingleResult(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	desc := tool.DescribeResult(Response{
		Results: []SearchResult{
			{
				Title:       "Go Testing Guide",
				URL:         "https://example.com/go-testing",
				Description: "A comprehensive guide to testing in Go",
				Language:    "en",
				ContentType: "article",
				ExtraSnippets: []string{
					"Table-driven tests",
					"Benchmark tests",
				},
			},
		},
	})

	assert.Contains(desc, "Go Testing Guide")
	assert.Contains(desc, "https://example.com/go-testing")
	assert.Contains(desc, "A comprehensive guide to testing in Go")
	assert.Contains(desc, "en")
	assert.Contains(desc, "article")
	assert.Contains(desc, "Table-driven tests")
	assert.Contains(desc, "Benchmark tests")
	assert.Contains(desc, "result 1")
}

func TestWebSearchDescribeResultMultipleResults(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	desc := tool.DescribeResult(Response{
		Results: []SearchResult{
			{Title: "Result One", URL: "https://one.example.com", Description: "First result", Language: "en", ContentType: "text"},
			{Title: "Result Two", URL: "https://two.example.com", Description: "Second result", Language: "fr", ContentType: "text"},
		},
	})

	assert.Contains(desc, "Result One")
	assert.Contains(desc, "Result Two")
	assert.Contains(desc, "result 1")
	assert.Contains(desc, "result 2")
	// Separator should appear between results
	assert.Contains(desc, "-----")
}

func TestWebSearchDescribeResultEmpty(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	desc := tool.DescribeResult(Response{
		Results: []SearchResult{},
	})

	// Should not panic and return empty string for no results
	assert.Empty(desc)
}

func TestWebSearchDescribeResultWithEmptySnippets(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	desc := tool.DescribeResult(Response{
		Results: []SearchResult{
			{
				Title:         "Simple Result",
				URL:           "https://simple.example.com",
				Description:   "A simple result",
				Language:      "en",
				ContentType:   "page",
				ExtraSnippets: nil,
			},
		},
	})

	assert.Contains(desc, "Simple Result")
	// Should not include "extra_snippets:" label entries when there are no snippets
	assert.NotContains(desc, "- ")
}

func TestWebSearchCreateReturnsSyncTool(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	var _ agent.SyncTool[Request, Response] = tool
	assert.NotNil(tool)
}

func TestAddToContainerRegistersWebSearchTool(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	AddToContainer(container, Configuration{
		BaseUrl: "https://api.example.com/",
		Secret:  "test-secret",
	})

	assert.True(container.Exists(agent.ToolCall{Name: "web_search"}))
}
