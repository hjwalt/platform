package llm

import (
	"context"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func createOpenAi(config ModelConfig, tools agent.ToolContainer) agent.LanguageModel {
	return &openAiModel{
		Model:    config.Model,
		Endpoint: config.Endpoint,
		Secret:   config.Secret,
		Tools:    tools,
	}
}

type openAiModel struct {
	Model    string
	Endpoint string
	Secret   string
	Tools    agent.ToolContainer
	client   openai.Client
}

func (r *openAiModel) Start() error {
	r.client = openai.NewClient(
		option.WithBaseURL(r.Endpoint),
		option.WithAPIKey(r.Secret), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return nil
}

func (r *openAiModel) Stop() {
}

func (r *openAiModel) Chat(ctx context.Context, messages []agent.Message, allowedTools []string) ([]agent.Message, error) {
	modelMessage := make([]openai.ChatCompletionMessageParamUnion, 0)
	for _, message := range messages {
		switch message.Type {
		case agent.MessageType_System:
			{
				modelMessage = append(modelMessage, OpenAiSystemMessage(message))
			}
		case agent.MessageType_Assistant:
			{
				modelMessage = append(modelMessage, OpenAiAssistantMessage(message))
			}
		case agent.MessageType_User:
			{
				modelMessage = append(modelMessage, OpenAiUserMessage(message))
			}
		case agent.MessageType_ToolRequest:
			{
				modelMessage = append(modelMessage, OpenAiToolRequestMessage(message))
			}
		case agent.MessageType_ToolResult:
			{
				modelMessage = append(modelMessage, OpenAiToolResultMessage(message))
			}
		}
	}

	params := openai.ChatCompletionNewParams{
		Messages: modelMessage,
		Model:    r.Model,
		Tools:    r.Tools.OpenAiParamsFiltered(allowedTools),
	}

	completion, err := r.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return []agent.Message{agent.NewMessage(
			messages[0].Context,
			agent.MessageType_Error,
			err.Error(),
			"",
			agent.ToolCall{},
		)}, err
	}

	outputMessages := make([]agent.Message, 0)
	for _, choice := range completion.Choices {
		switch choice.FinishReason {
		case "stop":
			{
				outputMessages = append(outputMessages, agent.NewMessage(
					messages[0].Context,
					agent.MessageType_Agent,
					choice.Message.Content,
					"",
					agent.ToolCall{},
				))
			}
		case "tool_calls":
			{
				for _, toolCall := range choice.Message.ToolCalls {
					toolData := agent.ToolCall{
						Id:        toolCall.ID,
						Name:      toolCall.Function.Name,
						Arguments: toolCall.Function.Arguments,
					}

					if toolRequestMessage, messageErr := r.Tools.DescribeRequest(toolData); messageErr == nil {
						outputMessages = append(outputMessages, agent.NewMessage(
							messages[0].Context,
							agent.MessageType_ToolRequest,
							toolRequestMessage,
							"",
							toolData,
						))
					} else {
						// TODO: do something with error
					}
				}
			}
		}
	}

	return outputMessages, nil
}

func OpenAiToolSchema[M any](name string, description string) openai.ChatCompletionToolUnionParam {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[M](opts)
	return OpenAiFromJsonSchema(name, description, toolSchema)
}

func OpenAiFromJsonSchema(name string, description string, toolSchema *jsonschema.Schema) openai.ChatCompletionToolUnionParam {
	unmarshalled, _ := format.Convert(toolSchema, schemaFormat, openAiFormat)
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        name,
		Description: openai.String(description),
		Parameters:  unmarshalled,
	})
}

func OpenAiSystemMessage(message agent.Message) openai.ChatCompletionMessageParamUnion {
	return openai.SystemMessage(message.Message)
}

func OpenAiAssistantMessage(message agent.Message) openai.ChatCompletionMessageParamUnion {
	return openai.AssistantMessage(message.Message)
}

func OpenAiUserMessage(message agent.Message) openai.ChatCompletionMessageParamUnion {
	return openai.UserMessage(message.Message)
}

func OpenAiToolRequestMessage(message agent.Message) openai.ChatCompletionMessageParamUnion {
	return openai.ChatCompletionMessageParamUnion{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			ToolCalls: []openai.ChatCompletionMessageToolCallUnionParam{
				{
					OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
						ID: message.Tool.Id,
						Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
							Arguments: message.Tool.Arguments,
							Name:      message.Tool.Name,
						},
					},
				},
			},
		},
	}
}

func OpenAiToolResultMessage(message agent.Message) openai.ChatCompletionMessageParamUnion {
	return openai.ToolMessage(message.Message, message.Tool.Id)
}

var (
	OpenAiMessageFormat = format.Json[openai.ChatCompletionMessage]()
	schemaFormat        = format.Json[*jsonschema.Schema]()
	openAiFormat        = format.Json[openai.FunctionParameters]()
)
