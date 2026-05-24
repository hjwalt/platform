package tool_container

import (
	"github.com/hjwalt/platform/agent"
	agent_tool "github.com/hjwalt/platform/agent/tool"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
)

func New() agent_tool.Container {
	return &container{
		sync: make(map[string]agent_tool.SyncWrapper),
	}
}

type container struct {
	sync map[string]agent_tool.SyncWrapper
}

func (r *container) AddSync(tool agent_tool.SyncWrapper) {
	r.sync[tool.Name()] = tool
}

func (r *container) GetSync() map[string]agent_tool.SyncWrapper {
	return r.sync
}

func (r *container) AsToolMap() map[string]agent.Tool {
	toolMap := make(map[string]agent.Tool)
	for k, v := range r.sync {
		toolMap[k] = tool_string_wrapper.WrapWrapper(v)
	}
	return toolMap
}
