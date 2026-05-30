package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/flow/stateful"
	"github.com/hjwalt/platform/flow/stateless"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/hjwalt/platform/runtime"
	file_store "github.com/hjwalt/platform/state/file"
	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
)

type TestState struct {
	Total int32
}

type TestMessage struct {
	Id      string
	Message string
	Value   int32
}

type Completed struct {
	Id      string
	Message string
}

type TestError struct {
	Id      string
	Message string
}

type Metric struct {
	Id    string
	Total int32
}

func KeyFunction(ctx context.Context, in TestMessage) (string, error) {
	return in.Id, nil
}

func StateUpdate(ctx context.Context, in TestMessage, st TestState) either.Either[TestState, TestError] {
	return either.Left[TestState, TestError](TestState{
		Total: st.Total + 1,
	})
}

func Accumulate(ctx context.Context, in TestMessage, st TestState) (optional.Optional[Metric], optional.Optional[TestError]) {
	return optional.Of(Metric{
		Id:    in.Id,
		Total: st.Total,
	}), optional.Empty[TestError]()
}

func Increment(ctx context.Context, in TestMessage) (optional.Optional[TestMessage], optional.Optional[TestError]) {
	if in.Value <= 10 {
		slog.InfoContext(ctx, "increment")
		return optional.Of(TestMessage{
			Id:      in.Id,
			Message: in.Message,
			Value:   in.Value + 1,
		}), optional.Empty[TestError]()
	}

	return optional.Empty[TestMessage](), optional.Empty[TestError]()
}

func Complete(ctx context.Context, in TestMessage) (optional.Optional[Completed], optional.Optional[TestError]) {
	if in.Value > 10 {
		slog.InfoContext(ctx, "larger than 10")
		return optional.Of(Completed{
			Id:      in.Id,
			Message: "larger than ten",
		}), optional.Empty[TestError]()
	}

	return optional.Empty[Completed](), optional.Empty[TestError]()
}

func LogCompleted(ctx context.Context, in Completed) optional.Optional[TestError] {
	slog.InfoContext(ctx, "completed message", "message", in.Message)
	return optional.Empty[TestError]()
}

func LogError(ctx context.Context, in TestError) optional.Optional[TestError] {
	slog.InfoContext(ctx, "error message", "message", in.Message)
	return optional.Empty[TestError]()
}

func LogMetric(ctx context.Context, in Metric) optional.Optional[TestError] {
	slog.InfoContext(ctx, "consumed count", "count", in.Total)
	return optional.Empty[TestError]()
}

func main() {
	kafkaProducer := kafka.NewProducer(
		kafka.KafkaProducerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "output_producer",
		},
	)

	outputProducer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New("test"),
			format.Json[TestMessage](),
		),
		metadata.IdUpdate,
	)

	completedProducer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New("completed"),
			format.Json[Completed](),
		),
		metadata.IdUpdate,
	)

	metricProducer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New("metric"),
			format.Json[Metric](),
		),
		metadata.IdUpdate,
	)

	errorProducer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New("error"),
			format.Json[TestError](),
		),
		metadata.IdUpdate,
	)

	incrementConsumer := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "increment_consumer",
			Topic:    "test",
			GroupId:  "increment_consumer",
		},
		converter.FlowToRuntimeHandler(
			stateless.NewOperator(
				"Increment",
				Increment,
				metadata.IdUpdate,
				metadata.AttemptIncrement,
				outputProducer,
				errorProducer,
			),
			converter.NewConverter(
				flow_runtime_kafka.New("test"),
				format.Json[TestMessage](),
			),
		),
	)

	completeNextConsumer := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "complete_consumer",
			Topic:    "test",
			GroupId:  "complete_consumer",
		},
		converter.FlowToRuntimeHandler(
			stateless.NewOperator(
				"Complete",
				Complete,
				metadata.IdUpdate,
				metadata.AttemptIncrement,
				completedProducer,
				errorProducer,
			),
			converter.NewConverter(
				flow_runtime_kafka.New("test"),
				format.Json[TestMessage](),
			),
		),
	)

	completedConsumer := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "completed_consumer",
			Topic:    "completed",
			GroupId:  "completed_consumer",
		},
		converter.FlowToRuntimeHandler(
			stateless.NewConsumer(
				LogCompleted,
				metadata.AttemptIncrement,
				errorProducer,
			),
			converter.NewConverter(
				flow_runtime_kafka.New("completed"),
				format.Json[Completed](),
			),
		),
	)

	errorConsumer := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "error_consumer",
			Topic:    "error",
			GroupId:  "error_consumer",
		},
		converter.FlowToRuntimeHandler(
			stateless.NewConsumer(
				LogError,
				metadata.AttemptIncrement,
				errorProducer,
			),
			converter.NewConverter(
				flow_runtime_kafka.New("test"),
				format.Json[TestError](),
			),
		),
	)

	metricConsumer := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "metric_consumer",
			Topic:    "metric",
			GroupId:  "metric_consumer",
		},
		converter.FlowToRuntimeHandler(
			stateless.NewConsumer(
				LogMetric,
				metadata.AttemptIncrement,
				errorProducer,
			),
			converter.NewConverter(
				flow_runtime_kafka.New("metric"),
				format.Json[Metric](),
			),
		),
	)

	fileStore := converter.RuntimeToFlowStore(
		file_store.New(file_store.Configuration{
			Path: "/home/hjwalt/Projects/platform/tmp/",
		}),
		format.Json[TestState](),
	)

	accumulator := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "accumulate_consumer",
			Topic:    "test",
			GroupId:  "accumulate_consumer",
		},
		converter.FlowToRuntimeHandler(
			stateful.NewOperator(
				"accumulate",
				KeyFunction,
				StateUpdate,
				Accumulate,
				metadata.IdUpdate,
				metadata.AttemptIncrement,
				metricProducer,
				errorProducer,
				fileStore,
			),
			converter.NewConverter(
				flow_runtime_kafka.New("test"),
				format.Json[TestMessage](),
			),
		),
	)

	runtime.Start(
		[]runtime.Runtime{
			kafkaProducer,
			fileStore,
			metricProducer,
			outputProducer,
			completedProducer,
			errorProducer,
			completeNextConsumer,
			incrementConsumer,
			completedConsumer,
			errorConsumer,
			accumulator,
			metricConsumer,
		},
		100*time.Millisecond,
	)

	runtime.Wait()
}
