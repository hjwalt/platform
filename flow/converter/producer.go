package converter

import (
	"context"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/message"
)

func RuntimeToFlowProducer[M any, V any](
	producer message.Producer[M],
	converter flow.Converter[M, V],
	metadata flow.ExtractMetadata[V],
) flow.Producer[V] {
	return &FlowProducer[M, V]{
		Converter: converter,
		Producer:  producer,
		Metadata:  metadata,
	}
}

type FlowProducer[M any, V any] struct {
	Converter flow.Converter[M, V]
	Producer  message.Producer[M]
	Metadata  flow.ExtractMetadata[V]
}

func (r *FlowProducer[M, V]) Produce(ctx context.Context, msgs []V) error {
	timestamp := time.Now()
	messageWithMetadata := make([]flow.Message[V], len(msgs))
	for i := 0; i < len(msgs); i++ {
		messageWithMetadata[i] = flow.Message[V]{
			Metadata:  r.Metadata(ctx, metadata.Default(), msgs[i]),
			Value:     msgs[i],
			Timestamp: timestamp,
		}
	}
	return r.ProduceMessage(ctx, messageWithMetadata)
}

func (r *FlowProducer[M, V]) ProduceMessage(ctx context.Context, msgs []flow.Message[V]) error {
	convertedMessages := make([]message.Message[M], len(msgs))
	for i := 0; i < len(msgs); i++ {
		convertedMsg, convertErr := r.Converter.FlowToRuntime(msgs[i])
		if convertErr != nil {
			return convertErr
		}
		convertedMessages[i] = convertedMsg
	}
	return r.Producer.Produce(ctx, convertedMessages)
}

func (r *FlowProducer[M, V]) Start() error {
	return nil
}

func (r *FlowProducer[M, V]) Stop() {
}
