package stateless_test

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
