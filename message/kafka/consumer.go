package kafka

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/message"
)

func NewConsumer(configuration KafkaConfiguration, handler message.Handler[KafkaMetadata, structure.Bytes]) message.Consumer[KafkaMetadata] {
	return runtime.NewLoop(&KafkaConsumer{
		Brokers:  configuration.Brokers,
		Topics:   configuration.Consumer.Topics,
		ClientId: configuration.ClientId,
		GroupId:  configuration.Consumer.GroupId,
		Handler:  handler,
	})
}

type KafkaConsumer struct {
	// required
	Brokers  string
	ClientId string
	GroupId  string
	Topics   []string
	Handler  message.Handler[KafkaMetadata, structure.Bytes]

	consumer *kafka.Consumer
	mu       sync.Mutex
}

func (r *KafkaConsumer) Start() error {
	// basic validations
	if r == nil {
		return ErrKafkaConsumerNil
	}

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":        r.Brokers,
		"client.id":                r.ClientId,
		"group.id":                 r.GroupId,
		"auto.offset.reset":        "smallest",
		"enable.auto.commit":       "false",
		"allow.auto.create.topics": "true",
	}

	slog.Debug("starting kafka consumer")

	if consumer, err := kafka.NewConsumer(kafkaConfig); err != nil {
		return errors.Join(err, ErrKafkaConsumerConnectFail)
	} else {
		r.consumer = consumer
	}

	if err := r.consumer.SubscribeTopics(r.Topics, nil); err != nil {
		return errors.Join(err, ErrKafkaConsumerSubscribeFail)
	}

	slog.Debug("started kafka consumer")

	return nil
}

func (r *KafkaConsumer) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	slog.Debug("stopping kafka consumer")
	if err := r.consumer.Unsubscribe(); err != nil {
		slog.Error("failed kafka consumer unsubscribe", "error", err)
	}
	slog.Debug("unsubscribed")
	if err := r.consumer.Close(); err != nil {
		slog.Error("failed kafka consumer close", "error", err)
	}
	slog.Debug("stopped kafka consumer")
}

func (r *KafkaConsumer) Loop(ctx context.Context, cancel context.CancelFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		slog.Info("kafka consumer exitting loop")
		return nil
	default:
	}

	ev := r.consumer.Poll(100)
	switch e := ev.(type) {
	case *kafka.Message:
		headers := make(map[string]string)

		for _, header := range e.Headers {
			headers[header.Key] = string(header.Value)
		}

		msg := message.Message[KafkaMetadata, structure.Bytes]{
			Metadata: KafkaMetadata{
				Topic:     *e.TopicPartition.Topic,
				Partition: e.TopicPartition.Partition,
				Offset:    int64(e.TopicPartition.Offset),
				Key:       string(e.Key),
				Headers:   headers,
			},
			Value: e.Value,
		}

		if err := r.Handler.Handle(ctx, msg); err != nil {
			return errors.Join(err, ErrKafkaConsumerConsume)
		}

		r.consumer.Commit()
	case kafka.Error:
		return errors.New(e.Error())
	default:
		// slog.Warn("kafka consumer loop type not managed")
	}

	return nil
}
