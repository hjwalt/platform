package memory_tool

import (
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/hjwalt/platform/agent"
	memory_clear_tool "github.com/hjwalt/platform/agent/tool/memory_clear"
	memory_get_tool "github.com/hjwalt/platform/agent/tool/memory_get"
	memory_update_tool "github.com/hjwalt/platform/agent/tool/memory_update"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
)

const (
	MemoryFileName   = "memory.md"
	MemoryGetName    = "memory_get"
	MemoryUpdateName = "memory_update"
	MemoryClearName  = "memory_clear"
)

var validPrefixPattern = regexp.MustCompile(`^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`)

type Configuration struct {
	BaseDir string
	Prefix  string
}

type Tools struct {
	Get    agent.SyncTool[GetRequest, GetResponse]
	Update agent.SyncTool[UpdateRequest, UpdateResponse]
	Clear  agent.SyncTool[ClearRequest, ClearResponse]
}

type UpdateMode = memory_update_tool.UpdateMode

const (
	UpdateModeReplace = memory_update_tool.UpdateModeReplace
	UpdateModeAppend  = memory_update_tool.UpdateModeAppend
)

type GetRequest = memory_get_tool.Request
type GetResponse = memory_get_tool.Response
type UpdateRequest = memory_update_tool.Request
type UpdateResponse = memory_update_tool.Response
type ClearRequest = memory_clear_tool.Request
type ClearResponse = memory_clear_tool.Response

func Create(config Configuration) (Tools, error) {
	baseDir := strings.TrimSpace(config.BaseDir)
	prefix := strings.TrimSpace(config.Prefix)
	if baseDir == "" {
		return Tools{}, ErrInvalidBaseDir
	}
	if prefix != "" && !validPrefixPattern.MatchString(prefix) {
		return Tools{}, ErrInvalidPrefix
	}

	sharedMutex := &sync.Mutex{}

	return Tools{
		Get: memory_get_tool.Create(memory_get_tool.Configuration{
			BaseDir:  baseDir,
			FileName: memoryFileName(prefix),
			Name:     toolName(prefix, MemoryGetName),
			Mutex:    sharedMutex,
		}),
		Update: memory_update_tool.Create(memory_update_tool.Configuration{
			BaseDir:  baseDir,
			FileName: memoryFileName(prefix),
			Name:     toolName(prefix, MemoryUpdateName),
			Mutex:    sharedMutex,
		}),
		Clear: memory_clear_tool.Create(memory_clear_tool.Configuration{
			BaseDir:  baseDir,
			FileName: memoryFileName(prefix),
			Name:     toolName(prefix, MemoryClearName),
			Mutex:    sharedMutex,
		}),
	}, nil
}

func AddToContainer(container agent.ToolContainer, config Configuration) error {
	tools, err := Create(config)
	if err != nil {
		return err
	}

	container.AddSync(tool_string_wrapper.StringWrapSync(tools.Get))
	container.AddSync(tool_string_wrapper.StringWrapSync(tools.Update))
	container.AddSync(tool_string_wrapper.StringWrapSync(tools.Clear))

	return nil
}

func toolName(prefix string, base string) string {
	if prefix == "" {
		return base
	}
	return prefix + "_" + base
}

func memoryFileName(prefix string) string {
	if prefix == "" {
		return MemoryFileName
	}
	return prefix + ".md"
}

var (
	ErrInvalidBaseDir    = errors.New("memory root path cannot be empty")
	ErrInvalidPrefix     = errors.New("memory tool prefix must match ^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$")
	ErrInvalidUpdateMode = memory_update_tool.ErrInvalidUpdateMode
)
