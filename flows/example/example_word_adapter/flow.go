package example_word_adapter

import (
	"context"

	"github.com/hjwalt/platform/commons/format"
	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/logger"
	"github.com/hjwalt/platform/flows"
	"github.com/hjwalt/platform/flows/adapter"
	"github.com/hjwalt/platform/flows/flow"
)

const (
	Instance = "flows-word-adapter"
)

func InstanceFunction(c context.Context, m flow.Message[[]byte, adapter.BasicResponse]) (*flow.Message[string, string], error) {
	logger.Info("posted")
	return &flow.Message[string, string]{
		Topic:   "word",
		Key:     m.Value.Message,
		Value:   m.Value.Message,
		Headers: m.Headers,
	}, nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.RouterAdapterConfiguration[adapter.BasicResponse, string, string]{
		Name:               Instance,
		Path:               "/adapter",
		ProduceTopic:       flow.StringTopic("word"),
		ProduceBroker:      "localhost:9092",
		RequestBodyFormat:  format.Json[adapter.BasicResponse](),
		RequestMapFunction: InstanceFunction,
		HttpPort:           8081,
	}
}

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
