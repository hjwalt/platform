package state

import (
	"context"
	"errors"
	"time"

	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/format"
)

type FormattedStore[V any] struct {
	Format format.Format[V]
	Store  Store
}

func (r *FormattedStore[V]) Read(ctx context.Context, id string) (State[V], error) {
	s, err := r.Store.Read(ctx, id)
	if err != nil {
		return State[V]{}, err
	}

	unmarshalled, unmarshalErr := r.Format.Unmarshal(s.Value)
	if unmarshalErr != nil {
		return State[V]{}, errors.Join(unmarshalErr, ErrWriteUnmarshal)
	}

	return State[V]{
		Id:        id,
		Value:     unmarshalled,
		Timestamp: time.Now(),
	}, nil
}

func (r *FormattedStore[V]) Write(ctx context.Context, s State[V]) error {

	bytes, marshalErr := r.Format.Marshal(s.Value)
	if marshalErr != nil {
		return errors.Join(ErrWriteMarshal, marshalErr)
	}

	return r.Store.Write(ctx, State[structure.Bytes]{
		Id:        s.Id,
		Value:     bytes,
		Timestamp: s.Timestamp,
	})
}

func (r *FormattedStore[V]) Start() error {
	return r.Store.Start()
}

func (r *FormattedStore[V]) Stop() {
	r.Store.Stop()
}

var (
	ErrWriteMarshal   = errors.New("cannot marshal value")
	ErrWriteUnmarshal = errors.New("cannot unmarshal value")
)
