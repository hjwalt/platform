package tool_container

import (
	"context"
	"errors"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/type/optional"
	"github.com/openai/openai-go/v3"
)

func New() agent.ToolContainer {
	return &container{
		sync:  make(map[string]agent.SyncToolWrapper),
		async: make(map[string]agent.AsyncToolWrapper),
	}
}

type container struct {
	sync  map[string]agent.SyncToolWrapper
	async map[string]agent.AsyncToolWrapper
}

func (r *container) AddSync(tool agent.SyncToolWrapper) {
	r.sync[tool.Name()] = tool
}

func (r *container) AddAsync(tool agent.AsyncToolWrapper) {
	r.async[tool.Name()] = tool
}

func (r *container) OpenAiParams() []openai.ChatCompletionToolUnionParam {
	tools := make([]openai.ChatCompletionToolUnionParam, 0)
	for _, v := range r.sync {
		tools = append(tools, llm.FromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
	}
	for _, v := range r.async {
		tools = append(tools, llm.FromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
	}
	return tools
}

func (r *container) Execute(ctx context.Context, in agent.Message, call agent.ToolCall) (optional.Optional[string], error) {
	if tool, exists := r.sync[call.Name]; exists {
		response, internalErr := tool.Apply(ctx, call.Arguments)
		if internalErr != nil {
			return optional.Empty[string](), internalErr
		}
		return optional.Of(tool.DescribeResult(response)), nil
	} else if asyncTool, exists := r.async[call.Name]; exists {
		internalErr := asyncTool.Send(ctx, agent.Parent{Context: in.Context}, call, call.Arguments)
		if internalErr != nil {
			return optional.Empty[string](), internalErr
		}
		return optional.Empty[string](), nil
	} else {
		return optional.Empty[string](), errors.New("tool " + call.Name + " does not exist")
	}
}

func (r *container) DescribeRequest(call agent.ToolCall) (string, error) {
	if tool, exists := r.sync[call.Name]; exists {
		return tool.DescribeRequest(call.Arguments), nil
	} else if asyncTool, exists := r.async[call.Name]; exists {
		return asyncTool.DescribeRequest(call.Arguments), nil
	} else {
		return "", errors.New("tool " + call.Name + " does not exist")
	}
}

func (r *container) Exists(call agent.ToolCall) bool {
	if _, exists := r.sync[call.Name]; exists {
		return true
	} else if _, exists := r.async[call.Name]; exists {
		return true
	} else {
		return false
	}
}

func (r *container) Auto(call agent.ToolCall) bool {
	if tool, exists := r.sync[call.Name]; exists {
		return tool.Auto()
	} else if asyncTool, exists := r.async[call.Name]; exists {
		return asyncTool.Auto()
	} else {
		return false
	}
}
