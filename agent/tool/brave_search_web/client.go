package brave_search_web

import (
	"context"
	"io"
	"net/http"
	"net/url"

	brave "github.com/hjwalt/platform/agent/tool/brave_search"
)

const (
	webSearchPath = "web/search"
)

// PARAMS

func WebSearch(ctx context.Context, b *brave.BraveClient, params []brave.Param, headers []brave.Header) (*SuccessResponse, error) {
	u, err := url.Parse(b.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = u.Path + webSearchPath

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
		resp, unmarshalErr := SuccessResponseFormat.Unmarshal(resBody)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}

		return &resp, nil
	}
}
