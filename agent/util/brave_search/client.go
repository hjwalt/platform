package brave_search

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/hjwalt/platform/format"
)

const (
	webSearchPath   = "web/search"
	imageSearchPath = "images/search"
	spellcheckPath  = "spellcheck/search"
	videoSearchPath = "videos/search"
)

func WebSearch(ctx context.Context, b *BraveClient, params []Param, headers []Header) (*WebSearchResult, error) {
	return Search[WebSearchResult](ctx, b, params, headers, webSearchPath)
}

func ImageSearch(ctx context.Context, b *BraveClient, params []Param, headers []Header) (*ImageSearchResult, error) {
	return Search[ImageSearchResult](ctx, b, params, headers, imageSearchPath)
}

func VideoSearch(ctx context.Context, b *BraveClient, params []Param, headers []Header) (*VideoSearchResult, error) {
	return Search[VideoSearchResult](ctx, b, params, headers, videoSearchPath)
}

func Spellcheck(ctx context.Context, b *BraveClient, params []Param, headers []Header) (*SpellcheckResult, error) {
	return Search[SpellcheckResult](ctx, b, params, headers, spellcheckPath)
}

// actual calls

func Search[C any](ctx context.Context, b *BraveClient, params []Param, headers []Header, path string) (*C, error) {
	u, err := url.Parse(b.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = u.Path + path

	// add query parameters
	values := url.Values{}
	for _, param := range params {
		values.Add(param.Key, param.Value)
	}
	u.RawQuery = values.Encode()

	// build request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// add headers
	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}

	// execute
	res, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// read response body
	resBody, bodyReadErr := io.ReadAll(res.Body)
	if bodyReadErr != nil {
		return nil, bodyReadErr
	}

	// parse
	if res.StatusCode != http.StatusOK {
		resp, unmarshalErr := FailureResponseFormat.Unmarshal(resBody)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}

		resp.ErrorResponse.RawQuery = req.URL.RawQuery
		return nil, resp
	} else {
		resp, unmarshalErr := format.Json[C]().Unmarshal(resBody)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}

		return &resp, nil
	}
}
