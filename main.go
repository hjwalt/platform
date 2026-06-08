package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/agent/llm"
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	finance_stock_price_tool "github.com/hjwalt/platform/agent/tool/finance_stock_price"
	linux_shell_tool "github.com/hjwalt/platform/agent/tool/linux_shell"
	web_search_tool "github.com/hjwalt/platform/agent/tool/web_search"
	"github.com/hjwalt/platform/configuration"
	"github.com/hjwalt/platform/environment"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/hjwalt/platform/runtime"
	file_store "github.com/hjwalt/platform/state/file"
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
		Model: configuration.ModelConfiguration{
			Configurations: map[string]llm.ModelConfig{
				"openai": {
					Type:     llm.OpenAi,
					Model:    environment.GetString("OPENAI_API_MODEL", "gemma4-it-e4b-FLM"),
					Endpoint: environment.GetString("OPENAI_API_ENDPOINT", "http://localhost:13305/api/v1"),
					Secret:   environment.GetString("OPENAI_API_KEY", "lemonade"),
				},
				"lemonade": {
					Type:     llm.OpenAi,
					Model:    "gemma4-it-e4b-FLM",
					Endpoint: "http://localhost:13305/api/v1",
					Secret:   "lemonade",
				},
				"deepseek": {
					Type:     llm.DeepSeek,
					Model:    "deepseek-v4-flash",
					Endpoint: "https://api.deepseek.com",
					Secret:   environment.GetString("DEEPSEEK_TOKEN", "deepseek"),
				},
			},
			Parser: "lemonade",
			Agent:  "lemonade",
		},
		Tool: configuration.ToolConfiguration{
			WebSearch: web_search_tool.Configuration{
				BaseUrl: environment.GetString("BRAVE_SEARCH_URL", "https://api.search.brave.com/res/v1/"),
				Secret:  environment.GetString("BRAVE_TOKEN", ""),
			},
			Shell: linux_shell_tool.Configuration{
				BaseDir: "/home/hjwalt/Projects/platform/tmp/cmd",
			},
			FxPrice: finance_fx_price_tool.Configuration{
				Secret: environment.GetString("MASSIVE_TOKEN", ""),
			},
			StockPrice: finance_stock_price_tool.Configuration{
				Secret: environment.GetString("MASSIVE_TOKEN", ""),
			},
		},
		Server: configuration.WebServerConfiguration{
			Port:               3001,
			StaticResourcePath: "./web/static",
		},
		Flow: configuration.FlowConfiguration{
			Store: file_store.Configuration{
				Path: "/home/hjwalt/Projects/platform/tmp/agent/",
			},
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

	// Runtimes

	configuration.RegisterKafkaProducer(holder, config)
	configuration.RegisterKafkaAgentMessageProducer(holder, config)
	configuration.RegisterParserModel(holder, config)
	configuration.RegisterTools(holder, config)
	configuration.RegisterAgentHarnessStore(holder, config)
	configuration.RegisterAgentModel(holder, config)
	configuration.RegisterKafkaAgentFlow(holder, config)

	// HTTP

	httpBuilder := route.NewConfiguration[web.Context]()

	httpBuilder.AddMiddlewares(
		middleware.RequestID,
		middleware.CleanPath,
		middleware.Recoverer,
	)

	httpBuilder.AddDecorators(
		&decorator.RuntimeDecorator{
			Chat:                 holder.GetAgentModel(),
			AgentMessageProducer: holder.GetAgentMessageProducer(),
			AgentHarnessStore:    converter.RuntimeToFlowStore(holder.GetAgentHarnessStore(), format.Json[harness.ExecutionState]()),
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
