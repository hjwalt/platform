package stateful_test

import (
	"context"

	"github.com/hjwalt/platform/flow"
)

// mockProducer implements flow.Producer[V]. It records a copy of the slice
// passed to each ProduceMessage call and can be configured to return an error.
type mockProducer[V any] struct {
	produced [][]flow.Message[V]
	err      error
}

func newMockProducer[V any](err error) *mockProducer[V] {
	return &mockProducer[V]{err: err}
}

func (p *mockProducer[V]) ProduceMessage(_ context.Context, messages []flow.Message[V]) error {
	call := make([]flow.Message[V], len(messages))
	copy(call, messages)
	p.produced = append(p.produced, call)
	return p.err
}

func (p *mockProducer[V]) Produce(_ context.Context, _ []V) error { return nil }

func (p *mockProducer[V]) Start() error { return nil }

func (p *mockProducer[V]) Stop() {}

// mockStateStore implements flow.Store[ST], backed by a map of flow.State.
type mockStateStore[ST any] struct {
	states      map[string]flow.State[ST]
	readErr     error
	writeErr    error
	keysErr     error
	startErr    error
	readCalls   int
	writeCalls  int
	keysCalls   int
	startCalls  int
	stopCalls   int
	lastWritten flow.State[ST]
}

func newMockStateStore[ST any]() *mockStateStore[ST] {
	return &mockStateStore[ST]{states: map[string]flow.State[ST]{}}
}

func (s *mockStateStore[ST]) Read(_ context.Context, id string) (flow.State[ST], error) {
	s.readCalls++
	if s.readErr != nil {
		return flow.State[ST]{}, s.readErr
	}
	return s.states[id], nil
}

func (s *mockStateStore[ST]) Write(_ context.Context, state flow.State[ST]) error {
	s.writeCalls++
	s.lastWritten = state
	if s.writeErr != nil {
		return s.writeErr
	}
	s.states[state.Id] = state
	return nil
}

func (s *mockStateStore[ST]) Keys(_ context.Context) ([]string, error) {
	s.keysCalls++
	if s.keysErr != nil {
		return nil, s.keysErr
	}
	keys := make([]string, 0, len(s.states))
	for id := range s.states {
		keys = append(keys, id)
	}
	return keys, nil
}

func (s *mockStateStore[ST]) Start() error {
	s.startCalls++
	return s.startErr
}

func (s *mockStateStore[ST]) Stop() {
	s.stopCalls++
}

// entryMetadata stamps the produced/errored value into the Group field and
// carries the other fields of the incoming message metadata through, so tests
// can assert that the input metadata was forwarded into the extract function.
func entryMetadata(_ context.Context, previous flow.Metadata, value string) flow.Metadata {
	return flow.Metadata{
		Id:       previous.Id,
		Group:    value,
		Attempt:  previous.Attempt + 1,
		Sequence: previous.Sequence,
		Source:   previous.Source,
	}
}
