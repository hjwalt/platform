package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/format"
	"github.com/openai/openai-go/v3"
)

type ToolRequest struct {
	Location string `json:"location" jsonschema:"location to get weather from"`
}

func main() {

	opts := &jsonschema.ForOptions{}

	toolSchema, _ := jsonschema.For[ToolRequest](opts)

	schemaFormat := format.Json[*jsonschema.Schema]()
	openAiFormat := format.Json[openai.FunctionParameters]()

	unmarshalled, _ := format.Convert(toolSchema, schemaFormat, openAiFormat)

	tools := []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_weather",
			Description: openai.String("Get weather at the given location"),
			Parameters:  unmarshalled,
		}),
	}

	model := llm.OpenAi(llm.OpenAiModelConfig{
		Model:    "Gemma-4-26B-A4B-it-GGUF",
		Endpoint: "http://localhost:13305/api/v1",
		Secret:   "nothing",
		Tools:    tools,
	})

	model.Start()

	messages := []agent.Message{
		{
			Type:    agent.MessageType_User,
			Message: "what is the weather in Jakarta right now",
		},
	}

	executing := true

	for executing {
		executing = false
		result, err := model.Chat(context.Background(), messages)
		if err != nil {
			panic(err)
		}

		messages = append(messages, result...)
		for _, resultMessage := range result {
			switch resultMessage.Type {
			case agent.MessageType_ToolRequest:

				executing = true
				rawMessage, unmarshallErr := llm.OpenAiMessageFormat.Unmarshal([]byte(resultMessage.Raw))
				if unmarshallErr != nil {
					panic(unmarshallErr)
				}
				toolCalls := rawMessage.ToolCalls
				for _, toolCall := range toolCalls {
					if toolCall.Function.Name == "get_weather" {
						// Extract the location from the function call arguments
						var args map[string]interface{}
						err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						if err != nil {
							panic(err)
						}
						location := args["location"].(string)

						// Simulate getting weather data
						weatherData := getWeather(location)

						// Print the weather data
						slog.Info("Weather in", location, weatherData)

						toolDataRaw, toolMarshallErr := agent.ToolDataFormat.Marshal(agent.ToolData{
							Id: toolCall.ID,
						})
						if toolMarshallErr != nil {
							panic(toolMarshallErr)
						}

						messages = append(messages, agent.Message{
							Type:    agent.MessageType_ToolResult,
							Message: weatherData,
							Raw:     string(toolDataRaw),
						})
					}
				}
			case agent.MessageType_Agent:
				{
					slog.Info(resultMessage.Message)
				}
			case agent.MessageType_Error:
				{
					slog.Error("error", "error", resultMessage)
				}
			}
		}
	}

	// marshalled, _ := schemaFormat.Marshal(toolSchema)
	// unmarshalled, _ := openAiFormat.Unmarshal(marshalled)

	slog.Info("schemajson", "schema", unmarshalled)
}

func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}
