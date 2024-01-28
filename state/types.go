package state

import (
	"context"
	"time"
)

type State struct {
	Id        string
	Value     []byte
	Timestamp time.Time
}

type Store interface {
	Read(context.Context, string) (State, error)
	Write(context.Context, State) error
	Start() error
	Stop()
}
