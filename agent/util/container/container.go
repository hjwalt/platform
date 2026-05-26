package tool_container

import (
	"context"
	"errors"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/openai/openai-go/v3"
)

func New() agent.ToolContainer {
	return &container{
		sync: make(map[string]agent.SyncToolWrapper),
	}
}

type container struct {
	sync map[string]agent.SyncToolWrapper
}

func (r *container) AddSync(tool agent.SyncToolWrapper) {
	r.sync[tool.Name()] = tool
}

func (r *container) GetSync() map[string]agent.SyncToolWrapper {
	return r.sync
}

func (r *container) OpenAiParams() []openai.ChatCompletionToolUnionParam {
	tools := make([]openai.ChatCompletionToolUnionParam, 0)
	for _, v := range r.sync {
		tools = append(tools, llm.FromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
	}
	return tools
}

func (r *container) Execute(ctx context.Context, call agent.ToolCall) (string, error) {
	if tool, exists := r.sync[call.Name]; exists {
		response, internalErr := tool.Apply(ctx, call.Arguments)
		if internalErr != nil {
			return "", internalErr
		}
		return tool.DescribeResult(response), nil
	} else {
		return "", errors.New("tool " + call.Name + " does not exist")
	}
}

func (r *container) DescribeRequest(call agent.ToolCall) (string, error) {
	if tool, exists := r.sync[call.Name]; exists {
		return tool.DescribeRequest(call.Arguments), nil
	} else {
		return "", errors.New("tool " + call.Name + " does not exist")
	}
}

func (r *container) Exists(call agent.ToolCall) bool {
	_, exists := r.sync[call.Name]
	return exists
}

func (r *container) Auto(call agent.ToolCall) bool {
	if tool, exists := r.sync[call.Name]; exists {
		return tool.Auto()
	}
	return false
}
