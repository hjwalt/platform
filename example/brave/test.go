package main

import (
	"context"
	"log/slog"
	"net/http"

	brave "github.com/hjwalt/platform/agent/tool/brave_search"
	"github.com/hjwalt/platform/agent/tool/brave_search_web"
	"github.com/hjwalt/platform/environment"
)

func main() {
	apiKey := environment.GetString("BRAVE_TOKEN", "")

	braveClient := brave.BraveClient{
		Client:  http.DefaultClient,
		BaseUrl: "https://api.search.brave.com/res/v1/",
	}

	success, err := brave_search_web.WebSearch(
		context.Background(),
		&braveClient,
		[]brave.Param{
			brave_search_web.WithTerm("the latest samsung phone"),
		},
		[]brave.Header{
			brave_search_web.WithSubscriptionToken(apiKey),
		},
	)
	if err != nil {
		slog.Error("error", "error", err)
	} else {
		slog.Info("success", "response", success)
		slog.Info("success", "response", success.Web)
	}
}
