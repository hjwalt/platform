package kafka

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	kafka_integration "github.com/hjwalt/platform/integration/kafka"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/runtime"
)

func NewConsumer(configuration kafka_integration.Configuration, topic string, handler message.Handler[KafkaMetadata]) message.Consumer[KafkaMetadata] {
	return runtime.NewLoop(&KafkaConsumer{
		configuration: configuration,
		topic:         topic,
		handler:       handler,
	})
}

type KafkaConsumer struct {
	// required
	configuration kafka_integration.Configuration
	topic         string
	handler       message.Handler[KafkaMetadata]

	consumer *kafka.Consumer
	mu       sync.Mutex
}

func (r *KafkaConsumer) Start() error {
	// basic validations
	if r == nil {
		return ErrKafkaConsumerNil
	}

	slog.Debug("starting kafka consumer")

	if consumer, err := kafka_integration.CreateConsumer(r.configuration); err != nil {
		return errors.Join(err, ErrKafkaConsumerConnectFail)
	} else {
		r.consumer = consumer
	}

	if err := r.consumer.SubscribeTopics([]string{r.topic}, nil); err != nil {
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
		if err := r.handler.Handle(handlerCtx, msg); err != nil {
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
