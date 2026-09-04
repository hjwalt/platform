package flow_runtime_kafka_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/stretchr/testify/assert"
)

func TestKafkaRuntimeFlowToRuntimeMetadata(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_kafka.New("topic-a")

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

	// topic comes from the runtime configuration, key from the source group
	assert.Equal("topic-a", result.Topic)
	assert.Equal("grp-1", result.Key)
	assert.Equal(1, len(result.Headers))
	assert.Contains(result.Headers, metadata.MetadataHeaderKey)
	assert.Equal(string(expected), result.Headers[metadata.MetadataHeaderKey])
	// offset and partition are not carried in this direction
	assert.Equal(int64(0), result.Offset)
	assert.Equal(int32(0), result.Partition)
}

func TestKafkaRuntimeRuntimeToFlowMetadata(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_kafka.New("topic-a")

	header := flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 9,
		Source:   "src-1",
	}
	content, err := metadata.Format.Marshal(header)
	assert.NoError(err)

	result := runtime.RuntimeToFlowMetadata(kafka.KafkaMetadata{
		Topic:     "other-topic",
		Offset:    42,
		Partition: 7,
		Key:       "key-from-runtime",
		Headers:   map[string]string{metadata.MetadataHeaderKey: string(content)},
	})

	// Id, Attempt and Source come from the header; Group is overridden by the
	// kafka message key and Sequence by the kafka offset
	assert.Equal(flow.Metadata{
		Id:       "mid-1",
		Group:    "key-from-runtime",
		Attempt:  3,
		Sequence: 42,
		Source:   "src-1",
	}, result)
}

func TestKafkaRuntimeMetadataRoundTrip(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_kafka.New("topic-a")

	source := flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 9,
		Source:   "src-1",
	}

	back := runtime.RuntimeToFlowMetadata(runtime.FlowToRuntimeMetadata(source))

	// Id, Group, Attempt and Source survive; Sequence does not because
	// FlowToRuntimeMetadata never sets Offset (it stays 0) and
	// RuntimeToFlowMetadata takes Sequence from Offset
	assert.Equal("mid-1", back.Id)
	assert.Equal("grp-1", back.Group)
	assert.Equal(int32(3), back.Attempt)
	assert.Equal("src-1", back.Source)
	assert.Equal(int64(0), back.Sequence)
}

func TestKafkaRuntimeMetadataMissingHeaderReturnsDefault(t *testing.T) {
	assert := assert.New(t)

	runtime := flow_runtime_kafka.New("topic-a")

	result := runtime.RuntimeToFlowMetadata(kafka.KafkaMetadata{
		Topic:   "topic-a",
		Headers: map[string]string{},
	})

	assertDefaultMetadata(assert, result)
}

func TestKafkaRuntimeMetadataGarbageHeaderReturnsDefault(t *testing.T) {
	assert := assert.New(t)

	// silence the slog.Error the implementation logs on the garbage path
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	runtime := flow_runtime_kafka.New("topic-a")

	result := runtime.RuntimeToFlowMetadata(kafka.KafkaMetadata{
		Topic:   "topic-a",
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
