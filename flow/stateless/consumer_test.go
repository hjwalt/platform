package stateless_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/stateless"
	"github.com/hjwalt/platform/type/optional"
	"github.com/stretchr/testify/assert"
)

func TestConsumerHandleNoError(t *testing.T) {
	assert := assert.New(t)

	errorProducer := newMockProducer[string](nil)

	consumer := stateless.NewConsumer(
		func(context.Context, string) optional.Optional[string] {
			return optional.Empty[string]()
		},
		entryMetadata,
		errorProducer,
	)

	err := consumer.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Empty(errorProducer.produced)
}

func TestConsumerHandleErrorProduced(t *testing.T) {
	assert := assert.New(t)

	errorProducer := newMockProducer[string](nil)
	errorValue := "boom"

	consumer := stateless.NewConsumer(
		func(_ context.Context, _ string) optional.Optional[string] {
			return optional.Of(errorValue)
		},
		entryMetadata,
		errorProducer,
	)

	msg := flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	}

	err := consumer.Handle(context.Background(), msg)

	assert.NoError(err)
	assert.Len(errorProducer.produced, 1)

	produced := errorProducer.produced[0]
	assert.Len(produced, 1)
	assert.Equal(errorValue, produced[0].Value)
	// metadata built from ErrorMetadata(ctx, msg.Metadata, errorValue)
	assert.Equal(msg.Metadata.Id, produced[0].Metadata.Id)
	assert.Equal(errorValue, produced[0].Metadata.Group)
	assert.Equal(msg.Metadata.Attempt+1, produced[0].Metadata.Attempt)
	assert.Equal(msg.Metadata.Sequence, produced[0].Metadata.Sequence)
	assert.Equal(msg.Metadata.Source, produced[0].Metadata.Source)
	assert.False(produced[0].Timestamp.IsZero())
}

func TestConsumerHandleProducerErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	producerErr := errors.New("producer failed")
	errorProducer := newMockProducer[string](producerErr)

	consumer := stateless.NewConsumer(
		func(_ context.Context, _ string) optional.Optional[string] {
			return optional.Of("boom")
		},
		entryMetadata,
		errorProducer,
	)

	err := consumer.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.ErrorIs(err, producerErr)
	assert.Len(errorProducer.produced, 1)
}
