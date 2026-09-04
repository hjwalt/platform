package converter_test

import (
	"context"
	"errors"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/state"
)

// testErr is a shared sentinel error for producer/store failure paths.
var testErr = errors.New("test error")

// mockMeta is the runtime metadata type used throughout the converter tests.
// mockMessageRuntime maps it to/from flow.Metadata field-by-field.
type mockMeta struct {
	Id       string
	Group    string
	Attempt  int32
	Sequence int64
	Source   string
}

type mockMessageRuntime struct{}

func (mockMessageRuntime) FlowToRuntimeMetadata(m flow.Metadata) mockMeta {
	return mockMeta{
		Id:       m.Id,
		Group:    m.Group,
		Attempt:  m.Attempt,
		Sequence: m.Sequence,
		Source:   m.Source,
	}
}

func (mockMessageRuntime) RuntimeToFlowMetadata(m mockMeta) flow.Metadata {
	return flow.Metadata{
		Id:       m.Id,
		Group:    m.Group,
		Attempt:  m.Attempt,
		Sequence: m.Sequence,
		Source:   m.Source,
	}
}

// mockMessageProducer implements message.Producer[M]. It records a copy of
// the slice passed to each Produce call and can return a configurable error.
type mockMessageProducer[M any] struct {
	produced [][]message.Message[M]
	err      error
}

func newMockMessageProducer[M any](err error) *mockMessageProducer[M] {
	return &mockMessageProducer[M]{err: err}
}

func (p *mockMessageProducer[M]) Produce(_ context.Context, messages []message.Message[M]) error {
	call := make([]message.Message[M], len(messages))
	copy(call, messages)
	p.produced = append(p.produced, call)
	return p.err
}

func (p *mockMessageProducer[M]) Start() error { return nil }

func (p *mockMessageProducer[M]) Stop() {}

// mockFlowHandler implements flow.Handler[V] and captures the messages it
// receives.
type mockFlowHandler[V any] struct {
	messages []flow.Message[V]
	err      error
}

func (h *mockFlowHandler[V]) Handle(_ context.Context, msg flow.Message[V]) error {
	h.messages = append(h.messages, msg)
	return h.err
}

// mockRawStateStore implements state.Store, backed by a map.
type mockRawStateStore struct {
	states     map[string]state.State
	readErr    error
	writeErr   error
	keysErr    error
	startErr   error
	readCalls  int
	writeCalls int
	keysCalls  int
	startCalls int
	stopCalls  int
}

func newMockRawStateStore() *mockRawStateStore {
	return &mockRawStateStore{states: map[string]state.State{}}
}

func (s *mockRawStateStore) Read(_ context.Context, id string) (state.State, error) {
	s.readCalls++
	if s.readErr != nil {
		return state.State{}, s.readErr
	}
	return s.states[id], nil
}

func (s *mockRawStateStore) Write(_ context.Context, entry state.State) error {
	s.writeCalls++
	if s.writeErr != nil {
		return s.writeErr
	}
	s.states[entry.Id] = entry
	return nil
}

func (s *mockRawStateStore) Delete(_ context.Context, _ string) error { return nil }

func (s *mockRawStateStore) Keys(_ context.Context) ([]string, error) {
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

func (s *mockRawStateStore) Start() error {
	s.startCalls++
	return s.startErr
}

func (s *mockRawStateStore) Stop() {
	s.stopCalls++
}

// myMessage is a local value type with JSON tags used with format.Json.
type myMessage struct {
	Id    string `json:"id"`
	Count int    `json:"count"`
}
