package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kultivator-consulting/goharmony"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	ctx := context.Background()
	client := openai.NewClient(
		option.WithBaseURL("http://localhost:13305/api/v1"),
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	question := "Write me a haiku about computers"

	resp, err := client.Completions.New(ctx, openai.CompletionNewParams{
		Prompt: openai.CompletionNewParamsPromptUnion{OfString: openai.String(question)},
		Model:  "gpt-oss-20b-FLM",
	})

	if err != nil {
		panic(err)
	}

	slog.Info("output", "resp", resp)

	parser := goharmony.NewParser()

	response := resp.Choices[0].Text
	// Extract only the user-facing message
	finalMessage := parser.ExtractFinalMessage(response)
	fmt.Println(finalMessage)
	// Output: Hello! How can I help you today?

	// Parse all messages
	messages, _ := parser.ParseResponse(response)
	for _, msg := range messages {
		fmt.Printf("Channel: %s, Content: %s\n", msg.Channel, msg.Content)
	}
}
