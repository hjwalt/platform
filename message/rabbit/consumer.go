package rabbit

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/runtime"
	"github.com/rabbitmq/amqp091-go"
)

func NewConsumer(configuration RabbitConfiguration, handler message.Handler[RabbitMetadata]) message.Consumer[RabbitMetadata] {
	return runtime.NewLoop(&RabbitConsumer{
		Name:             configuration.Name,
		ConnectionString: configuration.ConnectionString,
		QueueName:        configuration.Consumer.QueueName,
		QueueDurable:     configuration.Consumer.QueueDurable,
		Handler:          handler,
	})
}

type RabbitConsumer struct {
	Name             string
	ConnectionString string // amqp://guest:guest@localhost:5672/
	QueueName        string
	QueueDurable     bool
	Handler          message.Handler[RabbitMetadata]

	connection *amqp091.Connection
	channel    *amqp091.Channel
	queue      *amqp091.Queue
	messages   <-chan amqp091.Delivery
}

func (r *RabbitConsumer) Start() error {
	config := amqp091.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Properties: amqp091.Table{
			"connection_name": r.Name,
		},
	}

	slog.Debug("rabbit consumer starting")

	if conn, err := amqp091.DialConfig(r.ConnectionString, config); err != nil {
		return errors.Join(err, ErrRabbitConnection)
	} else {
		r.connection = conn
	}

	if ch, err := r.connection.Channel(); err != nil {
		return errors.Join(err, ErrRabbitChannel)
	} else {
		r.channel = ch
	}

	if err := r.channel.Confirm(false); err != nil {
		return errors.Join(err, ErrRabbitConfirmMode)
	}

	if err := r.channel.Qos(1, 0, false); err != nil {
		return errors.Join(err, ErrRabbitPrefetch)
	}

	if q, err := r.channel.QueueDeclare(r.QueueName, r.QueueDurable, false, false, false, nil); err != nil {
		return errors.Join(err, ErrRabbitQueue)
	} else {
		r.queue = &q
	}

	if msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	); err != nil {
		return errors.Join(err, ErrRabbitChannel)
	} else {
		r.messages = msgs
	}

	slog.Debug("rabbit consumer started")

	return nil
}

func (r *RabbitConsumer) Stop() {
	slog.Debug("rabbit consumer stopping")
	r.channel.Close()
	r.connection.Close()
	slog.Debug("rabbit consumer stopped")
}

func (r *RabbitConsumer) Loop(ctx context.Context, cancel context.CancelFunc) error {
	m, hasMore := <-r.messages
	if !hasMore {
		cancel()
		return nil
	}

	t := message.Message[RabbitMetadata]{
		Metadata: RabbitMetadata{
			Queue:    r.queue.Name,
			Headers:  m.Headers,
			Original: m,
		},
		Value:     m.Body,
		Timestamp: m.Timestamp,
	}

	if err := r.Handler.Handle(ctx, t); err != nil {
		return errors.Join(err, ErrRabbitConsume)
	}

	m.Ack(false)

	return nil
}
