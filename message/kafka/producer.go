package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hjwalt/platform/message"
)

func NewProducer(configuration KafkaProducerConfiguration) message.Producer[KafkaMetadata] {
	return &KafkaProducer{
		Brokers:  configuration.Brokers,
		ClientId: configuration.ClientId,
	}
}

type KafkaProducer struct {
	// required
	Brokers  string
	ClientId string

	// set in start
	producer *kafka.Producer
}

func (r *KafkaProducer) Start() error {
	slog.Debug("starting kafka producer")

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": r.Brokers,
		"client.id":         r.ClientId,
		"acks":              "all",
	}

	if producer, err := kafka.NewProducer(kafkaConfig); err != nil {
		return errors.Join(err, ErrKafkaProducerConnectFail)
	} else {
		r.producer = producer
	}

	slog.Debug("stopping kafka producer")
	return nil
}

func (r *KafkaProducer) Stop() {
	slog.Debug("stopping kafka producer")

	if remaining := r.producer.Flush(60000); remaining > 0 {
		slog.Error("producer flush timed out")
	}
	r.producer.Close()

	slog.Debug("stopped kafka producer")
}

func (r *KafkaProducer) Produce(c context.Context, sources []message.Message[KafkaMetadata]) error {
	if len(sources) == 0 {
		return nil
	}
	for _, original := range sources {
		delivery_chan := make(chan kafka.Event, 10000)

		headers := make([]kafka.Header, 0)
		for k, v := range original.Metadata.Headers {
			headers = append(headers, kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
		}

		kafkaMsg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &original.Metadata.Topic,
				Partition: kafka.PartitionAny,
			},
			Headers: headers,
			Key:     []byte(original.Metadata.Key),
			Value:   original.Value,
		}

		if err := r.producer.Produce(kafkaMsg, delivery_chan); err != nil {
			return errors.Join(err, ErrKafkaProducerFail)
		}

		e := <-delivery_chan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return errors.Join(m.TopicPartition.Error, ErrKafkaProducerFail)
		}
	}

	return nil
}
