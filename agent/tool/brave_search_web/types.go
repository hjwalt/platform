package brave_search_web

import (
	brave "github.com/hjwalt/platform/agent/tool/brave_search"
	"github.com/hjwalt/platform/format"
)

var (
	FailureResponseFormat = format.Json[FailureResponse]()
	SuccessResponseFormat = format.Json[SuccessResponse]()
)

type FailureResponse struct {
	ErrorResponse brave.ErrorResponse `json:"error"`
	Time          string              `json:"time"`
}

func (er FailureResponse) Error() string {
	return er.ErrorResponse.Detail
}

type SuccessResponse struct {
	Type        string                                         `json:"type"`
	Discussions *brave.ResultContainer[brave.DiscussionResult] `json:"discussions"`
	FAQ         *brave.ResultContainer[brave.QA]               `json:"faq"`
	InfoBox     *brave.ResultContainer[brave.GraphInfoBox]     `json:"infobox"`
	Locations   *brave.ResultContainer[brave.LocationResult]   `json:"locations"`
	Mixed       *brave.Mixed                                   `json:"mixed"`
	News        *brave.ResultContainer[brave.NewsResult]       `json:"news"`
	Query       *brave.Query                                   `json:"query"`
	Videos      *brave.ResultContainer[brave.VideoResult]      `json:"videos"`
	Web         *brave.ResultContainer[brave.SearchResult]     `json:"web"`
	Summarizer  *brave.Summarizer                              `json:"summarizer"`
}
