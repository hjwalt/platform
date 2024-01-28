package brave

import (
	"encoding/json"
	"net/http"
)

type BraveClient struct {
	Client  *http.Client
	BaseUrl string
}

type Header struct {
	Key   string
	Value string
}

type Param struct {
	Key   string
	Value string
}

func handleRequest[T any](client *http.Client, req *http.Request) (*T, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var resp errorResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			return nil, err
		}

		resp.Error.Time = resp.Time
		resp.Error.RawQuery = req.URL.RawQuery
		return nil, resp.Error
	}

	var resp T
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type errorResponse struct {
	Error ErrorResponse `json:"error"`
	Time  string        `json:"time"`
}
