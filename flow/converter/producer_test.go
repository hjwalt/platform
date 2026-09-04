package converter_test

import (
	"context"
	"testing"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/format"
	"github.com/stretchr/testify/assert"
)

// metadataFunc prefixes the value id into the Group and keeps the rest of the
// previous metadata intact, so tests can assert both the Default() base and
// the extraction result.
func metadataFunc(_ context.Context, previous flow.Metadata, value myMessage) flow.Metadata {
	return flow.Metadata{
		Id:       previous.Id,
		Group:    "g-" + value.Id,
		Attempt:  previous.Attempt,
		Sequence: previous.Sequence,
		Source:   previous.Source,
	}
}

func TestFlowProducerProduce(t *testing.T) {
	assert := assert.New(t)

	rawProducer := newMockMessageProducer[mockMeta](nil)
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	producer := converter.RuntimeToFlowProducer(rawProducer, conv, metadataFunc)

	values := []myMessage{
		{Id: "value-1", Count: 2},
		{Id: "value-2", Count: 5},
	}

	err := producer.Produce(context.Background(), values)

	assert.NoError(err)
	assert.Len(rawProducer.produced, 1)
	assert.Len(rawProducer.produced[0], 2)

	first := rawProducer.produced[0][0]
	second := rawProducer.produced[0][1]

	// metadata base comes from metadata.Default() (Attempt 0, Sequence -1,
	// Source UNKNOWN) with a fresh Id, then the extract func ran on top
	assert.NotEmpty(first.Metadata.Id)
	assert.Equal("g-value-1", first.Metadata.Group)
	assert.Equal(int32(0), first.Metadata.Attempt)
	assert.Equal(int64(-1), first.Metadata.Sequence)
	assert.Equal("UNKNOWN", first.Metadata.Source)

	assert.NotEmpty(second.Metadata.Id)
	assert.Equal("g-value-2", second.Metadata.Group)

	// one timestamp per Produce call
	assert.True(first.Timestamp.Equal(second.Timestamp))
	assert.False(first.Timestamp.IsZero())

	// values converted via the value format
	valueFormat := format.Json[myMessage]()
	for i, produced := range rawProducer.produced[0] {
		unmarshalled, err := valueFormat.Unmarshal(produced.Value)
		assert.NoError(err)
		assert.Equal(values[i], unmarshalled)
	}
}

func TestFlowProducerProduceMessage(t *testing.T) {
	assert := assert.New(t)

	rawProducer := newMockMessageProducer[mockMeta](nil)
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	producer := converter.RuntimeToFlowProducer(rawProducer, conv, metadataFunc)

	timestamp := time.Now()

	err := producer.ProduceMessage(context.Background(), []flow.Message[myMessage]{
		{
			Metadata: flow.Metadata{
				Id:       "mid-1",
				Group:    "grp-1",
				Attempt:  3,
				Sequence: 7,
				Source:   "src-1",
			},
			Value:     myMessage{Id: "value-1", Count: 2},
			Timestamp: timestamp,
		},
	})

	assert.NoError(err)
	assert.Len(rawProducer.produced, 1)
	assert.Len(rawProducer.produced[0], 1)

	produced := rawProducer.produced[0][0]
	assert.Equal(mockMeta{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 7,
		Source:   "src-1",
	}, produced.Metadata)
	assert.Equal(timestamp, produced.Timestamp)

	valueFormat := format.Json[myMessage]()
	unmarshalled, err := valueFormat.Unmarshal(produced.Value)
	assert.NoError(err)
	assert.Equal(myMessage{Id: "value-1", Count: 2}, unmarshalled)
}

func TestFlowProducerProduceConversionErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	rawProducer := newMockMessageProducer[mockMeta](nil)
	conv := converter.NewConverter(mockMessageRuntime{}, format.Broken())

	producer := converter.RuntimeToFlowProducer(
		rawProducer,
		conv,
		func(_ context.Context, previous flow.Metadata, _ string) flow.Metadata {
			return previous
		},
	)

	// second value fails to marshal -> conversion error, underlying producer
	// must not be called
	err := producer.Produce(context.Background(), []string{"ok", "marshal"})

	assert.ErrorIs(err, format.ErrMarshal)
	assert.Empty(rawProducer.produced)
}

func TestFlowProducerProduceMessageConversionErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	rawProducer := newMockMessageProducer[mockMeta](nil)
	conv := converter.NewConverter(mockMessageRuntime{}, format.Broken())

	producer := converter.RuntimeToFlowProducer(
		rawProducer,
		conv,
		func(_ context.Context, previous flow.Metadata, _ string) flow.Metadata {
			return previous
		},
	)

	err := producer.ProduceMessage(context.Background(), []flow.Message[string]{
		{Metadata: flow.Metadata{Id: "mid-1"}, Value: "marshal"},
	})

	assert.ErrorIs(err, format.ErrMarshal)
	assert.Empty(rawProducer.produced)
}

func TestFlowProducerUnderlyingProducerErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	producerErr := testErr
	rawProducer := newMockMessageProducer[mockMeta](producerErr)
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	producer := converter.RuntimeToFlowProducer(rawProducer, conv, metadataFunc)

	err := producer.ProduceMessage(context.Background(), []flow.Message[myMessage]{
		{
			Metadata: flow.Metadata{Id: "mid-1"},
			Value:    myMessage{Id: "value-1", Count: 2},
		},
	})

	assert.ErrorIs(err, producerErr)
	assert.Len(rawProducer.produced, 1)
}
