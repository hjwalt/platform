package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/format"
	kafka_integration "github.com/hjwalt/platform/integration/kafka"
	"github.com/hjwalt/platform/message/kafka"
)

type TestMessage struct {
	Id      string
	Message string
	Value   int32
}

func main() {
	kafkaProducer := kafka.NewProducer(
		kafka_integration.Configuration{
			Brokers:  "localhost:9092",
			ClientId: "test_consumer",
		},
	)

	producer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New("test"),
			format.Json[TestMessage](),
		),
		metadata.IdUpdate,
	)

	kafkaProducer.Start()

	producer.ProduceMessage(context.Background(), []flow.Message[TestMessage]{
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

	kafkaProducer.Stop()
}
