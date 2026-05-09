package stateless

import (
	"context"
	"time"

	"github.com/hjwalt/platform/flow"
)

func NewConsumer[IV any, ERR any](
	handler Consume[IV, ERR],
	errorMetadata flow.ExtractMetadata[ERR],
	errorProducer flow.Producer[ERR],
) flow.Handler[IV] {
	return &Consumer[IV, ERR]{
		HandlerFunction: handler,
		ErrorMetadata:   errorMetadata,
		ErrorProducer:   errorProducer,
	}
}

type Consumer[IV any, ERR any] struct {
	HandlerFunction Consume[IV, ERR]
	ErrorMetadata   flow.ExtractMetadata[ERR]
	ErrorProducer   flow.Producer[ERR]
}

func (r *Consumer[IV, ERR]) Handle(ctx context.Context, msg flow.Message[IV]) error {
	result := r.HandlerFunction(ctx, msg.Value)

	if result.IsPresent() {
		errorMessage := flow.Message[ERR]{
			Metadata:  r.ErrorMetadata(ctx, msg.Metadata, result.Get()),
			Value:     result.Get(),
			Timestamp: time.Now(),
		}
		return r.ErrorProducer.Produce(ctx, []flow.Message[ERR]{errorMessage})
	}

	return nil
}
