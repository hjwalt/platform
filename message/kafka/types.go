package kafka

import "errors"

type KafkaMetadata struct {
	Topic     string
	Offset    int64
	Partition int32
	Key       string
	Headers   map[string]string
}

type KafkaProducerConfiguration struct {
	Brokers  string
	ClientId string
}

type KafkaConsumerConfiguration struct {
	Brokers  string
	ClientId string
	GroupId  string
	Topic    string
}

var (
	ErrKafkaConsumerNil           = errors.New("kafka consumer is nil")
	ErrKafkaConsumerConnectFail   = errors.New("kafka consumer connect failed")
	ErrKafkaConsumerSubscribeFail = errors.New("kafka consumer subscribe failed")
	ErrKafkaConsumerConsume       = errors.New("kafka consumer handler failed")
	ErrKafkaProducerConnectFail   = errors.New("kafka producer connect failed")
	ErrKafkaProducerFail          = errors.New("kafka producer failed to write")
)
