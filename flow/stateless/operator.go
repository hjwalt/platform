package stateless

import (
	"context"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/logger"
)

func NewOperator[IV any, OV any, ERR any](
	name string,
	handler Operate[IV, OV, ERR],
	outputMetadata flow.ExtractMetadata[OV],
	errorMetadata flow.ExtractMetadata[ERR],
	outputProducer flow.Producer[OV],
	errorProducer flow.Producer[ERR],
) flow.Handler[IV] {
	return &Operator[IV, OV, ERR]{
		Name:            name,
		HandlerFunction: handler,
		OutputMetadata:  outputMetadata,
		ErrorMetadata:   errorMetadata,
		OutputProducer:  outputProducer,
		ErrorProducer:   errorProducer,
	}
}

type Operator[IV any, OV any, ERR any] struct {
	Name            string
	HandlerFunction Operate[IV, OV, ERR]
	OutputMetadata  flow.ExtractMetadata[OV]
	ErrorMetadata   flow.ExtractMetadata[ERR]
	OutputProducer  flow.Producer[OV]
	ErrorProducer   flow.Producer[ERR]
}

func (r *Operator[IV, OV, ERR]) Handle(parentCtx context.Context, msg flow.Message[IV]) error {
	ctx := logger.WithContext(parentCtx, "function", r.Name)

	result, handlerError := r.HandlerFunction(ctx, msg.Value)
	if result.IsPresent() {
		outputMessage := flow.Message[OV]{
			Metadata:  r.OutputMetadata(ctx, msg.Metadata, result.Get()),
			Value:     result.Get(),
			Timestamp: time.Now(),
		}
		return r.OutputProducer.ProduceMessage(ctx, []flow.Message[OV]{outputMessage})
	}
	if handlerError.IsPresent() {
		errorMessage := flow.Message[ERR]{
			Metadata:  r.ErrorMetadata(ctx, msg.Metadata, handlerError.Get()),
			Value:     handlerError.Get(),
			Timestamp: time.Now(),
		}
		return r.ErrorProducer.ProduceMessage(ctx, []flow.Message[ERR]{errorMessage})
	}
	return nil
}
