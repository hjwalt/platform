package state

import (
	"context"
	"time"

	"github.com/hjwalt/platform/commons/structure"
)

type State[V any] struct {
	Id        string
	Value     V
	Timestamp time.Time
}

type Store interface {
	Read(context.Context, string) (State[structure.Bytes], error)
	Write(context.Context, State[structure.Bytes]) error
	Start() error
	Stop()
}
