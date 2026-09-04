package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hjwalt/platform/message"
	"github.com/hjwalt/platform/message/memory"
	"github.com/stretchr/testify/assert"
)

var errHandlerFailure = errors.New("handler failure")

// recordingHandler records every handled message.
type recordingHandler struct {
	messages []message.Message[memory.MemoryMetadata]
}

func (r *recordingHandler) Handle(_ context.Context, m message.Message[memory.MemoryMetadata]) error {
	r.messages = append(r.messages, m)
	return nil
}

// failingHandler records the number of calls and always returns an error.
type failingHandler struct {
	err   error
	calls int
}

func (r *failingHandler) Handle(_ context.Context, _ message.Message[memory.MemoryMetadata]) error {
	r.calls++
	return r.err
}

func sampleMessage(id string) message.Message[memory.MemoryMetadata] {
	return message.Message[memory.MemoryMetadata]{
		Metadata:  memory.MemoryMetadata{Headers: map[string]string{"id": id}},
		Value:     []byte(id),
		Timestamp: time.Now(),
	}
}

func TestNewConsumerReturnsConsumer(t *testing.T) {
	assert := assert.New(t)

	consumer := memory.NewConsumer(
		memory.MemoryConfiguration{Channel: make(chan message.Message[memory.MemoryMetadata], 1)},
		&recordingHandler{},
	)

	assert.NotNil(consumer)
}

func TestMemoryConsumerStartStop(t *testing.T) {
	assert := assert.New(t)

	consumer := &memory.MemoryConsumer{
		Channel: make(chan message.Message[memory.MemoryMetadata], 1),
		Handler: &recordingHandler{},
	}

	// Start and Stop are no-ops.
	assert.NoError(consumer.Start())
	consumer.Stop()
	assert.NoError(consumer.Start())
	consumer.Stop()
}

func TestMemoryConsumerLoopHandlesSingleMessage(t *testing.T) {
	assert := assert.New(t)

	want := sampleMessage("1")
	channel := make(chan message.Message[memory.MemoryMetadata], 1)
	channel <- want

	handler := &recordingHandler{}
	consumer := &memory.MemoryConsumer{Channel: channel, Handler: handler}

	err := consumer.Loop(context.Background(), func() {})

	assert.NoError(err)
	assert.Equal(1, len(handler.messages))
	assert.Equal(want, handler.messages[0])
	assert.Equal(0, len(channel), "the message should have been consumed from the channel")
}

func TestMemoryConsumerLoopEmptyChannelWaitsThenReturns(t *testing.T) {
	assert := assert.New(t)

	consumer := &memory.MemoryConsumer{
		Channel: make(chan message.Message[memory.MemoryMetadata], 1),
		Handler: &recordingHandler{},
	}

	start := time.Now()
	err := consumer.Loop(context.Background(), func() {})
	elapsed := time.Since(start)

	// With no message available, Loop blocks on a 1 second timeout then
	// returns nil without calling the handler.
	assert.NoError(err)
	assert.GreaterOrEqual(elapsed, time.Second)
	assert.Less(elapsed, 5*time.Second)
	assert.Equal(0, len(consumer.Channel))
}

func TestMemoryConsumerLoopIgnoresHandlerError(t *testing.T) {
	assert := assert.New(t)

	channel := make(chan message.Message[memory.MemoryMetadata], 1)
	channel <- sampleMessage("2")

	handler := &failingHandler{err: errHandlerFailure}
	consumer := &memory.MemoryConsumer{Channel: channel, Handler: handler}

	err := consumer.Loop(context.Background(), func() {})

	// Oddity: the error returned by handler.Handle is discarded; Loop always
	// returns nil. Asserting actual behavior here.
	assert.NoError(err)
	assert.Equal(1, handler.calls)
}
