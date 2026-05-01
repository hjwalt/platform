package stateless

import (
	"context"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/logger"
)

func NewExploder[IV any, OV any, ERR any](name string, handler Explode[IV, OV, ERR], metadataOperation flow.MetadataOperation, outputProducer flow.Producer[OV], errorProducer flow.Producer[ERR]) flow.Handler[IV] {
	return &Exploder[IV, OV, ERR]{
		Name:              name,
		HandlerFunction:   handler,
		MetadataOperation: metadataOperation,
		OutputProducer:    outputProducer,
		ErrorProducer:     errorProducer,
	}
}

type Exploder[IV any, OV any, ERR any] struct {
	Name              string
	HandlerFunction   Explode[IV, OV, ERR]
	MetadataOperation flow.MetadataOperation
	OutputProducer    flow.Producer[OV]
	ErrorProducer     flow.Producer[ERR]
}

func (r *Exploder[IV, OV, ERR]) Handle(parentCtx context.Context, msg flow.Message[IV]) error {
	ctx := logger.WithContext(parentCtx, "function", r.Name)

	result, error := r.HandlerFunction(ctx, msg.Value)

	if result.IsPresent() {
		for _, entry := range result.Get() {
			outputMessage := flow.Message[OV]{
				Metadata:  r.MetadataOperation.OnSuccess(msg.Metadata),
				Value:     entry,
				Timestamp: time.Now(),
			}
			return r.OutputProducer.Produce(ctx, []flow.Message[OV]{outputMessage})
		}
	}
	if error.IsPresent() {
		errorMessage := flow.Message[ERR]{
			Metadata:  r.MetadataOperation.OnError(msg.Metadata),
			Value:     error.Get(),
			Timestamp: time.Now(),
		}
		return r.ErrorProducer.Produce(ctx, []flow.Message[ERR]{errorMessage})
	}
	return nil
}
