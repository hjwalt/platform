package main

import (
	"context"

	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/kafka"
)

func main() {

	producer := kafka.NewProducer(
		kafka.KafkaConfiguration{
			Brokers:  "localhost:9092",
			ClientId: "test_consumer",
		},
	)

	producer.Start()

	producer.Produce(context.Background(), []message.Message[kafka.KafkaMetadata, structure.Bytes]{
		{
			Metadata: kafka.KafkaMetadata{
				Topic: "test",
				Key:   "test",
			},
			Value: []byte("test"),
		},
	})

	producer.Stop()
}
