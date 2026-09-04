package brave_search_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hjwalt/platform/agent/util/brave_search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// roundTripFunc adapts a func to http.RoundTripper.
type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newClient(baseUrl string) *brave_search.BraveClient {
	return &brave_search.BraveClient{
		Client:  http.DefaultClient,
		BaseUrl: baseUrl,
	}
}

func TestWebSearchSendsParamsAndHeadersAndParses(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/web/search", r.URL.Path)

		query := r.URL.Query()
		assert.Equal("golang", query.Get("q"))
		assert.Equal("5", query.Get("count"))
		assert.Equal("us", query.Get("country"))
		assert.Equal("news,videos", query.Get("result_filter"))

		assert.Equal("tok-abc", r.Header.Get("X-Subscription-Token"))
		assert.Equal("no-cache", r.Header.Get("Cache-Control"))
		assert.Equal("test-agent/1.0", r.Header.Get("User-Agent"))
		assert.Equal("47.000", r.Header.Get("X-Loc-Lat"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"type": "search",
			"query": {"original": "golang", "country": "us"},
			"web": {
				"type": "search",
				"mutated_by_goggles": false,
				"results": [
					{"title": "The Go Programming Language", "url": "https://go.dev", "family_friendly": true}
				]
			}
		}`))
	}))
	defer server.Close()

	client := newClient(server.URL)
	params := []brave_search.Param{
		brave_search.WithTerm("golang"),
		brave_search.WithCount(5),
		brave_search.WithCountry("us"),
		brave_search.WithResultFilter([]string{"news", "videos"}),
	}
	headers := []brave_search.Header{
		brave_search.WithSubscriptionToken("tok-abc"),
		brave_search.WithNoCache(),
		brave_search.WithUserAgent("test-agent/1.0"),
		brave_search.WithLocLatitude(float32(47)),
	}

	result, err := brave_search.WebSearch(context.Background(), client, params, headers)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal("search", result.Type)
	require.NotNil(t, result.Query)
	assert.Equal("golang", result.Query.Original)
	require.NotNil(t, result.Web)
	require.Len(t, result.Web.Results, 1)
	assert.Equal("The Go Programming Language", result.Web.Results[0].Title)
	assert.Equal("https://go.dev", result.Web.Results[0].URL)
	assert.True(result.Web.Results[0].FamilyFriendly)
}

func TestImageSearchUsesImagesPathAndParses(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/images/search", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"type": "images",
			"mutated_by_goggles": false,
			"results": [{"type": "image", "title": "a sunset", "url": "https://img.example/sunset.jpg"}]
		}`))
	}))
	defer server.Close()

	result, err := brave_search.ImageSearch(context.Background(), newClient(server.URL),
		[]brave_search.Param{brave_search.WithTerm("sunset")}, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal("images", result.Type)
	require.Len(t, result.Results, 1)
	assert.Equal("a sunset", result.Results[0].Title)
	assert.Equal("https://img.example/sunset.jpg", result.Results[0].URL)
}

func TestVideoSearchUsesVideosPathAndParses(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/videos/search", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"type": "videos",
			"results": [{"type": "video_result", "title": "how to test"}]
		}`))
	}))
	defer server.Close()

	result, err := brave_search.VideoSearch(context.Background(), newClient(server.URL), nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Results, 1)
	assert.Equal("how to test", result.Results[0].Title)
}

func TestSpellcheckUsesSpellcheckPathAndParses(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/spellcheck/search", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"type": "spellcheck",
			"results": [{"query": "golang"}]
		}`))
	}))
	defer server.Close()

	result, err := brave_search.Spellcheck(context.Background(), newClient(server.URL),
		[]brave_search.Param{brave_search.WithSpellCheckLang("en")}, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Results, 1)
	assert.Equal("golang", result.Results[0].Query)
}

func TestSearchNon200ReturnsFailureResponseWithRawQuery(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{
			"error": {
				"id": "err-1",
				"status": 400,
				"code": "InvalidQuery",
				"detail": "query is not valid",
				"meta": {"component": "search", "errors": []}
			},
			"time": "2026-09-04T00:00:00Z"
		}`))
	}))
	defer server.Close()

	_, err := brave_search.WebSearch(context.Background(), newClient(server.URL),
		[]brave_search.Param{brave_search.WithTerm("bad query"), brave_search.WithCount(5)}, nil)

	require.Error(t, err)
	assert.Equal("query is not valid", err.Error())

	var failure brave_search.FailureResponse
	require.True(t, errors.As(err, &failure))
	assert.Equal("err-1", failure.ErrorResponse.ID)
	assert.Equal(400, failure.ErrorResponse.Status)
	assert.Equal("InvalidQuery", failure.ErrorResponse.Code)
	// the raw query string sent to the server is captured for diagnostics
	assert.Contains(failure.ErrorResponse.RawQuery, "q=bad+query")
	assert.Contains(failure.ErrorResponse.RawQuery, "count=5")
	assert.Equal("2026-09-04T00:00:00Z", failure.Time)
}

func TestSearchNon200WithNonJsonBodyReturnsUnmarshalError(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server exploded"))
	}))
	defer server.Close()

	_, err := brave_search.WebSearch(context.Background(), newClient(server.URL), nil, nil)

	require.Error(t, err)
	var failure brave_search.FailureResponse
	assert.False(errors.As(err, &failure))
}

func TestSearch200WithNonJsonBodyReturnsUnmarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html>not json</html>"))
	}))
	defer server.Close()

	_, err := brave_search.WebSearch(context.Background(), newClient(server.URL), nil, nil)

	require.Error(t, err)
}

func TestSearchClientDoErrorIsPropagated(t *testing.T) {
	assert := assert.New(t)

	sentinel := errors.New("connection refused")
	client := &brave_search.BraveClient{
		Client: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, sentinel
		})},
		BaseUrl: "http://brave.invalid",
	}

	_, err := brave_search.WebSearch(context.Background(), client, nil, nil)

	assert.ErrorIs(err, sentinel)
}

func TestSearchInvalidBaseUrlFailsToParse(t *testing.T) {
	assert := assert.New(t)

	client := newClient("://missing-scheme")

	_, err := brave_search.WebSearch(context.Background(), client, nil, nil)

	require.Error(t, err)
	assert.Contains(err.Error(), "missing protocol scheme")
}

func TestSearchAppendsPathToBaseUrlWithExistingPrefix(t *testing.T) {
	assert := assert.New(t)

	// note: Search concatenates u.Path + path with no joining slash, so the
	// base url prefix must end in "/" (e.g. https://host/res/v1/).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/res/v1/web/search", r.URL.Path)
		_, _ = w.Write([]byte(`{"type": "search"}`))
	}))
	defer server.Close()

	_, err := brave_search.WebSearch(context.Background(), newClient(server.URL+"/res/v1/"), nil, nil)

	assert.NoError(err)
}
