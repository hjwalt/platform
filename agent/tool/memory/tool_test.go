package memory_tool

import (
	"context"
	"testing"

	"github.com/hjwalt/platform/agent"
	tool_container "github.com/hjwalt/platform/agent/util/container"
	memory_store "github.com/hjwalt/platform/state/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryToolNameByPrefix(t *testing.T) {
	assert := assert.New(t)

	store := memory_store.New()
	defaultTool := Create(Configuration{}, store)
	prefixedTool := Create(Configuration{Key: "corrections"}, store)

	assert.Equal("memory", defaultTool.Name())
	assert.Equal("corrections_memory", prefixedTool.Name())
}

func TestMemoryToolGetUpdateClearFlow(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := memory_store.New()
	memory := Create(Configuration{Key: "preferences"}, store)

	getEmpty, getEmptyErr := memory.Apply(ctx, Request{Operation: OperationGet})
	assert.NoError(getEmptyErr)
	assert.Equal(OperationGet, getEmpty.Operation)
	assert.Equal("preferences", getEmpty.Key)
	assert.False(getEmpty.Exists)
	assert.Equal("", getEmpty.Content)

	updateReplace, updateReplaceErr := memory.Apply(ctx, Request{
		Operation: OperationUpdate,
		Content:   "A",
		Mode:      UpdateModeReplace,
	})
	assert.NoError(updateReplaceErr)
	assert.Equal(OperationUpdate, updateReplace.Operation)
	assert.Equal(UpdateModeReplace, updateReplace.Mode)
	assert.Equal(1, updateReplace.Bytes)

	updateAppend, updateAppendErr := memory.Apply(ctx, Request{
		Operation: OperationUpdate,
		Content:   "B",
		Mode:      UpdateModeAppend,
	})
	assert.NoError(updateAppendErr)
	assert.Equal(OperationUpdate, updateAppend.Operation)
	assert.Equal(UpdateModeAppend, updateAppend.Mode)
	assert.Equal(2, updateAppend.Bytes)

	getFilled, getFilledErr := memory.Apply(ctx, Request{Operation: OperationGet})
	assert.NoError(getFilledErr)
	assert.True(getFilled.Exists)
	assert.Equal("AB", getFilled.Content)

	clearResponse, clearErr := memory.Apply(ctx, Request{Operation: OperationClear})
	assert.NoError(clearErr)
	assert.Equal(OperationClear, clearResponse.Operation)
	assert.True(clearResponse.Cleared)

	getAfterClear, getAfterClearErr := memory.Apply(ctx, Request{Operation: OperationGet})
	assert.NoError(getAfterClearErr)
	assert.False(getAfterClear.Exists)
	assert.Equal("", getAfterClear.Content)
}

func TestMemoryToolValidationAndErrors(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	err := Validate(Configuration{Key: "bad-prefix!"}, memory_store.New())
	assert.ErrorIs(err, ErrInvalidPrefix)

	err = Validate(Configuration{}, nil)
	assert.ErrorIs(err, ErrNilStore)

	memory := Create(Configuration{}, memory_store.New())

	_, applyErr := memory.Apply(ctx, Request{Operation: Operation("invalid")})
	assert.ErrorIs(applyErr, ErrInvalidOperation)

	_, applyErr = memory.Apply(ctx, Request{
		Operation: OperationUpdate,
		Content:   "x",
		Mode:      UpdateMode("invalid"),
	})
	assert.ErrorIs(applyErr, ErrInvalidMode)
}

func TestAddToContainerRegistersSingleMemoryTool(t *testing.T) {
	assert := assert.New(t)
	container := tool_container.New()

	AddToContainer(container, Configuration{Key: "improvements"}, memory_store.New())

	assert.True(container.Exists(agent.ToolCall{Name: "improvements_memory"}))
	assert.False(container.Exists(agent.ToolCall{Name: "improvements_memory_get"}))
	assert.False(container.Exists(agent.ToolCall{Name: "improvements_memory_update"}))
	assert.False(container.Exists(agent.ToolCall{Name: "improvements_memory_clear"}))
}
