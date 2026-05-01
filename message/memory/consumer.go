package memory

import (
	"context"
	"time"

	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/runtime"
)

func NewConsumer(configuration MemoryConfiguration, handler message.Handler[MemoryMetadata]) message.Consumer[MemoryMetadata] {
	return runtime.NewLoop(&MemoryConsumer{
		Channel: configuration.Channel,
		Handler: handler,
	})
}

type MemoryConsumer struct {
	Channel <-chan message.Message[MemoryMetadata]
	Handler message.Handler[MemoryMetadata]
}

func (r *MemoryConsumer) Start() error {
	return nil
}

func (r *MemoryConsumer) Stop() {

}

func (r *MemoryConsumer) Loop(ctx context.Context, cancel context.CancelFunc) error {
	select {
	case msg := <-r.Channel:
		r.Handler.Handle(ctx, msg)
	case <-time.After(1 * time.Second):
		// slog.Info("memory consumer timeout, iterate")
	}
	return nil
}
