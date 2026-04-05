package example_word_collect

import (
	"context"
	"time"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/logger"
	"github.com/hjwalt/platform/commons/reflect"
	"github.com/hjwalt/platform/flows"
	"github.com/hjwalt/platform/flows/example"
	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/flows/runtime_sarama"
	"github.com/hjwalt/platform/flows/stateful"
	"github.com/hjwalt/platform/format"
	"go.uber.org/zap"
)

const (
	Instance = "flows-word-collect"
)

func key(ctx context.Context, m flow.Message[string, string]) (string, error) {
	return m.Key, nil
}

func aggregator(c context.Context, m flow.Message[string, string], s stateful.State[*example.WordCollectState]) (stateful.State[*example.WordCollectState], error) {
	logger.Info("collect aggregate")

	// setting defaults
	if s.Content == nil {
		s.Content = &example.WordCollectState{Count: 0}
	}

	// update state
	s.Content.Word = m.Key
	s.Content.Count += 1

	logger.Info("count", zap.Int64("count", s.Content.Count), zap.String("key", m.Key))

	return s, nil
}

func collector(c context.Context, persistenceId string, s stateful.State[*example.WordCollectState]) (*flow.Message[string, string], error) {

	logger.Info("collect")

	// create output message
	outMessage := flow.Message[string, string]{
		Topic: "word-count",
		Key:   s.Content.Word,
		Value: reflect.GetString(s.Content.Count),
	}

	return &outMessage, nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	flows.RegisterKafkaConsumerKeyedHandlerConfig(
		ci,
		runtime_sarama.WithKeyedHandlerMaxPerKey(1000),
		runtime_sarama.WithKeyedHandlerMaxBufferred(10000),
		runtime_sarama.WithKeyedHandlerMaxDelay(1*time.Second),
	)

	return flows.CollectorOneToOneConfiguration[*example.WordCollectState, string, string, string, string]{
		Name:             Instance,
		InputTopic:       flow.StringTopic("word"),
		OutputTopic:      flow.StringTopic("word-collect"),
		Aggregator:       aggregator,
		Collector:        collector,
		InputBroker:      "localhost:9092",
		OutputBroker:     "localhost:9092",
		HttpPort:         8081,
		StateFormat:      format.Protobuf[*example.WordCollectState](),
		StateKeyFunction: key,
	}
}

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
