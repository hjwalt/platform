package stateful_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/stateful"
	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
	"github.com/stretchr/testify/assert"
)

const testInputMessage = "input"

var testMessageMetadata = flow.Metadata{Id: "id-1", Group: "grp-1", Attempt: 1, Sequence: 2, Source: "src"}

func testOperator(
	stateKey func(context.Context, string) (string, error),
	stateUpdate func(context.Context, string, string) either.Either[string, string],
	operate func(context.Context, string, string) (optional.Optional[string], optional.Optional[string]),
	store *mockStateStore[string],
) (*mockProducer[string], *mockProducer[string], flow.Handler[string]) {
	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](nil)

	operator := stateful.NewOperator(
		"test-operator",
		stateKey,
		stateUpdate,
		operate,
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
		store,
	)
	return outputProducer, errorProducer, operator
}

func TestOperatorStateKeyError(t *testing.T) {
	assert := assert.New(t)

	keyErr := errors.New("state key failed")
	store := newMockStateStore[string]()

	stateUpdateCalled := false
	handlerCalled := false

	_, errorProducer, operator := testOperator(
		func(_ context.Context, _ string) (string, error) { return "", keyErr },
		func(_ context.Context, _ string, _ string) either.Either[string, string] {
			stateUpdateCalled = true
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			handlerCalled = true
			return optional.Empty[string](), optional.Empty[string]()
		},
		store,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	})

	assert.ErrorIs(err, keyErr)
	assert.Zero(store.readCalls)
	assert.Zero(store.writeCalls)
	assert.False(stateUpdateCalled)
	assert.False(handlerCalled)
	assert.Empty(errorProducer.produced)
}

func TestOperatorStateStoreReadError(t *testing.T) {
	assert := assert.New(t)

	readErr := errors.New("state read failed")
	store := newMockStateStore[string]()
	store.readErr = readErr

	stateUpdateCalled := false
	handlerCalled := false

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, _ string) either.Either[string, string] {
			stateUpdateCalled = true
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			handlerCalled = true
			return optional.Empty[string](), optional.Empty[string]()
		},
		store,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	})

	assert.ErrorIs(err, readErr)
	assert.Equal(1, store.readCalls)
	assert.Zero(store.writeCalls)
	assert.False(stateUpdateCalled)
	assert.False(handlerCalled)
	assert.Empty(outputProducer.produced)
	assert.Empty(errorProducer.produced)
}

func TestOperatorStateUpdateRightProducesError(t *testing.T) {
	assert := assert.New(t)

	store := newMockStateStore[string]()

	handlerCalled := false

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, _ string) either.Either[string, string] {
			return either.Right[string, string]("state-error")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			handlerCalled = true
			return optional.Empty[string](), optional.Empty[string]()
		},
		store,
	)

	msg := flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	}

	err := operator.Handle(context.Background(), msg)

	assert.NoError(err)
	assert.Empty(outputProducer.produced)
	assert.Len(errorProducer.produced, 1)
	assert.Equal("state-error", errorProducer.produced[0][0].Value)
	assert.Equal(msg.Metadata.Id, errorProducer.produced[0][0].Metadata.Id)
	assert.Equal("state-error", errorProducer.produced[0][0].Metadata.Group)
	assert.Equal(msg.Metadata.Attempt+1, errorProducer.produced[0][0].Metadata.Attempt)
	assert.False(errorProducer.produced[0][0].Timestamp.IsZero())
	assert.Zero(store.writeCalls)
	assert.False(handlerCalled)
}

func TestOperatorStateUpdateRightProducerErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	producerErr := errors.New("error producer failed")
	store := newMockStateStore[string]()

	outputProducer := newMockProducer[string](nil)
	errorProducer := newMockProducer[string](producerErr)

	operator := stateful.NewOperator(
		"test-operator",
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, _ string) either.Either[string, string] {
			return either.Right[string, string]("state-error")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Empty[string](), optional.Empty[string]()
		},
		entryMetadata,
		entryMetadata,
		outputProducer,
		errorProducer,
		store,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	})

	assert.ErrorIs(err, producerErr)
	assert.Len(errorProducer.produced, 1)
	assert.Empty(outputProducer.produced)
	assert.Zero(store.writeCalls)
}

func TestOperatorStateStoreWriteErrorStopsBeforeHandler(t *testing.T) {
	assert := assert.New(t)

	writeErr := errors.New("state write failed")
	store := newMockStateStore[string]()
	store.writeErr = writeErr

	handlerCalled := false
	updateCalledWith := ""

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, current string) either.Either[string, string] {
			updateCalledWith = current
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			handlerCalled = true
			return optional.Of("output-value"), optional.Empty[string]()
		},
		store,
	)

	// pre-existing state so Read yields a value
	store.states["key-"+testInputMessage] = flow.State[string]{Id: "key-" + testInputMessage, Value: "state-old"}

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	})

	assert.ErrorIs(err, writeErr)
	assert.Equal(1, store.readCalls)
	assert.Equal(1, store.writeCalls)
	// handler must not run before the write has succeeded
	assert.False(handlerCalled)
	assert.Equal("state-old", updateCalledWith)
	assert.Empty(outputProducer.produced)
	assert.Empty(errorProducer.produced)
}

