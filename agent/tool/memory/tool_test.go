package memory_tool

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/hjwalt/platform/agent"
	tool_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

func TestCreateValidation(t *testing.T) {
	assert := assert.New(t)

	_, rootErr := Create(Configuration{RootPath: "", Prefix: "session"})
	assert.ErrorIs(rootErr, ErrInvalidRootPath)

	_, prefixErr := Create(Configuration{RootPath: "/tmp", Prefix: "bad-prefix"})
	assert.ErrorIs(prefixErr, ErrInvalidPrefix)
}

func TestCreateNames(t *testing.T) {
	assert := assert.New(t)

	withoutPrefix, withoutPrefixErr := Create(Configuration{RootPath: "/tmp"})
	assert.NoError(withoutPrefixErr)
	assert.Equal(MemoryGetName, withoutPrefix.Get.Name())
	assert.Equal(MemoryUpdateName, withoutPrefix.Update.Name())
	assert.Equal(MemoryClearName, withoutPrefix.Clear.Name())

	withPrefix, withPrefixErr := Create(Configuration{RootPath: "/tmp", Prefix: "session"})
	assert.NoError(withPrefixErr)
	assert.Equal("session_"+MemoryGetName, withPrefix.Get.Name())
	assert.Equal("session_"+MemoryUpdateName, withPrefix.Update.Name())
	assert.Equal("session_"+MemoryClearName, withPrefix.Clear.Name())
}

func TestGetMissingFile(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	root := t.TempDir()

	tools, createErr := Create(Configuration{RootPath: root})
	assert.NoError(createErr)

	result, getErr := tools.Get.Apply(ctx, GetRequest{})

	assert.NoError(getErr)
	assert.False(result.Exists)
	assert.Equal("", result.Content)
	assert.Equal(filepath.Join(root, MemoryFileName), result.Path)
}

func TestUpdateReplaceAndAppend(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	root := t.TempDir()

	tools, createErr := Create(Configuration{RootPath: root})
	assert.NoError(createErr)

	replaceResponse, replaceErr := tools.Update.Apply(ctx, UpdateRequest{
		Content: "line-one",
		Mode:    UpdateModeReplace,
	})
	assert.NoError(replaceErr)
	assert.Equal(UpdateModeReplace, replaceResponse.Mode)
	assert.Equal(len([]byte("line-one")), replaceResponse.Bytes)

	appendResponse, appendErr := tools.Update.Apply(ctx, UpdateRequest{
		Content: "\nline-two",
		Mode:    UpdateModeAppend,
	})
	assert.NoError(appendErr)
	assert.Equal(UpdateModeAppend, appendResponse.Mode)
	assert.Equal(len([]byte("line-one\nline-two")), appendResponse.Bytes)

	getResponse, getErr := tools.Get.Apply(ctx, GetRequest{})
	assert.NoError(getErr)
	assert.True(getResponse.Exists)
	assert.Equal("line-one\nline-two", getResponse.Content)
}

func TestUpdateDefaultModeReplace(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	root := t.TempDir()

	tools, createErr := Create(Configuration{RootPath: root})
	assert.NoError(createErr)

	_, firstErr := tools.Update.Apply(ctx, UpdateRequest{Content: "A"})
	assert.NoError(firstErr)

	_, secondErr := tools.Update.Apply(ctx, UpdateRequest{Content: "B"})
	assert.NoError(secondErr)

	getResponse, getErr := tools.Get.Apply(ctx, GetRequest{})
	assert.NoError(getErr)
	assert.Equal("B", getResponse.Content)
}

func TestUpdateInvalidMode(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	root := t.TempDir()

	tools, createErr := Create(Configuration{RootPath: root})
	assert.NoError(createErr)

	_, updateErr := tools.Update.Apply(ctx, UpdateRequest{
		Content: "text",
		Mode:    "merge",
	})

	assert.ErrorIs(updateErr, ErrInvalidUpdateMode)
}

func TestClear(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	root := t.TempDir()

	tools, createErr := Create(Configuration{RootPath: root})
	assert.NoError(createErr)

	_, updateErr := tools.Update.Apply(ctx, UpdateRequest{Content: "content"})
	assert.NoError(updateErr)

	clearResponse, clearErr := tools.Clear.Apply(ctx, ClearRequest{})
	assert.NoError(clearErr)
	assert.True(clearResponse.Cleared)
	assert.Equal(filepath.Join(root, MemoryFileName), clearResponse.Path)

	getResponse, getErr := tools.Get.Apply(ctx, GetRequest{})
	assert.NoError(getErr)
	assert.True(getResponse.Exists)
	assert.Equal("", getResponse.Content)
}

func TestAddToContainerWithMultiplePrefixes(t *testing.T) {
	assert := assert.New(t)
	container := tool_container.New()

	sessionAddErr := AddToContainer(container, Configuration{RootPath: t.TempDir(), Prefix: "session"})
	repoAddErr := AddToContainer(container, Configuration{RootPath: t.TempDir(), Prefix: "repo"})
	baseAddErr := AddToContainer(container, Configuration{RootPath: t.TempDir()})

	assert.NoError(sessionAddErr)
	assert.NoError(repoAddErr)
	assert.NoError(baseAddErr)

	assert.True(container.Exists(agent.ToolCall{Name: "session_memory_get"}))
	assert.True(container.Exists(agent.ToolCall{Name: "session_memory_update"}))
	assert.True(container.Exists(agent.ToolCall{Name: "session_memory_clear"}))
	assert.True(container.Exists(agent.ToolCall{Name: "repo_memory_get"}))
	assert.True(container.Exists(agent.ToolCall{Name: "repo_memory_update"}))
	assert.True(container.Exists(agent.ToolCall{Name: "repo_memory_clear"}))
	assert.True(container.Exists(agent.ToolCall{Name: "memory_get"}))
	assert.True(container.Exists(agent.ToolCall{Name: "memory_update"}))
	assert.True(container.Exists(agent.ToolCall{Name: "memory_clear"}))
}
