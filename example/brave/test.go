package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/hjwalt/platform/agent/tool/brave_search"
	"github.com/hjwalt/platform/environment"
)

func main() {
	apiKey := environment.GetString("BRAVE_TOKEN", "")

	braveClient := brave_search.BraveClient{
		Client:  http.DefaultClient,
		BaseUrl: "https://api.search.brave.com/res/v1/",
	}

	success, err := brave_search.WebSearch(
		context.Background(),
		&braveClient,
		[]brave_search.Param{
			brave_search.WithTerm("the latest samsung phone"),
		},
		[]brave_search.Header{
			brave_search.WithSubscriptionToken(apiKey),
		},
	)
	if err != nil {
		slog.Error("error", "error", err)
	} else {
		slog.Info("success", "response", success)
		slog.Info("success", "response", success.Web)
	}
}
