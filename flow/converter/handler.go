package converter

import (
	"context"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message"
)

func FlowToRuntimeHandler[M any, V any](handler flow.Handler[V], converter flow.Converter[M, V]) message.Handler[M] {
	return &FlowHandler[M, V]{
		Converter: converter,
		Handler:   handler,
	}
}

type FlowHandler[M any, V any] struct {
	Converter flow.Converter[M, V]
	Handler   flow.Handler[V]
}

func (r *FlowHandler[M, V]) Handle(parentCtx context.Context, msg message.Message[M]) error {
	convertedMsg, convertErr := r.Converter.RuntimeToFlow(msg)
	if convertErr != nil {
		return convertErr
	}
	messageContext := make(map[string]string)
	messageContext["message_id"] = convertedMsg.Metadata.Id
	messageContext["message_group"] = convertedMsg.Metadata.Group

	ctx := logger.WithContexts(parentCtx, messageContext)
	return r.Handler.Handle(ctx, convertedMsg)
}
