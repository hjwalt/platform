package kafka

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/runtime"
)

func NewConsumer(configuration KafkaConsumerConfiguration, handler message.Handler[KafkaMetadata]) message.Consumer[KafkaMetadata] {
	return runtime.NewLoop(&KafkaConsumer{
		Brokers:  configuration.Brokers,
		Topic:    configuration.Topic,
		ClientId: configuration.ClientId,
		GroupId:  configuration.GroupId,
		Handler:  handler,
	})
}

type KafkaConsumer struct {
	// required
	Brokers  string
	ClientId string
	GroupId  string
	Topic    string
	Handler  message.Handler[KafkaMetadata]

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
		"allow.auto.create.topics": "true",
		"max.poll.interval.ms":     "1800000",
		// allows using async commit and manual offset storing
		"enable.auto.commit":       "true",
		"enable.auto.offset.store": "false",
	}

	slog.Debug("starting kafka consumer")

	if consumer, err := kafka.NewConsumer(kafkaConfig); err != nil {
		return errors.Join(err, ErrKafkaConsumerConnectFail)
	} else {
		r.consumer = consumer
	}

	if err := r.consumer.SubscribeTopics([]string{r.Topic}, nil); err != nil {
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

		msg := message.Message[KafkaMetadata]{
			Metadata: KafkaMetadata{
				Topic:     *e.TopicPartition.Topic,
				Partition: e.TopicPartition.Partition,
				Offset:    int64(e.TopicPartition.Offset),
				Key:       string(e.Key),
				Headers:   headers,
			},
			Value: e.Value,
		}

		handlerCtx := logger.WithContext(ctx, "topic", *e.TopicPartition.Topic)
		if err := r.Handler.Handle(handlerCtx, msg); err != nil {
			return errors.Join(err, ErrKafkaConsumerConsume)
		}

		r.consumer.StoreMessage(e)
	case kafka.Error:
		slog.Error(e.Error())
		return errors.New(e.Error())
	default:
		// slog.Info("kafka consumer loop type not managed", "type", ev)
	}

	return nil
}
