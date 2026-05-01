package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/agent/mcp/mcp_brave_search_web"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/example/route/decorators"
	"github.com/hjwalt/platform/example/route/page_billing"
	"github.com/hjwalt/platform/example/route/page_chat"
	"github.com/hjwalt/platform/example/route/page_home"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_memory"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/flow/stateless"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/memory"
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
		// Model: "Gemma-4-26B-A4B-it-GGUF",
		Model:    "gpt-oss-20b-FLM",
		Endpoint: "http://localhost:13305/api/v1",
		Secret:   "nothing",
		Tools: []openai.ChatCompletionToolUnionParam{
			llm.OpenAiToolSchema[mcp_brave_search_web.Request]("web_search", "search the internet with the terms for the more information"),
		},
	})

	store := rag.Memory()

	ragModel := rag.Rag(model, store)

	agentFlow := harness.OpenAiFlow[context.Context]{
		Tools: map[string]agent.Tool{
			"web_search": mcp_brave_search_web.Instance(),
		},
		Model: ragModel,
	}

	// In Memory Messaging

	agentMessageChannel := memory.MemoryConfiguration{
		Channel: make(chan message.Message[memory.MemoryMetadata], 100),
	}

	messageProducer := converter.RuntimeToFlowProducer(
		memory.NewProducer(agentMessageChannel),
		converter.NewConverter(
			flow_runtime_memory.New(),
			format.Json[agent.Message](),
		),
	)

	chatConsumer := memory.NewConsumer(
		agentMessageChannel,
		converter.FlowToRuntimeHandler(
			stateless.NewExploder(
				"Increment",
				agentFlow.Handle,
				metadata.MessageUpdate(),
				messageProducer,
				messageProducer,
			),
			converter.NewConverter(
				flow_runtime_memory.New(),
				format.Json[agent.Message](),
			),
		),
	)

	// HTTP

	httpBuilder := route.NewConfiguration[example.Context]()

	httpBuilder.AddMiddlewares(
		middleware.RequestID,
		middleware.RealIP,
		middleware.CleanPath,
		middleware.Recoverer,
	)

	httpBuilder.AddDecorators(
		&decorators.RuntimeDecorator{
			Chat:                 ragModel,
			RagStore:             store,
			AgentMessageProducer: messageProducer,
		},
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
			store,
			ragModel,
			messageProducer,
			chatConsumer,
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
