package metadata

import (
	"context"

	"github.com/google/uuid"
	"github.com/hjwalt/platform/flow"
)

func IdUpdate[V any](ctx context.Context, pref flow.Metadata, value V) flow.Metadata {
	return flow.Metadata{
		Id:       uuid.New().String(),
		Group:    pref.Group,
		Attempt:  0,
		Sequence: pref.Sequence,
		Source:   pref.Source,
	}
}

func AttemptIncrement[V any](ctx context.Context, pref flow.Metadata, value V) flow.Metadata {
	return flow.Metadata{
		Id:       uuid.New().String(),
		Group:    pref.Group,
		Attempt:  pref.Attempt + 1,
		Sequence: pref.Sequence,
		Source:   pref.Source,
	}
}
