package main

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/example/route/decorators"
	"github.com/hjwalt/platform/example/route/page_billing"
	"github.com/hjwalt/platform/example/route/page_chat"
	"github.com/hjwalt/platform/example/route/page_home"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/web/route"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
)

func main() {
	logger.Default()
	godotenv.Load()

	// Chat

	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    "Gemma-4-26B-A4B-it-GGUF",
		Endpoint: "http://localhost:13305/api/v1",
		Secret:   "nothing",
		Tools:    make([]openai.ChatCompletionToolUnionParam, 0),
	})

	// HTTP

	httpBuilder := route.NewConfiguration[example.Context]()

	httpBuilder.AddMiddlewares(
		middleware.RequestID,
		middleware.RealIP,
		middleware.CleanPath,
		middleware.Recoverer,
	)

	httpBuilder.AddDecorators(
		&decorators.RuntimeDecorator{Chat: model},
	)

	httpBuilder.HandleStatic("/static/", "./example/route/static")
	page_home.Add(httpBuilder)
	page_billing.Add(httpBuilder)
	page_chat.Add(httpBuilder)

	httpHandler := httpBuilder.Build()

	httpRuntime := runtime.NewHttp(
		runtime.HttpWithPort(3001),
		runtime.HttpWithHandler(httpHandler),
	)

	// start

	startErr := runtime.Start(
		[]runtime.Runtime{
			model,
			httpRuntime,
		},
		time.Second,
	)

	if startErr != nil {
		panic(startErr)
	}

	slog.Info("started")

	runtime.Wait()

	slog.Info("stopped")
}