func TestOperatorSuccessResultPresent(t *testing.T) {
	assert := assert.New(t)

	store := newMockStateStore[string]()

	handlerCalledWithState := ""
	stateUpdateCalledWith := ""

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, current string) either.Either[string, string] {
			stateUpdateCalledWith = current
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, state string) (optional.Optional[string], optional.Optional[string]) {
			handlerCalledWithState = state
			return optional.Of("output-value"), optional.Empty[string]()
		},
		store,
	)

	store.states["key-"+testInputMessage] = flow.State[string]{Id: "key-" + testInputMessage, Value: "state-old"}

	msg := flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	}

	err := operator.Handle(context.Background(), msg)

	assert.NoError(err)

	// state was written first: Id from StateKey, Value from StateUpdate left side
	assert.Equal(1, store.writeCalls)
	assert.Equal("key-"+testInputMessage, store.lastWritten.Id)
	assert.Equal("state-new", store.lastWritten.Value)
	assert.False(store.lastWritten.Timestamp.IsZero())
	assert.Equal("state-old", stateUpdateCalledWith)
	assert.Equal("state-new", store.states["key-"+testInputMessage].Value)

	// handler ran with the updated state and its result was produced
	assert.Equal("state-new", handlerCalledWithState)
	assert.Len(outputProducer.produced, 1)
	assert.Equal("output-value", outputProducer.produced[0][0].Value)
	assert.Equal(msg.Metadata.Id, outputProducer.produced[0][0].Metadata.Id)
	assert.Equal("output-value", outputProducer.produced[0][0].Metadata.Group)
	assert.False(outputProducer.produced[0][0].Timestamp.IsZero())
	assert.Empty(errorProducer.produced)
}

func TestOperatorSuccessErrorPresent(t *testing.T) {
	assert := assert.New(t)

	store := newMockStateStore[string]()

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, _ string) either.Either[string, string] {
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Empty[string](), optional.Of("handler-error")
		},
		store,
	)

	msg := flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	}

	err := operator.Handle(context.Background(), msg)

	assert.NoError(err)
	assert.Equal(1, store.writeCalls)
	assert.Equal("key-"+testInputMessage, store.lastWritten.Id)
	assert.Equal("state-new", store.lastWritten.Value)

	assert.Empty(outputProducer.produced)
	assert.Len(errorProducer.produced, 1)
	assert.Equal("handler-error", errorProducer.produced[0][0].Value)
	assert.Equal(msg.Metadata.Id, errorProducer.produced[0][0].Metadata.Id)
	assert.Equal("handler-error", errorProducer.produced[0][0].Metadata.Group)
	assert.False(errorProducer.produced[0][0].Timestamp.IsZero())
}

func TestOperatorSuccessNoResultNoError(t *testing.T) {
	assert := assert.New(t)

	store := newMockStateStore[string]()

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, _ string) either.Either[string, string] {
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Empty[string](), optional.Empty[string]()
		},
		store,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	})

	assert.NoError(err)
	assert.Equal(1, store.writeCalls)
	assert.Equal("key-"+testInputMessage, store.lastWritten.Id)
	assert.Equal("state-new", store.states["key-"+testInputMessage].Value)
	assert.Empty(outputProducer.produced)
	assert.Empty(errorProducer.produced)
}

func TestOperatorSuccessNoPriorState(t *testing.T) {
	assert := assert.New(t)

	store := newMockStateStore[string]()

	stateUpdateCalledWith := "unset"

	outputProducer, errorProducer, operator := testOperator(
		func(_ context.Context, value string) (string, error) { return "key-" + value, nil },
		func(_ context.Context, _ string, current string) either.Either[string, string] {
			stateUpdateCalledWith = current
			return either.Left[string, string]("state-new")
		},
		func(_ context.Context, _ string, _ string) (optional.Optional[string], optional.Optional[string]) {
			return optional.Empty[string](), optional.Empty[string]()
		},
		store,
	)

	err := operator.Handle(context.Background(), flow.Message[string]{
		Metadata: testMessageMetadata,
		Value:    testInputMessage,
	})

	assert.NoError(err)
	// no prior state -> Read yields the zero flow.State
	assert.Equal("", stateUpdateCalledWith)
	assert.Equal(1, store.writeCalls)
	assert.Empty(outputProducer.produced)
	assert.Empty(errorProducer.produced)
}
