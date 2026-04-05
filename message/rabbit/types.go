package rabbit

import (
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMetadata struct {
	Queue    string
	Headers  map[string]any
	Original amqp091.Delivery
}

type RabbitConfiguration struct {
	Name             string
	ConnectionString string // amqp://guest:guest@localhost:5672/
	Consumer         RabbitConsumerConfiguration
}

type RabbitConsumerConfiguration struct {
	QueueName    string
	QueueDurable bool
}

var (
	ErrRabbitChannel             = errors.New("rabbitmq channel failed")
	ErrRabbitConfirmMode         = errors.New("rabbitmq channel confirm mode failed")
	ErrRabbitConnection          = errors.New("rabbitmq connection failed")
	ErrRabbitConsume             = errors.New("rabbit messages consume")
	ErrRabbitMessages            = errors.New("rabbit messages start consume")
	ErrRabbitPrefetch            = errors.New("rabbit prefetch setting")
	ErrRabbitProduce             = errors.New("rabbitmq producing")
	ErrRabbitProduceNotConfirmed = errors.New("rabbitmq produce not confirmed")
	ErrRabbitQueue               = errors.New("rabbitmq queue declaration failed")
)
