package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/agent/rag"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_memory"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/flow/stateless"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/memory"
	"github.com/hjwalt/platform/runtime"
	"github.com/openai/openai-go/v3"
)

type ToolRequest struct {
	Location string `json:"location" jsonschema:"location to get weather from"`
}

func main() {
	agentMessageChannel := memory.MemoryConfiguration{
		Channel: make(chan message.Message[memory.MemoryMetadata], 100),
	}

	messageFormat := format.Json[agent.Message]()

	initialMessage, _ := messageFormat.Marshal(agent.Message{
		Type:    agent.MessageType_User,
		Message: "what is the weather in Jakarta right now",
		Raw:     "",
	})

	agentMessageChannel.Channel <- message.Message[memory.MemoryMetadata]{
		Value: initialMessage,
	}

	weatherTool := WeatherTool{}

	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model: "Gemma-4-26B-A4B-it-GGUF",
		// Model:    "Gemma-4-E4B-it-GGUF",
		Endpoint: "http://localhost:13305/api/v1",
		Secret:   "nothing",
		Tools: []openai.ChatCompletionToolUnionParam{
			weatherTool.Schema(),
		},
	})

	store := rag.Memory()

	ragModel := rag.Rag(model, store)

	agentFlow := harness.OpenAiFlow[context.Context]{
		Tools: map[string]agent.Tool{
			weatherTool.Name(): weatherTool,
		},
		Model: ragModel,
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

	runtime.Start(
		[]runtime.Runtime{
			model,
			store,
			ragModel,
			messageProducer,
			chatConsumer,
		},
		100*time.Millisecond,
	)

	runtime.Wait()
}

type WeatherTool struct {
}

func (t WeatherTool) Name() string {
	return "get_weather"
}

func (t WeatherTool) Description() string {
	return "Get weather at the given location"
}

func (t WeatherTool) Schema() openai.ChatCompletionToolUnionParam {
	return llm.OpenAiToolSchema[ToolRequest](t.Name(), t.Description())
}

func (t WeatherTool) Execute(input string) (string, error) {
	var args map[string]interface{}
	err := json.Unmarshal([]byte(input), &args)
	if err != nil {
		panic(err)
	}
	location := args["location"].(string)

	// Print the weather data
	slog.Info("getting weather in", "location", location)

	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C", nil
}
