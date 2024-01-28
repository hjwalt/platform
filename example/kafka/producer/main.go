package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message/kafka"
)

type TestMessage struct {
	Id      string
	Message string
	Value   int32
}

func main() {
	producer := converter.RuntimeToFlowProducer(
		kafka.NewProducer(
			kafka.KafkaConfiguration{
				Brokers:  "localhost:9092",
				ClientId: "test_consumer",
			},
		),
		converter.NewConverter(
			flow_runtime_kafka.New("test"),
			format.Json[TestMessage](),
		),
	)

	producer.Start()

	producer.Produce(context.Background(), []flow.Message[TestMessage]{
		{
			Metadata: flow.Metadata{
				Id:     uuid.New().String(),
				Group:  "hello",
				Source: "MANUAL",
			},
			Value: TestMessage{
				Id:      "first",
				Message: "hello world",
				Value:   1,
			},
		},
	})

	producer.Stop()
}
