package message

import (
	"context"
	"time"

	"github.com/hjwalt/platform/commons/structure"
)

type Message[M any, V any] struct {
	Id        string
	Metadata  M
	Value     V
	Timestamp time.Time
}

type Handler[M any, V any] interface {
	Handle(context.Context, Message[M, V]) error
}

type Producer[M any] interface {
	Produce(context.Context, []Message[M, structure.Bytes]) error
	Start() error
	Stop()
}

type Consumer[M any] interface {
	Start() error
	Stop()
}
