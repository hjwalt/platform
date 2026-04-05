package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/kafka"
)

type TestHandler struct {
}

func (r *TestHandler) Handle(ctx context.Context, msg message.Message[kafka.KafkaMetadata, structure.Bytes]) error {
	slog.Info("new message", "message", msg)
	return nil
}

func main() {
	handler := &TestHandler{}

	kafkaRuntime := kafka.NewConsumer(
		kafka.KafkaConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "test_consumer",
			Consumer: kafka.KafkaConsumerConfiguration{
				Topics:  []string{"test"},
				GroupId: "test",
			},
		},
		handler,
	)

	runtime.Start(
		[]runtime.Runtime{
			kafkaRuntime,
		},
		100*time.Millisecond,
	)

	runtime.Wait()
}
