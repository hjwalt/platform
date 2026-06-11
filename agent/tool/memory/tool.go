package memory_tool

import (
	"errors"
	"regexp"
	"sync"

	"github.com/hjwalt/platform/agent"
	memory_clear_tool "github.com/hjwalt/platform/agent/tool/memory_clear"
	memory_get_tool "github.com/hjwalt/platform/agent/tool/memory_get"
	memory_update_tool "github.com/hjwalt/platform/agent/tool/memory_update"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/state"
)

const (
	MemoryKey        = "memory"
	MemoryGetName    = "memory_get"
	MemoryUpdateName = "memory_update"
	MemoryClearName  = "memory_clear"
)

var validPrefixPattern = regexp.MustCompile(`^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`)

type Configuration struct {
	Key string
}

type Tools struct {
	Get    agent.SyncTool[memory_get_tool.Request, memory_get_tool.Response]
	Update agent.SyncTool[memory_update_tool.Request, memory_update_tool.Response]
	Clear  agent.SyncTool[memory_clear_tool.Request, memory_clear_tool.Response]
}

type UpdateMode = memory_update_tool.UpdateMode

const (
	UpdateModeReplace = memory_update_tool.UpdateModeReplace
	UpdateModeAppend  = memory_update_tool.UpdateModeAppend
)

func Create(config Configuration, store state.Store) Tools {

	sharedMutex := &sync.Mutex{}
	key := memoryKey(config.Key)

	return Tools{
		Get: memory_get_tool.Create(memory_get_tool.Configuration{
			Store: store,
			Key:   key,
			Name:  toolName(config.Key, MemoryGetName),
			Mutex: sharedMutex,
		}),
		Update: memory_update_tool.Create(memory_update_tool.Configuration{
			Store: store,
			Key:   key,
			Name:  toolName(config.Key, MemoryUpdateName),
			Mutex: sharedMutex,
		}),
		Clear: memory_clear_tool.Create(memory_clear_tool.Configuration{
			Store: store,
			Key:   key,
			Name:  toolName(config.Key, MemoryClearName),
			Mutex: sharedMutex,
		}),
	}
}

func AddToContainer(container agent.ToolContainer, config Configuration, store state.Store) {
	tools := Create(config, store)

	container.AddSync(tool_string_wrapper.StringWrapSync(tools.Get))
	container.AddSync(tool_string_wrapper.StringWrapSync(tools.Update))
	container.AddSync(tool_string_wrapper.StringWrapSync(tools.Clear))
}

func toolName(prefix string, base string) string {
	if prefix == "" {
		return base
	}
	return prefix + "_" + base
}

func memoryKey(prefix string) string {
	if prefix == "" {
		return MemoryKey
	}
	return prefix
}

var (
	ErrNilStore          = errors.New("memory store cannot be nil")
	ErrInvalidPrefix     = errors.New("memory tool prefix must match ^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$")
	ErrInvalidUpdateMode = memory_update_tool.ErrInvalidUpdateMode
)
