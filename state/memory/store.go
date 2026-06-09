package memory_store

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/hjwalt/platform/state"
)

func New() state.Store {
	return &store{
		states: map[string][]byte{},
	}
}

type store struct {
	mutex sync.RWMutex

	states map[string][]byte
}

func (s *store) Read(ctx context.Context, id string) (state.State, error) {
	if err := ctx.Err(); err != nil {
		return state.State{}, err
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.states[id]

	if !ok {
		return state.State{
			Id:        id,
			Value:     []byte{},
			Timestamp: time.Now(),
		}, nil
	}

	copied := append([]byte(nil), value...)
	return state.State{
		Id:        id,
		Value:     copied,
		Timestamp: time.Now(),
	}, nil
}

func (s *store) Write(ctx context.Context, input state.State) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if input.Id == "" {
		return ErrInvalidID
	}

	copied := append([]byte(nil), input.Value...)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.states[input.Id] = copied

	return nil
}

func (s *store) Delete(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.states, id)
	return nil
}

func (s *store) Keys(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return []string{}, err
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	keys := make([]string, 0, len(s.states))
	for key := range s.states {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys, nil
}

func (s *store) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.states == nil {
		s.states = map[string][]byte{}
	}

	return nil
}

func (s *store) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.states = nil
}

var ErrInvalidID = errors.New("state id cannot be empty")
