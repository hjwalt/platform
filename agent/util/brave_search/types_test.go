package brave_search_test

import (
	"testing"

	"github.com/hjwalt/platform/agent/util/brave_search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailureResponseErrorReturnsDetail(t *testing.T) {
	assert := assert.New(t)

	resp := brave_search.FailureResponse{
		ErrorResponse: brave_search.ErrorResponse{
			ID:     "id-1",
			Status: 429,
			Code:   "RateLimited",
			Detail: "you have been rate limited",
			Time:   "2026-09-04T00:00:00Z",
		},
	}

	assert.Equal("you have been rate limited", resp.Error())
	assert.Equal("you have been rate limited", resp.ErrorResponse.Error())
}

func TestFailureResponseSatisfiesErrorInterface(t *testing.T) {
	var err error = brave_search.FailureResponse{ErrorResponse: brave_search.ErrorResponse{Detail: "d"}}
	require.NotNil(t, err)
	assert.Equal(t, "d", err.Error())
}

func TestFailureResponseJsonRoundTrip(t *testing.T) {
	assert := assert.New(t)

	original := brave_search.FailureResponse{
		ErrorResponse: brave_search.ErrorResponse{
			ID:       "id-9",
			Status:   400,
			Code:     "InvalidQuery",
			Detail:   "bad query",
			Meta:     brave_search.ErrorMeta{Component: "search"},
			RawQuery: "q=test",
			Time:     "internal-only", // json:"-" on the nested struct
		},
		Time: "2026-09-04T01:02:03Z",
	}

	body, err := brave_search.FailureResponseFormat.Marshal(original)
	assert.NoError(err)

	back, err := brave_search.FailureResponseFormat.Unmarshal(body)
	assert.NoError(err)
	assert.Equal(original.ErrorResponse.ID, back.ErrorResponse.ID)
	assert.Equal(original.ErrorResponse.Status, back.ErrorResponse.Status)
	assert.Equal(original.ErrorResponse.Code, back.ErrorResponse.Code)
	assert.Equal(original.ErrorResponse.Detail, back.ErrorResponse.Detail)
	assert.Equal(original.ErrorResponse.Meta, back.ErrorResponse.Meta)
	assert.Equal("q=test", back.ErrorResponse.RawQuery)
	// the nested ErrorResponse.Time is json:"-" and is never serialised
	assert.Empty(back.ErrorResponse.Time)
}
