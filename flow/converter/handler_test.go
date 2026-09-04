package converter_test

import (
	"context"
	"testing"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message"
	"github.com/stretchr/testify/assert"
)

func TestFlowHandlerHandleSuccess(t *testing.T) {
	assert := assert.New(t)

	innerHandler := &mockFlowHandler[myMessage]{}
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	handler := converter.FlowToRuntimeHandler(innerHandler, conv)

	timestamp := time.Now()

	err := handler.Handle(context.Background(), message.Message[mockMeta]{
		Metadata: mockMeta{
			Id:       "mid-1",
			Group:    "grp-1",
			Attempt:  3,
			Sequence: 7,
			Source:   "src-1",
		},
		Value:     []byte(`{"id":"value-1","count":2}`),
		Timestamp: timestamp,
	})

	assert.NoError(err)
	assert.Len(innerHandler.messages, 1)

	received := innerHandler.messages[0]
	assert.Equal(flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 7,
		Source:   "src-1",
	}, received.Metadata)
	assert.Equal(myMessage{Id: "value-1", Count: 2}, received.Value)
	assert.Equal(timestamp, received.Timestamp)
}

func TestFlowHandlerHandleConvertErrorHandlerNotCalled(t *testing.T) {
	assert := assert.New(t)

	innerHandler := &mockFlowHandler[myMessage]{}
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	handler := converter.FlowToRuntimeHandler(innerHandler, conv)

	err := handler.Handle(context.Background(), message.Message[mockMeta]{
		Metadata: mockMeta{Id: "mid-1"},
		Value:    []byte("this is not valid json"),
	})

	assert.Error(err)
	assert.Empty(innerHandler.messages)
}

func TestFlowHandlerHandleInnerHandlerErrorPropagates(t *testing.T) {
	assert := assert.New(t)

	innerErr := testErr
	innerHandler := &mockFlowHandler[myMessage]{err: innerErr}
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	handler := converter.FlowToRuntimeHandler(innerHandler, conv)

	err := handler.Handle(context.Background(), message.Message[mockMeta]{
		Metadata: mockMeta{Id: "mid-1"},
		Value:    []byte(`{"id":"value-1","count":2}`),
	})

	assert.ErrorIs(err, innerErr)
	assert.Len(innerHandler.messages, 1)
}

func TestFlowHandlerHandleWithLoggerContext(t *testing.T) {
	assert := assert.New(t)

	// pre-populate the parent context the way logger.WithContext would, to
	// ensure the handler's context enrichment layers onto it without breaking
	innerHandler := &mockFlowHandler[myMessage]{}
	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	handler := converter.FlowToRuntimeHandler(innerHandler, conv)

	err := handler.Handle(context.Background(), message.Message[mockMeta]{
		Metadata: mockMeta{
			Id:       "mid-1",
			Group:    "grp-1",
			Attempt:  3,
			Sequence: 7,
			Source:   "src-1",
		},
		Value: []byte(`{"id":"value-1","count":2}`),
	})

	assert.NoError(err)
	assert.Len(innerHandler.messages, 1)
	assert.Equal("mid-1", innerHandler.messages[0].Metadata.Id)
}
