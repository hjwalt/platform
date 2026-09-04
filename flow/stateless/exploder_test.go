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

func TestExploderHandleSingleEntry(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Of([]string{"only"}), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := exploder.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Len(outputProducer.produced, 1)
	assert.Empty(errorProducer.produced)

	produced := outputProducer.produced[0]
	assert.Len(produced, 1)
	assert.Equal("only", produced[0].Value)
	assert.Equal("only", produced[0].Metadata.Group)
	assert.Equal("id-1", produced[0].Metadata.Id)
	assert.False(produced[0].Timestamp.IsZero())
}

// TestExploderHandleMultipleEntries pins down the CURRENT behaviour: the
// implementation returns from inside the result loop, so only the first entry
// of a multi-entry result is ever produced. The remaining entries are dropped.
func TestExploderHandleMultipleEntriesOnlyFirstProduced(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Of([]string{"first", "second", "third"}), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := exploder.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	// actual behaviour: one call, one message, only the first entry
	assert.Len(outputProducer.produced, 1)
	assert.Len(outputProducer.produced[0], 1)
	assert.Equal("first", outputProducer.produced[0][0].Value)
	assert.Empty(errorProducer.produced)
}

func TestExploderHandleErrorPresent(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Empty[[]string](), optional.Of("explode-error")
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

	err := exploder.Handle(context.Background(), msg)

	assert.NoError(err)
	assert.Empty(outputProducer.produced)
	assert.Len(errorProducer.produced, 1)
	assert.Equal("explode-error", errorProducer.produced[0][0].Value)
	assert.Equal(msg.Metadata.Id, errorProducer.produced[0][0].Metadata.Id)
	assert.Equal("explode-error", errorProducer.produced[0][0].Metadata.Group)
	assert.False(errorProducer.produced[0][0].Timestamp.IsZero())
}

func TestExploderHandleNoResultNoError(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Empty[[]string](), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := exploder.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Empty(outputProducer.produced)
	assert.Empty(errorProducer.produced)
}

// TestExploderHandleEmptyResultWithError: an empty-but-present result does not
// enter the loop, so control falls through to the error branch.
func TestExploderHandleEmptyResultWithError(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Of([]string{}), optional.Of("explode-error")
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := exploder.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Empty(outputProducer.produced)
	assert.Len(errorProducer.produced, 1)
	assert.Equal("explode-error", errorProducer.produced[0][0].Value)
}

// TestExploderHandleEntriesAndError: when the result slice is non-empty AND an
// error is present, the result branch is taken first and the error is dropped
// after the first entry is produced.
func TestExploderHandleEntriesAndError(t *testing.T) {
	assert := assert.New(t)

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Of([]string{"first", "second"}), optional.Of("explode-error")
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := exploder.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.NoError(err)
	assert.Len(outputProducer.produced, 1)
	assert.Equal("first", outputProducer.produced[0][0].Value)
	assert.Empty(errorProducer.produced)
}

func TestExploderHandleProducerErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	producerErr := errors.New("output producer failed")
	outputProducer := newMockProducer[string](producerErr)
	errorProducer := newMockProducer[string](nil)

	exploder := stateless.NewExploder(
		"test-exploder",
		func(_ context.Context, _ string) (optional.Optional[[]string], optional.Optional[string]) {
			return optional.Of([]string{"first", "second"}), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
	)

	err := exploder.Handle(context.Background(), flow.Message[string]{
		Metadata: flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"},
		Value:    "input",
	})

	assert.ErrorIs(err, producerErr)
	assert.Len(outputProducer.produced, 1)
}
