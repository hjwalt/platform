package converter

import (
	"context"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/message"
)

func RuntimeToFlowProducer[M any, V any](producer message.Producer[M], converter flow.Converter[M, V]) flow.Producer[V] {
	return &FlowProducer[M, V]{
		Converter: converter,
		Producer:  producer,
	}
}

type FlowProducer[M any, V any] struct {
	Converter flow.Converter[M, V]
	Producer  message.Producer[M]
}

func (r *FlowProducer[M, V]) Produce(ctx context.Context, msgs []flow.Message[V]) error {
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
