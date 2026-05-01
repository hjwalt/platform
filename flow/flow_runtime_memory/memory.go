package flow_runtime_memory

import (
	"log/slog"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/message/memory"
)

func New() flow.MessageRuntime[memory.MemoryMetadata] {
	return &MemoryRuntime{}
}

type MemoryRuntime struct {
}

func (r *MemoryRuntime) FlowToRuntimeMetadata(source flow.Metadata) memory.MemoryMetadata {
	headers := make(map[string]string)

	marshalled, marshalErr := metadata.Format.Marshal(source)
	if marshalErr != nil {
		slog.Error("failed to marshal metadata, resetting", "error", marshalErr)
	} else {
		headers[metadata.MetadataHeaderKey] = string(marshalled)
	}

	return memory.MemoryMetadata{
		Headers: headers,
	}
}

func (r *MemoryRuntime) RuntimeToFlowMetadata(source memory.MemoryMetadata) flow.Metadata {
	metadataContent, exists := source.Headers[metadata.MetadataHeaderKey]
	if !exists {
		return metadata.Default()
	}

	unmarshalled, unmarshalErr := metadata.Format.Unmarshal([]byte(metadataContent))
	if unmarshalErr != nil {
		slog.Error("failed to unmarshal metadata, resetting", "error", unmarshalErr)
		return metadata.Default()
	}

	return unmarshalled
}
