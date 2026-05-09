package stateful

import (
	"context"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/logger"
)

func NewOperator[IV any, OV any, ST any, ERR any](
	name string,
	stateKey StateKey[IV],
	stateUpdate StateUpdate[IV, ST, ERR],
	handler Operate[IV, OV, ST, ERR],
	outputMetadata flow.ExtractMetadata[OV],
	errorMetadata flow.ExtractMetadata[ERR],
	outputProducer flow.Producer[OV],
	errorProducer flow.Producer[ERR],
	stateStore flow.Store[ST],
) flow.Handler[IV] {
	return &Operator[IV, OV, ST, ERR]{
		Name:            name,
		StateKey:        stateKey,
		StateUpdate:     stateUpdate,
		HandlerFunction: handler,
		OutputMetadata:  outputMetadata,
		ErrorMetadata:   errorMetadata,
		OutputProducer:  outputProducer,
		ErrorProducer:   errorProducer,
		StateStore:      stateStore,
	}
}

type Operator[IV any, OV any, ST any, ERR any] struct {
	Name            string
	StateKey        StateKey[IV]
	StateUpdate     StateUpdate[IV, ST, ERR]
	HandlerFunction Operate[IV, OV, ST, ERR]
	OutputMetadata  flow.ExtractMetadata[OV]
	ErrorMetadata   flow.ExtractMetadata[ERR]
	OutputProducer  flow.Producer[OV]
	ErrorProducer   flow.Producer[ERR]
	StateStore      flow.Store[ST]
}

func (r *Operator[IV, OV, ST, ERR]) Handle(parentCtx context.Context, msg flow.Message[IV]) error {
	ctx := logger.WithContext(parentCtx, "function", r.Name)

	key, keyErr := r.StateKey(ctx, msg.Value)
	if keyErr != nil {
		return keyErr
	}

	state, stateReadErr := r.StateStore.Read(ctx, key)
	if stateReadErr != nil {
		return stateReadErr
	}

	nextState := r.StateUpdate(ctx, msg.Value, state.Value)
	if nextState.IsRight() {
		errorMessage := flow.Message[ERR]{
			Metadata:  r.ErrorMetadata(ctx, msg.Metadata, nextState.Right()),
			Value:     nextState.Right(),
			Timestamp: time.Now(),
		}
		return r.ErrorProducer.Produce(ctx, []flow.Message[ERR]{errorMessage})
	}

	stateWriteErr := r.StateStore.Write(ctx, flow.State[ST]{
		Id:        key,
		Value:     nextState.Left(),
		Timestamp: time.Now(),
	})
	if stateWriteErr != nil {
		return stateWriteErr
	}

	result, handlerError := r.HandlerFunction(ctx, msg.Value, nextState.Left())
	if result.IsPresent() {
		outputMessage := flow.Message[OV]{
			Metadata:  r.OutputMetadata(ctx, msg.Metadata, result.Get()),
			Value:     result.Get(),
			Timestamp: time.Now(),
		}
		return r.OutputProducer.Produce(ctx, []flow.Message[OV]{outputMessage})
	}
	if handlerError.IsPresent() {
		errorMessage := flow.Message[ERR]{
			Metadata:  r.ErrorMetadata(ctx, msg.Metadata, handlerError.Get()),
			Value:     handlerError.Get(),
			Timestamp: time.Now(),
		}
		return r.ErrorProducer.Produce(ctx, []flow.Message[ERR]{errorMessage})
	}
	return nil
}
