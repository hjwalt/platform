package converter

import (
	"context"
	"errors"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/state"
)

func RuntimeToFlowStore[V any](store state.Store, format format.Format[V]) flow.Store[V] {
	return &FlowStore[V]{
		Format: format,
		Store:  store,
	}
}

type FlowStore[V any] struct {
	Format format.Format[V]
	Store  state.Store
}

func (r *FlowStore[V]) Read(ctx context.Context, id string) (flow.State[V], error) {
	s, err := r.Store.Read(ctx, id)
	if err != nil {
		return flow.State[V]{}, err
	}

	unmarshalled, unmarshalErr := r.Format.Unmarshal(s.Value)
	if unmarshalErr != nil {
		return flow.State[V]{}, errors.Join(unmarshalErr, ErrWriteUnmarshal)
	}

	return flow.State[V]{
		Id:        id,
		Value:     unmarshalled,
		Timestamp: time.Now(),
	}, nil
}

func (r *FlowStore[V]) Write(ctx context.Context, s flow.State[V]) error {

	bytes, marshalErr := r.Format.Marshal(s.Value)
	if marshalErr != nil {
		return errors.Join(ErrWriteMarshal, marshalErr)
	}

	return r.Store.Write(ctx, state.State{
		Id:        s.Id,
		Value:     bytes,
		Timestamp: s.Timestamp,
	})
}

func (r *FlowStore[V]) Keys(ctx context.Context) ([]string, error) {
	return r.Store.Keys(ctx)
}

func (r *FlowStore[V]) Start() error {
	return r.Store.Start()
}

func (r *FlowStore[V]) Stop() {
	r.Store.Stop()
}

var (
	ErrWriteMarshal   = errors.New("cannot marshal value")
	ErrWriteUnmarshal = errors.New("cannot unmarshal value")
)
