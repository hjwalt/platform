package state_file

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/hjwalt/platform/state"
)

type FileStore struct {
	Path string
}

func (r *FileStore) Read(ctx context.Context, id string) (state.State, error) {
	file := r.Path + id + ".dat"

	bytes, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return state.State{
			Id:        id,
			Value:     []byte{},
			Timestamp: time.Now(),
		}, nil
	}
	if err != nil {
		return state.State{}, errors.Join(ErrReadFail, err)
	}
	return state.State{
		Id:        id,
		Value:     bytes,
		Timestamp: time.Now(),
	}, nil
}

func (r *FileStore) Write(ctx context.Context, state state.State) error {
	file := r.Path + state.Id + ".dat"

	writeErr := os.WriteFile(file, state.Value, 0644)
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
