package file

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/state"
)

type FileStore struct {
	Path string
}

func (r *FileStore) Read(ctx context.Context, id string) (state.State[structure.Bytes], error) {
	file := r.Path + id + ".dat"

	bytes, err := os.ReadFile(file)
	if err != nil {
		return state.State[structure.Bytes]{}, errors.Join(ErrReadFail, err)
	}
	return state.State[structure.Bytes]{
		Id:        id,
		Value:     bytes,
		Timestamp: time.Now(),
	}, nil
}

func (r *FileStore) Write(ctx context.Context, state state.State[structure.Bytes]) error {
	file := r.Path + state.Id + ".dat"

	writeErr := os.WriteFile(file, state.Value, os.ModeExclusive)
	if writeErr != nil {
		return errors.Join(ErrWriteFail, writeErr)
	}

	return nil
}

func (r *FileStore) Start() error {
	return nil
}

func (r *FileStore) Stop() {

}

var (
	ErrWriteFail = errors.New("cannot write file")
	ErrReadFail  = errors.New("cannot read file")
)
