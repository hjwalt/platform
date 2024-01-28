package flow_runtime_kafka

import (
	"log/slog"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/message/kafka"
)

func New(topic string) flow.MessageRuntime[kafka.KafkaMetadata] {
	return &KafkaRuntime{
		Topic: topic,
	}
}

type KafkaRuntime struct {
	Topic string
}

func (r *KafkaRuntime) FlowToRuntimeMetadata(source flow.Metadata) kafka.KafkaMetadata {
	headers := make(map[string]string)

	marshalled, marshalErr := metadata.Format.Marshal(source)
	if marshalErr != nil {
		slog.Error("failed to marshal metadata, resetting", "error", marshalErr)
	} else {
		headers[metadata.MetadataHeaderKey] = string(marshalled)
	}

	return kafka.KafkaMetadata{
		Topic:   r.Topic,
		Key:     source.Group,
		Headers: headers,
	}
}

func (r *KafkaRuntime) RuntimeToFlowMetadata(source kafka.KafkaMetadata) flow.Metadata {
	metadataContent, exists := source.Headers[metadata.MetadataHeaderKey]
	if !exists {
		return metadata.Default()
	}

	unmarshalled, unmarshalErr := metadata.Format.Unmarshal([]byte(metadataContent))
	if unmarshalErr != nil {
		slog.Error("failed to unmarshal metadata, resetting", "error", unmarshalErr)
		return metadata.Default()
	}

	return flow.Metadata{
		Id:       unmarshalled.Id,
		Group:    source.Key,
		Attempt:  unmarshalled.Attempt,
		Sequence: source.Offset,
		Source:   unmarshalled.Source,
	}

}
