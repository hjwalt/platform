package harness_container

import (
	"context"
	"errors"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/type/optional"
	"github.com/openai/openai-go/v3"
)

func NewToolContainer() agent.ToolContainer {
	return &toolContainer{
		tools: make([]string, 0),
		sync:  make(map[string]agent.SyncToolWrapper),
		async: make(map[string]agent.AsyncToolWrapper),
	}
}

type toolContainer struct {
	tools []string
	sync  map[string]agent.SyncToolWrapper
	async map[string]agent.AsyncToolWrapper
}

func (r *toolContainer) AddSync(tool agent.SyncToolWrapper) {
	r.tools = append(r.tools, tool.Name())
	r.sync[tool.Name()] = tool
}

func (r *toolContainer) AddAsync(tool agent.AsyncToolWrapper) {
	r.tools = append(r.tools, tool.Name())
	r.async[tool.Name()] = tool
}

func (r *toolContainer) OpenAiParamsFiltered(allowed []string) []openai.ChatCompletionToolUnionParam {
	if len(allowed) == 0 {
		allowed = append(allowed, r.tools...)
	}

	tools := make([]openai.ChatCompletionToolUnionParam, 0)
	for _, toolName := range allowed {
		if v, syncExists := r.sync[toolName]; syncExists {
			tools = append(tools, llm.OpenAiFromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
		} else if v, asyncExists := r.async[toolName]; asyncExists {
			tools = append(tools, llm.OpenAiFromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
		}
	}
	return tools
}

func (r *toolContainer) DeepSeekParams(allowed []string) []deepseek.Tool {
	if len(allowed) == 0 {
		allowed = append(allowed, r.tools...)
	}

	tools := make([]deepseek.Tool, 0)
	for _, toolName := range allowed {
		if v, syncExists := r.sync[toolName]; syncExists {
			tools = append(tools, llm.DeepSeekToolFromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
		} else if v, asyncExists := r.async[toolName]; asyncExists {
			tools = append(tools, llm.DeepSeekToolFromJsonSchema(v.Name(), v.Description(), v.RequestSchema()))
		}
	}
	return tools
}

func (r *toolContainer) Execute(ctx context.Context, in agent.Message, call agent.ToolCall) (optional.Optional[string], error) {
	if tool, exists := r.sync[call.Name]; exists {
		response, internalErr := tool.Apply(ctx, call.Arguments)
		if internalErr != nil {
			return optional.Empty[string](), internalErr
		}
		return optional.Of(tool.DescribeResult(response)), nil
	} else if asyncTool, exists := r.async[call.Name]; exists {
		internalErr := asyncTool.Send(ctx, agent.AgentContext{ParentContext: in.Context}, call, call.Arguments)
		if internalErr != nil {
			return optional.Empty[string](), internalErr
		}
		return optional.Empty[string](), nil
	} else {
		return optional.Empty[string](), errors.New("tool " + call.Name + " does not exist")
	}
}

func (r *toolContainer) DescribeRequest(call agent.ToolCall) (string, error) {
	if tool, exists := r.sync[call.Name]; exists {
		return tool.DescribeRequest(call.Arguments), nil
	} else if asyncTool, exists := r.async[call.Name]; exists {
		return asyncTool.DescribeRequest(call.Arguments), nil
	} else {
		return "", errors.New("tool " + call.Name + " does not exist")
	}
}

func (r *toolContainer) Exists(call agent.ToolCall) bool {
	if _, exists := r.sync[call.Name]; exists {
		return true
	} else if _, exists := r.async[call.Name]; exists {
		return true
	} else {
		return false
	}
}

func (r *toolContainer) Auto(call agent.ToolCall) bool {
	if tool, exists := r.sync[call.Name]; exists {
		return tool.Auto()
	} else if asyncTool, exists := r.async[call.Name]; exists {
		return asyncTool.Auto()
	} else {
		return false
	}
}
