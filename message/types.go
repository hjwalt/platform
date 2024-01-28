package message

import (
	"context"
	"time"
)

type Message[M any] struct {
	Metadata  M
	Value     []byte
	Timestamp time.Time
}

type Handler[M any] interface {
	Handle(context.Context, Message[M]) error
}

type Producer[M any] interface {
	Produce(context.Context, []Message[M]) error
	Start() error
	Stop()
}

type Consumer[M any] interface {
	Start() error
	Stop()
}
