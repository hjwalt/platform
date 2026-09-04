package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	assert := assert.New(t)

	producer := memory.NewProducer(memory.MemoryConfiguration{
		Channel: make(chan message.Message[memory.MemoryMetadata], 1),
	})

	assert.NotNil(producer)
}

func TestProducerStartStop(t *testing.T) {
	assert := assert.New(t)

	producer := memory.NewProducer(memory.MemoryConfiguration{
		Channel: make(chan message.Message[memory.MemoryMetadata], 1),
	})

	// Start and Stop are no-ops returning nil / not panicking, even repeated.
	assert.NoError(producer.Start())
	producer.Stop()
	producer.Stop()
	assert.NoError(producer.Start())
	producer.Stop()
}

func TestProducerProduceSendsAllMessagesInOrder(t *testing.T) {
	assert := assert.New(t)

	channel := make(chan message.Message[memory.MemoryMetadata], 3)
	producer := memory.NewProducer(memory.MemoryConfiguration{Channel: channel})

	messages := []message.Message[memory.MemoryMetadata]{
		{
			Metadata:  memory.MemoryMetadata{Headers: map[string]string{"id": "1"}},
			Value:     []byte("first"),
			Timestamp: time.Now(),
		},
		{
			Metadata:  memory.MemoryMetadata{Headers: map[string]string{"id": "2"}},
			Value:     []byte("second"),
			Timestamp: time.Now(),
		},
		{
			Metadata:  memory.MemoryMetadata{Headers: map[string]string{"id": "3"}},
			Value:     []byte("third"),
			Timestamp: time.Now(),
		},
	}

	err := producer.Produce(context.Background(), messages)

	assert.NoError(err)
	assert.Equal(3, len(channel), "all messages should be sitting on the channel")

	for i, want := range messages {
		got := <-channel
		assert.Equal(want, got, "message at index %d was received out of order or altered", i)
	}
}

func TestProducerProduceEmptyList(t *testing.T) {
	assert := assert.New(t)

	channel := make(chan message.Message[memory.MemoryMetadata], 1)
	producer := memory.NewProducer(memory.MemoryConfiguration{Channel: channel})

	err := producer.Produce(context.Background(), []message.Message[memory.MemoryMetadata]{})

	assert.NoError(err)
	assert.Equal(0, len(channel))
}
