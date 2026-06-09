package kafka_integration

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Configuration struct {
	Brokers  string
	ClientId string
	GroupId  string
	Custom   map[string]string
}

type Consumer = kafka.Consumer

type Producer = kafka.Producer
