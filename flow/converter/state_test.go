package converter_test

import (
	"context"
	"testing"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/state"
	"github.com/stretchr/testify/assert"
)

func TestFlowStoreRead(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	valueFormat := format.Json[string]()
	store := converter.RuntimeToFlowStore(raw, valueFormat)

	valueBytes, err := valueFormat.Marshal("hello")
	assert.NoError(err)
	raw.states["key-1"] = state.State{Id: "key-1", Value: valueBytes, Timestamp: time.Now()}

	result, err := store.Read(context.Background(), "key-1")

	assert.NoError(err)
	assert.Equal("key-1", result.Id)
	assert.Equal("hello", result.Value)
	assert.False(result.Timestamp.IsZero())
}

func TestFlowStoreReadErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	raw.readErr = testErr
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	_, err := store.Read(context.Background(), "key-1")

	assert.ErrorIs(err, testErr)
}

func TestFlowStoreReadUnmarshalError(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	store := converter.RuntimeToFlowStore(raw, format.Broken())

	// BrokenFormat.Unmarshal fails on the value "unmarshal"
	raw.states["key-1"] = state.State{Id: "key-1", Value: []byte("unmarshal")}

	_, err := store.Read(context.Background(), "key-1")

	assert.Error(err)
	assert.ErrorIs(err, converter.ErrWriteUnmarshal)
	assert.ErrorIs(err, format.ErrUnmarshal)
}

func TestFlowStoreReadIdComesFromRequest(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	store := converter.RuntimeToFlowStore(raw, format.Json[myMessage]())

	// no entry under the requested id: empty bytes unmarshal to the zero
	// value without error
	result, err := store.Read(context.Background(), "requested-id")

	assert.NoError(err)
	assert.Equal("requested-id", result.Id)
	assert.Equal(myMessage{}, result.Value)
}

func TestFlowStoreWriteAndReadRoundTrip(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	timestamp := time.Now()

	err := store.Write(context.Background(), flow.State[string]{
		Id:        "key-1",
		Value:     "hello",
		Timestamp: timestamp,
	})
	assert.NoError(err)

	// underlying store received marshalled bytes and forwarded timestamp
	rawEntry, exists := raw.states["key-1"]
	assert.True(exists)
	assert.Equal(state.State{Id: "key-1", Value: []byte(`"hello"`), Timestamp: timestamp}, rawEntry)

	result, err := store.Read(context.Background(), "key-1")
	assert.NoError(err)
	assert.Equal("key-1", result.Id)
	assert.Equal("hello", result.Value)
}

func TestFlowStoreWriteMarshalError(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	store := converter.RuntimeToFlowStore(raw, format.Broken())

	// BrokenFormat.Marshal fails on the value "marshal"
	err := store.Write(context.Background(), flow.State[string]{
		Id:    "key-1",
		Value: "marshal",
	})

	assert.Error(err)
	assert.ErrorIs(err, converter.ErrWriteMarshal)
	assert.ErrorIs(err, format.ErrMarshal)
	// nothing delegated to the underlying store
	assert.Empty(raw.states)
	assert.Zero(raw.writeCalls)
}

func TestFlowStoreWriteUnderlyingErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	raw.writeErr = testErr
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	err := store.Write(context.Background(), flow.State[string]{
		Id:    "key-1",
		Value: "hello",
	})

	assert.ErrorIs(err, testErr)
}

func TestFlowStoreKeysDelegates(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	raw.states["key-1"] = state.State{Id: "key-1"}
	raw.states["key-2"] = state.State{Id: "key-2"}
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	keys, err := store.Keys(context.Background())

	assert.NoError(err)
	assert.ElementsMatch([]string{"key-1", "key-2"}, keys)
	assert.Equal(1, raw.keysCalls)
}

func TestFlowStoreKeysErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	raw.keysErr = testErr
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	_, err := store.Keys(context.Background())

	assert.ErrorIs(err, testErr)
}

func TestFlowStoreStartStopDelegate(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	assert.NoError(store.Start())
	assert.Equal(1, raw.startCalls)

	store.Stop()
	assert.Equal(1, raw.stopCalls)
}

func TestFlowStoreStartErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	raw := newMockRawStateStore()
	raw.startErr = testErr
	store := converter.RuntimeToFlowStore(raw, format.Json[string]())

	err := store.Start()

	assert.ErrorIs(err, testErr)
}
