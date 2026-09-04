package flow_runtime_memory_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/flow_runtime_memory"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/message/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRuntimeFlowToRuntimeMetadata(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_memory.New()

	source := flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 9,
		Source:   "src-1",
	}

	result := runtime.FlowToRuntimeMetadata(source)

	expected, err := metadata.Format.Marshal(source)
	assert.NoError(err)
	assert.Equal(1, len(result.Headers))
	assert.Contains(result.Headers, metadata.MetadataHeaderKey)
	assert.Equal(string(expected), result.Headers[metadata.MetadataHeaderKey])
}

func TestMemoryRuntimeRuntimeToFlowMetadata(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_memory.New()

	source := flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 9,
		Source:   "src-1",
	}

	content, err := metadata.Format.Marshal(source)
	assert.NoError(err)

	result := runtime.RuntimeToFlowMetadata(memory.MemoryMetadata{
		Headers: map[string]string{metadata.MetadataHeaderKey: string(content)},
	})

	assert.Equal(source, result)
}

func TestMemoryRuntimeMetadataRoundTrip(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_memory.New()

	source := flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 9,
		Source:   "src-1",
	}

	back := runtime.RuntimeToFlowMetadata(runtime.FlowToRuntimeMetadata(source))

	assert.Equal(source, back)
}

func TestMemoryRuntimeMetadataMissingHeaderReturnsDefault(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_memory.New()

	result := runtime.RuntimeToFlowMetadata(memory.MemoryMetadata{
		Headers: map[string]string{},
	})

	assertDefaultMetadata(assert, result)
}

func TestMemoryRuntimeMetadataGarbageHeaderReturnsDefault(t *testing.T) {
	assert := assert.New(t)

	// silence the slog.Error the implementation logs on the garbage path
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	runtime := flow_runtime_memory.New()

	result := runtime.RuntimeToFlowMetadata(memory.MemoryMetadata{
		Headers: map[string]string{metadata.MetadataHeaderKey: "{not valid json"},
	})

	assertDefaultMetadata(assert, result)
}

func assertDefaultMetadata(assert *assert.Assertions, m flow.Metadata) {
	assert.Equal(int32(0), m.Attempt)
	assert.Equal(int64(-1), m.Sequence)
	assert.Equal("UNKNOWN", m.Source)
	assert.NotEmpty(m.Id)
}
