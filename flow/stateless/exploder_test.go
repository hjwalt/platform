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

// TestExploderHandleMultipleEntriesAllProduced pins down the fixed behaviour:
// every entry of a multi-entry result is produced as its own message, in
// order, and no error is emitted when the handler returned none.
func TestExploderHandleMultipleEntriesAllProduced(t *testing.T) {
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
	assert.Len(outputProducer.produced, 3)
	assert.Empty(errorProducer.produced)

	wantValues := []string{"first", "second", "third"}
	for i, want := range wantValues {
		assert.Len(outputProducer.produced[i], 1)
		assert.Equal(want, outputProducer.produced[i][0].Value)
		assert.Equal(want, outputProducer.produced[i][0].Metadata.Group)
		assert.Equal("id-1", outputProducer.produced[i][0].Metadata.Id)
		assert.False(outputProducer.produced[i][0].Timestamp.IsZero())
	}
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
// error is present, all entries are produced and the error branch is then
// reached as well, so the error message is emitted after the entries.
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
	assert.Len(outputProducer.produced, 2)
	assert.Equal("first", outputProducer.produced[0][0].Value)
	assert.Equal("second", outputProducer.produced[1][0].Value)
	assert.Len(errorProducer.produced, 1)
	assert.Equal("explode-error", errorProducer.produced[0][0].Value)
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
