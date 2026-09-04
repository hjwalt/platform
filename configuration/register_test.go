package configuration

import (
	"context"
	"os"
	"testing"

	"github.com/hjwalt/platform/state"
	file_store "github.com/hjwalt/platform/state/file"
	"github.com/stretchr/testify/assert"
)

// tempStoreDir returns a fresh temporary directory with a trailing separator,
// because the file store concatenates Path + id + ".dat" without adding one.
func tempStoreDir(t *testing.T) string {
	t.Helper()
	return t.TempDir() + string(os.PathSeparator)
}

// exerciseStore verifies basic file store behaviour for a store obtained
// through the registration wiring.
func exerciseStore(t *testing.T, store state.Store, id string, value []byte) {
	t.Helper()
	assert := assert.New(t)

	writeErr := store.Write(context.Background(), state.State{Id: id, Value: value})
	assert.NoError(writeErr)

	got, readErr := store.Read(context.Background(), id)
	assert.NoError(readErr)
	assert.Equal(id, got.Id)
	assert.Equal(value, got.Value)

	keys, keysErr := store.Keys(context.Background())
	assert.NoError(keysErr)
	assert.Contains(keys, id)

	assert.NoError(store.Delete(context.Background(), id))

	keysAfter, keysErr := store.Keys(context.Background())
	assert.NoError(keysErr)
	assert.NotContains(keysAfter, id)
}

func TestRegisterAgentHarnessStore(t *testing.T) {
	assert := assert.New(t)

	conf := Configuration{
		Store: StoreConfiguration{
			Agent: file_store.Configuration{Path: tempStoreDir(t)},
		},
	}

	ctx := ContextBuilder()
	RegisterAgentHarnessStore(ctx, conf)

	store := ctx.GetAgentHarnessStore()
	assert.NotNil(store)

	// the store should also be registered as a runtime on the holder
	holder, ok := ctx.(*holder)
	assert.True(ok)
	assert.Len(holder.Runtimes, 1)
	assert.Same(store, holder.Runtimes[0])

	exerciseStore(t, store, "agent-state", []byte("harness value"))
}

func TestRegisterMemoryStore(t *testing.T) {
	assert := assert.New(t)

	conf := Configuration{
		Store: StoreConfiguration{
			Memory: file_store.Configuration{Path: tempStoreDir(t)},
		},
	}

	ctx := ContextBuilder()
	RegisterMemoryStore(ctx, conf)

	store := ctx.GetMemoryStore()
	assert.NotNil(store)

	holder, ok := ctx.(*holder)
	assert.True(ok)
	assert.Len(holder.Runtimes, 1)
	assert.Same(store, holder.Runtimes[0])

	exerciseStore(t, store, "memory-state", []byte("memory value"))
}
