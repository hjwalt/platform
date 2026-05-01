package llm

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func OpenAi(config OpenAiModelConfig) agent.LanguageModel {
	return &OpenAiModel{
		Model:    config.Model,
		Endpoint: config.Endpoint,
		Secret:   config.Secret,
		Tools:    config.Tools,
	}
}

type OpenAiModelConfig struct {
	Model    string
	Endpoint string
	Secret   string
	Tools    []openai.ChatCompletionToolUnionParam
}

type OpenAiModel struct {
	Model    string
	Endpoint string
	Secret   string
	Tools    []openai.ChatCompletionToolUnionParam
	client   openai.Client
}

func (r *OpenAiModel) Start() error {
	r.client = openai.NewClient(
		option.WithBaseURL(r.Endpoint),
		option.WithAPIKey(r.Secret), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return nil
}

func (r *OpenAiModel) Stop() {
}

func (r *OpenAiModel) Chat(ctx context.Context, messages []agent.Message) ([]agent.Message, error) {
	modelMessage := make([]openai.ChatCompletionMessageParamUnion, 0)
	for _, message := range messages {
		switch message.Type {
		case agent.MessageType_User:
			{
				modelMessage = append(modelMessage, openai.UserMessage(message.Message))
			}
		// having tool request in the history chain seems to throw the model in a pickle generating broken tool requests
		// case agent.MessageType_ToolRequest, agent.MessageType_Agent:
		// 	{
		// 		rawMessage, unmarshallErr := OpenAiMessageFormat.Unmarshal([]byte(message.Raw))
		// 		if unmarshallErr != nil {
		// 			return []agent.Message{}, unmarshallErr
		// 		}
		// 		modelMessage = append(modelMessage, rawMessage.ToParam())
		// 	}
		case agent.MessageType_ToolResult:
			{
				toolData, unmarshallErr := agent.ToolDataFormat.Unmarshal([]byte(message.Raw))
				if unmarshallErr != nil {
					return []agent.Message{}, unmarshallErr
				}
				modelMessage = append(modelMessage, openai.ToolMessage(message.Message, toolData.Id))
			}
		}
	}

	params := openai.ChatCompletionNewParams{
		Messages: modelMessage,
		Seed:     openai.Int(0),
		Model:    r.Model,
		Tools:    r.Tools,
	}

	completion, err := r.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return []agent.Message{{
			Context: messages[0].Context,
			Type:    agent.MessageType_Error,
			Message: err.Error(),
			Raw:     "",
		}}, err
	}

	outputMessages := make([]agent.Message, 0)
	for _, choice := range completion.Choices {
		raw, _ := OpenAiMessageFormat.Marshal(choice.Message)

		switch choice.FinishReason {
		case "stop":
			{
				outputMessages = append(outputMessages, agent.Message{
					Context: messages[0].Context,
					Type:    agent.MessageType_Agent,
					Message: choice.Message.Content,
					Raw:     string(raw),
				})
			}
		case "tool_calls":
			{
				outputMessages = append(outputMessages, agent.Message{
					Context: messages[0].Context,
					Type:    agent.MessageType_ToolRequest,
					Message: choice.Message.Content,
					Raw:     string(raw),
				})
			}
		}
	}

	return outputMessages, nil
}

func OpenAiToolSchema[M any](name string, description string) openai.ChatCompletionToolUnionParam {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[M](opts)
	unmarshalled, _ := format.Convert(toolSchema, schemaFormat, openAiFormat)
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        name,
		Description: openai.String(description),
		Parameters:  unmarshalled,
	})
}

var (
	OpenAiMessageFormat = format.Json[openai.ChatCompletionMessage]()
	schemaFormat        = format.Json[*jsonschema.Schema]()
	openAiFormat        = format.Json[openai.FunctionParameters]()
)
