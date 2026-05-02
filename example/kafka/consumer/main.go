package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message/kafka"
	"github.com/hjwalt/platform/runtime"
)

type TestMessage struct {
	Message string
}

type TestHandler struct {
}

func (r *TestHandler) Handle(ctx context.Context, msg flow.Message[TestMessage]) error {
	stringify := format.Json[flow.Message[TestMessage]]()
	msgJson, _ := stringify.Marshal(msg)

	slog.Info("new message", "message", string(msgJson))
	return nil
}

func main() {
	handler := converter.FlowToRuntimeHandler(
		&TestHandler{},
		converter.NewConverter(
			flow_runtime_kafka.New("test"),
			format.Json[TestMessage](),
		),
	)
	consumer := kafka.NewConsumer(
		kafka.KafkaConsumerConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "test_consumer",
			Topic:    "test",
			GroupId:  "test",
		},
		handler,
	)

	runtime.Start(
		[]runtime.Runtime{
			consumer,
		},
		100*time.Millisecond,
	)

	runtime.Wait()
}
