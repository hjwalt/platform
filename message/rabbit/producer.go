package rabbit

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/hjwalt/platform/message"
	"github.com/rabbitmq/amqp091-go"
)

func NewProducer(configuration RabbitConfiguration) message.Producer[RabbitMetadata] {
	return &RabbitProducer{
		Name:             configuration.Name,
		ConnectionString: configuration.ConnectionString,
	}
}

type RabbitProducer struct {
	Name             string
	ConnectionString string // amqp://guest:guest@localhost:5672/

	connection *amqp091.Connection
	channel    *amqp091.Channel
}

func (p *RabbitProducer) Start() error {
	config := amqp091.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Properties: amqp091.Table{
			"connection_name": p.Name,
		},
	}

	slog.Debug("rabbit producer starting")

	if conn, err := amqp091.DialConfig(p.ConnectionString, config); err != nil {
		return errors.Join(err, ErrRabbitConnection)
	} else {
		p.connection = conn
	}

	if ch, err := p.connection.Channel(); err != nil {
		return errors.Join(err, ErrRabbitChannel)
	} else {
		p.channel = ch
	}

	if err := p.channel.Confirm(false); err != nil {
		return errors.Join(err, ErrRabbitConfirmMode)
	}

	slog.Debug("rabbit producer started")

	return nil
}

func (p *RabbitProducer) Stop() {
	slog.Debug("rabbit producer stoppping")

	p.channel.Close()
	p.connection.Close()

	slog.Debug("rabbit producer stopped")
}

func (p *RabbitProducer) Produce(c context.Context, tarr []message.Message[RabbitMetadata]) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	for _, t := range tarr {

		confirm, err := p.channel.PublishWithDeferredConfirmWithContext(ctx,
			"",               // exchange
			t.Metadata.Queue, // routing key
			false,            // mandatory
			false,            // immediate
			amqp091.Publishing{
				DeliveryMode: amqp091.Persistent,
				Headers:      t.Metadata.Headers,
				ContentType:  "application/octet-stream",
				Body:         t.Value,
			},
		)
		if err != nil {
			return errors.Join(err, ErrRabbitProduce)
		}

		published, err := confirm.WaitContext(ctx)
		if err != nil {
			return errors.Join(err, ErrRabbitProduce)
		}

		if !published {
			return errors.Join(err, ErrRabbitProduce)
		}

	}

	return nil
}
