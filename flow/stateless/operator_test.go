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

func TestOperatorHandleResultPresent(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	operator := stateless.NewOperator(
		"test-operator",
		func(_ context.Context, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Of("output-value"), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	msg := flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	}

	err := operator.Handle(context.Background(), msg)

	assert.NoError(err)
	assert.Len(outputProducer.produced, 1)
	assert.Empty(errorProducer.produced)

	produced := outputProducer.produced[0]
	assert.Len(produced, 1)
	assert.Equal("output-value", produced[0].Value)
	assert.Equal(msg.Metadata.Id, produced[0].Metadata.Id)
	assert.Equal("output-value", produced[0].Metadata.Group)
	assert.Equal(msg.Metadata.Attempt+1, produced[0].Metadata.Attempt)
	assert.Equal(msg.Metadata.Sequence, produced[0].Metadata.Sequence)
	assert.Equal(msg.Metadata.Source, produced[0].Metadata.Source)
	assert.False(produced[0].Timestamp.IsZero())
}

func TestOperatorHandleOutputProducerErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	producerErr := errors.New("output producer failed")
	outputProducer := newMockProducer[string](producerErr)
	errorProducer := newMockProducer[string](nil)

	operator := stateless.NewOperator(
		"test-operator",
		func(_ context.Context, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Of("output-value"), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.ErrorIs(err, producerErr)
	assert.Len(outputProducer.produced, 1)
	assert.Empty(errorProducer.produced)
}

func TestOperatorHandleErrorPresent(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	operator := stateless.NewOperator(
		"test-operator",
		func(_ context.Context, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Empty[string](), optional.Of("handler-error")
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	msg := flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	}

	err := operator.Handle(context.Background(), msg)

	assert.NoError(err)
	assert.Empty(outputProducer.produced)
	assert.Len(errorProducer.produced, 1)

	produced := errorProducer.produced[0]
	assert.Len(produced, 1)
	assert.Equal("handler-error", produced[0].Value)
	assert.Equal(msg.Metadata.Id, produced[0].Metadata.Id)
	assert.Equal("handler-error", produced[0].Metadata.Group)
	assert.Equal(msg.Metadata.Attempt+1, produced[0].Metadata.Attempt)
	assert.Equal(msg.Metadata.Sequence, produced[0].Metadata.Sequence)
	assert.Equal(msg.Metadata.Source, produced[0].Metadata.Source)
	assert.False(produced[0].Timestamp.IsZero())
}

func TestOperatorHandleNoResultNoError(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	operator := stateless.NewOperator(
		"test-operator",
		func(_ context.Context, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Empty[string](), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Empty(outputProducer.produced)
	assert.Empty(errorProducer.produced)
}

func TestOperatorHandleResultAndErrorPresentResultWins(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	operator := stateless.NewOperator(
		"test-operator",
		func(_ context.Context, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Of("output-value"), optional.Of("handler-error")
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Len(outputProducer.produced, 1)
	assert.Equal("output-value", outputProducer.produced[0][0].Value)
	assert.Empty(errorProducer.produced)
}
