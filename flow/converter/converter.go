package converter

import (
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message"
)

func NewConverter[M any, V any](messageRuntime flow.MessageRuntime[M], valueFormat format.Format[V]) flow.Converter[M, V] {
	return &FlowConverter[M, V]{
		MessageRuntime: messageRuntime,
		Format:         valueFormat,
	}
}

type FlowConverter[M any, V any] struct {
	MessageRuntime flow.MessageRuntime[M]
	Format         format.Format[V]
}

func (r *FlowConverter[M, V]) RuntimeToFlow(msg message.Message[M]) (flow.Message[V], error) {
	value, unmarshalErr := r.Format.Unmarshal(msg.Value)
	if unmarshalErr != nil {
		return flow.Message[V]{}, unmarshalErr
	}

	meta := r.MessageRuntime.RuntimeToFlowMetadata(msg.Metadata)

	return flow.Message[V]{
		Metadata:  meta,
		Value:     value,
		Timestamp: msg.Timestamp,
	}, nil
}

func (r *FlowConverter[M, V]) FlowToRuntime(msg flow.Message[V]) (message.Message[M], error) {
	value, marshalErr := r.Format.Marshal(msg.Value)
	if marshalErr != nil {
		return message.Message[M]{}, marshalErr
	}

	meta := r.MessageRuntime.FlowToRuntimeMetadata(msg.Metadata)

	return message.Message[M]{
		Metadata:  meta,
		Value:     value,
		Timestamp: msg.Timestamp,
	}, nil
}
