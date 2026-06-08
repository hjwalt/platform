package memory_store

import (
	"context"
	"sync"
	"testing"

	"github.com/hjwalt/platform/state"
	"github.com/stretchr/testify/assert"
)

func TestStoreReadWrite(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()

	startErr := store.Start()
	writeErr := store.Write(ctx, state.State{Id: "my-id", Value: []byte("hello")})
	result, readErr := store.Read(ctx, "my-id")

	assert.NoError(startErr)
	assert.NoError(writeErr)
	assert.NoError(readErr)
	assert.Equal("my-id", result.Id)
	assert.Equal([]byte("hello"), result.Value)
	assert.False(result.Timestamp.IsZero())
}

func TestStoreReadMissing(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()

	result, readErr := store.Read(ctx, "missing")

	assert.NoError(readErr)
	assert.Equal("missing", result.Id)
	assert.Equal([]byte{}, result.Value)
	assert.False(result.Timestamp.IsZero())
}

func TestStoreWriteValidateID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()

	writeErr := store.Write(ctx, state.State{Id: "", Value: []byte("x")})

	assert.ErrorIs(writeErr, ErrInvalidID)
}

func TestStoreKeysDeterministic(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()

	assert.NoError(store.Write(ctx, state.State{Id: "b", Value: []byte("2")}))
	assert.NoError(store.Write(ctx, state.State{Id: "a", Value: []byte("1")}))
	assert.NoError(store.Write(ctx, state.State{Id: "c", Value: []byte("3")}))

	keys, keysErr := store.Keys(ctx)

	assert.NoError(keysErr)
	assert.Equal([]string{"a", "b", "c"}, keys)
}

func TestStoreReadUsesDefensiveCopy(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()

	input := []byte("value")
	assert.NoError(store.Write(ctx, state.State{Id: "copy", Value: input}))

	input[0] = 'X'

	result, readErr := store.Read(ctx, "copy")
	assert.NoError(readErr)
	assert.Equal([]byte("value"), result.Value)

	result.Value[1] = 'Y'
	again, readAgainErr := store.Read(ctx, "copy")
	assert.NoError(readAgainErr)
	assert.Equal([]byte("value"), again.Value)
}

func TestStoreConcurrentUpdates(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()

	const workers = 200
	const keyCount = 20

	wait := sync.WaitGroup{}
	wait.Add(workers)

	for i := 0; i < workers; i++ {
		index := i
		go func() {
			defer wait.Done()
			id := string(rune('a' + (index % keyCount)))
			_ = store.Write(ctx, state.State{Id: id, Value: []byte("v")})
			_, _ = store.Read(ctx, id)
			_, _ = store.Keys(ctx)
		}()
	}

	wait.Wait()

	keys, keysErr := store.Keys(ctx)
	assert.NoError(keysErr)
	assert.Len(keys, keyCount)

	for _, key := range keys {
		result, readErr := store.Read(ctx, key)
		assert.NoError(readErr)
		assert.Equal(key, result.Id)
		assert.NotEmpty(result.Value)
	}
}

func TestStoreStopClearsState(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	store := New()
	assert.NoError(store.Write(ctx, state.State{Id: "x", Value: []byte("1")}))

	store.Stop()

	result, readErr := store.Read(ctx, "x")
	assert.NoError(readErr)
	assert.Equal([]byte{}, result.Value)

	assert.NoError(store.Start())
	keys, keysErr := store.Keys(ctx)
	assert.NoError(keysErr)
	assert.Empty(keys)
}
