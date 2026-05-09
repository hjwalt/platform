package flow

import (
	"context"
	"time"

	"github.com/hjwalt/platform/message"
)

type Metadata struct {
	Id       string
	Group    string
	Attempt  int32
	Sequence int64
	Source   string
}

type Message[V any] struct {
	Metadata  Metadata
	Value     V
	Timestamp time.Time
}

type State[V any] struct {
	Id        string
	Value     V
	Timestamp time.Time
}

type Handler[V any] interface {
	Handle(context.Context, Message[V]) error
}

type Producer[V any] interface {
	Produce(context.Context, []Message[V]) error
	Start() error
	Stop()
}

type Consumer[V any] interface {
	Start() error
	Stop()
}

type Converter[M any, V any] interface {
	RuntimeToFlow(msg message.Message[M]) (Message[V], error)
	FlowToRuntime(msg Message[V]) (message.Message[M], error)
}

type MessageRuntime[M any] interface {
	FlowToRuntimeMetadata(Metadata) M
	RuntimeToFlowMetadata(M) Metadata
}

type Store[V any] interface {
	Read(ctx context.Context, id string) (State[V], error)
	Write(ctx context.Context, s State[V]) error
	Start() error
	Stop()
}

type ExtractMetadata[V any] func(ctx context.Context, prev Metadata, value V) Metadata
