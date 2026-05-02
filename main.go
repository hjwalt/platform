package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/mcp/mcp_brave_search_web"
	"github.com/hjwalt/platform/configuration"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/example/route/decorators"
	"github.com/hjwalt/platform/example/route/page_billing"
	"github.com/hjwalt/platform/example/route/page_chat"
	"github.com/hjwalt/platform/example/route/page_home"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/memory"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/web/route"
	"github.com/joho/godotenv"
)

func main() {
	logger.Default()
	godotenv.Load()

	// Config building

	config := configuration.Configuration{
		OpenAi: configuration.OpenAiConfiguration{
			Model:    environment.GetString("OPENAI_API_MODEL", "gemma4-it-e4b-FLM"),
			Endpoint: environment.GetString("OPENAI_API_ENDPOINT", "http://localhost:13305/api/v1"),
			Secret:   environment.GetString("OPENAI_API_KEY", "lemonade"),
		},
		BraveSearch: configuration.BraveSearchConfiguration{
			BaseUrl: environment.GetString("BRAVE_SEARCH_URL", "https://api.search.brave.com/res/v1/"),
			Secret:  environment.GetString("BRAVE_SEARCH_TOKEN", ""),
		},
		Server: configuration.WebServerConfiguration{
			Port:               3001,
			StaticResourcePath: "./web/static",
		},
	}

	holder := configuration.Holder()

	// Flow

	agentMessageChannel := memory.MemoryConfiguration{
		Channel: make(chan message.Message[memory.MemoryMetadata], 100),
	}

	// AI -- tools

	tools := []agent.Tool{
		mcp_brave_search_web.Instance(),
	}

	// Runtimes

	model := configuration.RegisterOpenAi(holder, config, tools)
	store := configuration.RegisterInMemoryRagMemory(holder, config)
	ragModel := configuration.RegisterRagModel(holder, config, model, store)
	messageProducer := configuration.RegisterInMemoryAgentMessageProducer(holder, config, agentMessageChannel)
	configuration.RegisterInMemoryAgentMessageConsumer(holder, config, agentMessageChannel, messageProducer, ragModel, tools)

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

	httpBuilder.HandleStatic("/static/", config.Server.StaticResourcePath)
	page_home.Add(httpBuilder)
	page_billing.Add(httpBuilder)
	page_chat.Add(httpBuilder)

	httpHandler := httpBuilder.Build()

	httpRuntime := runtime.NewHttp(
		runtime.HttpWithPort(config.Server.Port),
		runtime.HttpWithHandler(httpHandler),
	)

	// start

	holder.Add(httpRuntime)
	holder.Block()

}
