package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/format"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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

	// marshalled, _ := schemaFormat.Marshal(toolSchema)
	// unmarshalled, _ := openAiFormat.Unmarshal(marshalled)

	slog.Info("schemajson", "schema", unmarshalled)

	question := "what is the weather in Jakarta right now"
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("Get weather at the given location"),
				Parameters:  unmarshalled,
			}),
		},
		Seed:  openai.Int(0),
		Model: "Gemma-4-26B-A4B-it-GGUF",
		// Model: "gpt-oss-20b-FLM",
	}

	ctx := context.Background()
	client := openai.NewClient(
		option.WithBaseURL("http://localhost:13305/api/v1"),
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	// Make initial chat completion request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	slog.Info("completion", "completion", completion)

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}

	// If there is a was a function call, continue the conversation
	params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
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
			fmt.Printf("Weather in %s: %s\n", location, weatherData)

			params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
		}
	}

	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}

func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}
