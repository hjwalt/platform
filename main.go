package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/agent/mcp/mcp_brave_search_web"
	shell_tool "github.com/hjwalt/platform/agent/tool/shell"
	"github.com/hjwalt/platform/configuration"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/decorator"
	"github.com/hjwalt/platform/web/page/page_billing"
	"github.com/hjwalt/platform/web/page/page_chat"
	"github.com/hjwalt/platform/web/page/page_home"
	"github.com/hjwalt/platform/web/route"
	"github.com/joho/godotenv"
)

func main() {
	logger.Default()
	godotenv.Load()

	// Config building

	instanceId := uuid.New().String()

	config := configuration.Configuration{
		OpenAi: configuration.OpenAiConfiguration{
			Model:    environment.GetString("OPENAI_API_MODEL", "gemma4-it-e4b-FLM"),
			Endpoint: environment.GetString("OPENAI_API_ENDPOINT", "http://localhost:13305/api/v1"),
			Secret:   environment.GetString("OPENAI_API_KEY", "lemonade"),
		},
		Tool: configuration.ToolConfiguration{
			BraveSearch: mcp_brave_search_web.Configuration{
				BaseUrl: environment.GetString("BRAVE_SEARCH_URL", "https://api.search.brave.com/res/v1/"),
				Secret:  environment.GetString("BRAVE_TOKEN", ""),
			},
			Shell: shell_tool.Configuration{
				BaseDir: "/home/hjwalt/Projects/platform/tmp/cmd",
			},
		},
		Server: configuration.WebServerConfiguration{
			Port:               3001,
			StaticResourcePath: "./web/static",
		},
		Flow: configuration.FlowConfiguration{
			Agent: configuration.AgentFlowConfiguration{
				Topic: "AGENT",
				Producer: kafka.KafkaProducerConfiguration{
					Brokers:  "localhost:9092",
					ClientId: "agent-producer-" + instanceId,
				},
				Consumer: kafka.KafkaConsumerConfiguration{
					Brokers:  "localhost:9092",
					ClientId: "agent-consumer-" + instanceId,
					Topic:    "AGENT",
					GroupId:  "agent-consumer",
				},
			},
			Result: configuration.AgentFlowConfiguration{
				Topic: "AGENT-RESULT",
				Producer: kafka.KafkaProducerConfiguration{
					Brokers:  "localhost:9092",
					ClientId: "result-producer-" + instanceId,
				},
				Consumer: kafka.KafkaConsumerConfiguration{
					Brokers:  "localhost:9092",
					ClientId: "result-consumer-" + instanceId,
					Topic:    "AGENT-RESULT",
					GroupId:  "result-consumer",
				},
			},
		},
	}

	holder := configuration.ContextBuilder()

	// AI -- tools

	holder.AddTool(mcp_brave_search_web.Instance(config.Tool.BraveSearch))
	holder.AddTool(shell_tool.Instance(config.Tool.Shell))

	// Runtimes

	configuration.RegisterAgentHarnessStore(holder, config)
	configuration.RegisterOpenAi(holder, config)
	configuration.RegisterKafkaAgentFlow(holder, config)

	// HTTP

	httpBuilder := route.NewConfiguration[web.Context]()

	httpBuilder.AddMiddlewares(
		middleware.RequestID,
		middleware.RealIP,
		middleware.CleanPath,
		middleware.Recoverer,
	)

	httpBuilder.AddDecorators(
		&decorator.RuntimeDecorator{
			Chat:                 holder.GetLanguageModel(),
			AgentMessageProducer: holder.GetAgentMessageProducer(),
			AgentHarnessStore:    converter.RuntimeToFlowStore(holder.GetAgentHarnessStore(), format.Json[harness.ExecutionState]()),
			Tool:                 holder.GetTool(),
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
