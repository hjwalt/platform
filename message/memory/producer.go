package memory

import (
	"context"

	"github.com/hjwalt/platform/message"
)

func NewProducer(configuration MemoryConfiguration) message.Producer[MemoryMetadata] {
	return &MemoryProducer{
		Channel: configuration.Channel,
	}
}

type MemoryProducer struct {
	Channel chan<- message.Message[MemoryMetadata]
}

func (r *MemoryProducer) Start() error {
	return nil
}

func (r *MemoryProducer) Stop() {

}

func (r *MemoryProducer) Produce(c context.Context, sources []message.Message[MemoryMetadata]) error {
	for _, original := range sources {
		r.Channel <- original
	}
	return nil
}
