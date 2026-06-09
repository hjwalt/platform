package file_store

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/hjwalt/platform/state"
)

func New(config Configuration) state.Store {
	return &store{
		Path: config.Path,
	}
}

type Configuration struct {
	Path string
}

type store struct {
	Path string
}

func (r *store) Read(ctx context.Context, id string) (state.State, error) {
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

func (r *store) Write(ctx context.Context, state state.State) error {
	file := r.Path + state.Id + ".dat"

	writeErr := os.WriteFile(file, state.Value, 0644)
	if writeErr != nil {
		return errors.Join(ErrWriteFail, writeErr)
	}

	return nil
}

func (r *store) Delete(ctx context.Context, id string) error {
	file := r.Path + id + ".dat"

	return os.Remove(file)
}

func (r *store) Keys(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(r.Path)
	if err != nil {
		return []string{}, err
	}

	keys := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			keys = append(keys, strings.TrimSuffix(entry.Name(), ".dat"))
		}
	}

	return keys, nil
}

func (r *store) Start() error {
	return nil
}

func (r *store) Stop() {

}

var (
	ErrWriteFail = errors.New("cannot write file")
	ErrReadFail  = errors.New("cannot read file")
)
