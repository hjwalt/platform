package kafka_integration

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Configuration struct {
	Brokers  string
	ClientId string
	GroupId  string
	Custom   map[string]string
}

type Consumer = kafka.Consumer

type Producer = kafka.Producer

func CreateConsumer(configuration Configuration) (*Consumer, error) {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers":        configuration.Brokers,
		"client.id":                configuration.ClientId,
		"group.id":                 configuration.GroupId,
		"auto.offset.reset":        "smallest",
		"allow.auto.create.topics": "true",
		"max.poll.interval.ms":     "1800000",
		// allows using async commit and manual offset storing
		"enable.auto.commit":       "true",
		"enable.auto.offset.store": "false",
	}

	for k, v := range configuration.Custom {
		kafkaConfig[k] = v
	}

	return kafka.NewConsumer(&kafkaConfig)
}

func CreateProducer(configuration Configuration) (*Producer, error) {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": configuration.Brokers,
		"client.id":         configuration.ClientId,
		"acks":              "all",
	}

	for k, v := range configuration.Custom {
		kafkaConfig[k] = v
	}

	return kafka.NewProducer(&kafkaConfig)
}
