package file_store_test

import (
	"context"
	"os"
	"testing"

	"github.com/hjwalt/platform/state"
	file_store "github.com/hjwalt/platform/state/file"
	"github.com/stretchr/testify/assert"
)

// tempStorePath returns a fresh temporary directory with a trailing separator,
// because the store concatenates Path + id + ".dat" without adding one.
func tempStorePath(t *testing.T) string {
	t.Helper()
	return t.TempDir() + string(os.PathSeparator)
}

func TestReadMissingId(t *testing.T) {
	assert := assert.New(t)

	store := file_store.New(file_store.Configuration{Path: tempStorePath(t)})

	got, err := store.Read(context.Background(), "missing")

	// Reading a missing id returns the id with an empty, non-nil Value and no error.
	assert.NoError(err)
	assert.Equal("missing", got.Id)
	assert.NotNil(got.Value)
	assert.Equal([]byte{}, got.Value)
}

func TestWriteReadRoundtrip(t *testing.T) {
	assert := assert.New(t)

	store := file_store.New(file_store.Configuration{Path: tempStorePath(t)})

	want := state.State{Id: "roundtrip", Value: []byte("hello world")}
	err := store.Write(context.Background(), want)
	assert.NoError(err)

	got, err := store.Read(context.Background(), "roundtrip")

	assert.NoError(err)
	assert.Equal(want.Id, got.Id)
	assert.Equal(want.Value, got.Value)
}

func TestWriteOverwritesExistingValue(t *testing.T) {
	assert := assert.New(t)

	store := file_store.New(file_store.Configuration{Path: tempStorePath(t)})

	assert.NoError(store.Write(context.Background(), state.State{Id: "ow", Value: []byte("first")}))
	assert.NoError(store.Write(context.Background(), state.State{Id: "ow", Value: []byte("second")}))

	got, err := store.Read(context.Background(), "ow")

	assert.NoError(err)
	assert.Equal([]byte("second"), got.Value)
}

func TestWriteError(t *testing.T) {
	assert := assert.New(t)

	// Path points into a directory that does not exist, so os.WriteFile fails.
	store := file_store.New(file_store.Configuration{Path: tempStorePath(t) + "nope" + string(os.PathSeparator)})

	err := store.Write(context.Background(), state.State{Id: "id", Value: []byte("value")})

	assert.Error(err)
	assert.ErrorIs(err, file_store.ErrWriteFail)
}

func TestReadError(t *testing.T) {
	assert := assert.New(t)

	path := tempStorePath(t)
	// A directory named "<id>.dat" makes os.ReadFile fail with a non-NotExist error.
	assert.NoError(os.Mkdir(path+"blocked.dat", 0755))

	store := file_store.New(file_store.Configuration{Path: path})

	got, err := store.Read(context.Background(), "blocked")

	assert.Error(err)
	assert.ErrorIs(err, file_store.ErrReadFail)
	assert.Equal(state.State{}, got)
}

func TestDeleteRemovesFile(t *testing.T) {
	assert := assert.New(t)

	path := tempStorePath(t)
	store := file_store.New(file_store.Configuration{Path: path})

	assert.NoError(store.Write(context.Background(), state.State{Id: "delme", Value: []byte("value")}))

	_, statErr := os.Stat(path + "delme.dat")
	assert.NoError(statErr, "file should exist before delete")

	assert.NoError(store.Delete(context.Background(), "delme"))

	_, statErr = os.Stat(path + "delme.dat")
	assert.ErrorIs(statErr, os.ErrNotExist, "file should be gone after delete")

	// A subsequent read returns the missing-file state without error.
	got, err := store.Read(context.Background(), "delme")
	assert.NoError(err)
	assert.Equal("delme", got.Id)
	assert.Equal([]byte{}, got.Value)
}

func TestDeleteMissingIdReturnsError(t *testing.T) {
	assert := assert.New(t)

	store := file_store.New(file_store.Configuration{Path: tempStorePath(t)})

	// os.Remove's raw error is returned (not wrapped in a store error).
	err := store.Delete(context.Background(), "never-written")
	assert.Error(err)
	assert.ErrorIs(err, os.ErrNotExist)
}

func TestKeysListsWrittenIdsAndIgnoresDirectories(t *testing.T) {
	assert := assert.New(t)

	path := tempStorePath(t)
	store := file_store.New(file_store.Configuration{Path: path})

	// Write in a deliberately non-sorted order.
	assert.NoError(store.Write(context.Background(), state.State{Id: "beta", Value: []byte("1")}))
	assert.NoError(store.Write(context.Background(), state.State{Id: "alpha", Value: []byte("2")}))
	assert.NoError(store.Write(context.Background(), state.State{Id: "gamma", Value: []byte("3")}))

	// Directories (with and without a ".dat" suffix) must be ignored.
	assert.NoError(os.Mkdir(path+"subdir", 0755))
	assert.NoError(os.Mkdir(path+"subdir.dat", 0755))

	// Oddity: files without a ".dat" suffix are listed as keys (the suffix is
	// trimmed unconditionally from any non-directory entry name).
	assert.NoError(os.WriteFile(path+"notes.txt", []byte("hi"), 0644))

	keys, err := store.Keys(context.Background())

	assert.NoError(err)
	assert.ElementsMatch([]string{"alpha", "beta", "gamma", "notes.txt"}, keys)
}

func TestKeysOnMissingDirectoryReturnsError(t *testing.T) {
	assert := assert.New(t)

	store := file_store.New(file_store.Configuration{Path: tempStorePath(t) + "missing" + string(os.PathSeparator)})

	keys, err := store.Keys(context.Background())

	assert.Error(err)
	assert.ErrorIs(err, os.ErrNotExist)
	assert.Equal([]string{}, keys)
}

func TestStartStopNoOp(t *testing.T) {
	assert := assert.New(t)

	store := file_store.New(file_store.Configuration{Path: tempStorePath(t)})

	assert.NoError(store.Start())
	store.Stop()
	assert.NoError(store.Start())
	store.Stop()
}
