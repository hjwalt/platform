package converter_test

import (
	"testing"
	"time"

	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message"
	"github.com/stretchr/testify/assert"
)

func TestConverterRuntimeToFlow(t *testing.T) {
	assert := assert.New(t)

	valueFormat := format.Json[myMessage]()
	conv := converter.NewConverter(mockMessageRuntime{}, valueFormat)

	timestamp := time.Now()

	flowMsg, err := conv.RuntimeToFlow(message.Message[mockMeta]{
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
	assert.Equal(flow.Metadata{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 7,
		Source:   "src-1",
	}, flowMsg.Metadata)
	assert.Equal(myMessage{Id: "value-1", Count: 2}, flowMsg.Value)
	assert.Equal(timestamp, flowMsg.Timestamp)
}

func TestConverterRuntimeToFlowUnmarshalError(t *testing.T) {
	assert := assert.New(t)

	valueFormat := format.Json[myMessage]()
	conv := converter.NewConverter(mockMessageRuntime{}, valueFormat)

	_, err := conv.RuntimeToFlow(message.Message[mockMeta]{
		Metadata: mockMeta{Id: "mid-1"},
		Value:    []byte("this is not valid json"),
	})

	assert.Error(err)
}

func TestConverterFlowToRuntime(t *testing.T) {
	assert := assert.New(t)

	valueFormat := format.Json[myMessage]()
	conv := converter.NewConverter(mockMessageRuntime{}, valueFormat)

	timestamp := time.Now()

	runtimeMsg, err := conv.FlowToRuntime(flow.Message[myMessage]{
		Metadata: flow.Metadata{
			Id:       "mid-1",
			Group:    "grp-1",
			Attempt:  3,
			Sequence: 7,
			Source:   "src-1",
		},
		Value:     myMessage{Id: "value-1", Count: 2},
		Timestamp: timestamp,
	})

	assert.NoError(err)
	assert.Equal(mockMeta{
		Id:       "mid-1",
		Group:    "grp-1",
		Attempt:  3,
		Sequence: 7,
		Source:   "src-1",
	}, runtimeMsg.Metadata)
	assert.Equal(timestamp, runtimeMsg.Timestamp)

	// bytes roundtrip through the value format
	unmarshalled, err := valueFormat.Unmarshal(runtimeMsg.Value)
	assert.NoError(err)
	assert.Equal(myMessage{Id: "value-1", Count: 2}, unmarshalled)
}

func TestConverterFlowToRuntimeMarshalError(t *testing.T) {
	assert := assert.New(t)

	conv := converter.NewConverter(mockMessageRuntime{}, format.Broken())

	_, err := conv.FlowToRuntime(flow.Message[string]{
		Metadata: flow.Metadata{Id: "mid-1"},
		Value:    "marshal",
	})

	assert.ErrorIs(err, format.ErrMarshal)
}

func TestConverterRuntimeToFlowRoundTrip(t *testing.T) {
	assert := assert.New(t)

	conv := converter.NewConverter(mockMessageRuntime{}, format.Json[myMessage]())
	timestamp := time.Now()

	original := message.Message[mockMeta]{
		Metadata: mockMeta{
			Id:       "mid-1",
			Group:    "grp-1",
			Attempt:  3,
			Sequence: 7,
			Source:   "src-1",
		},
		Value:     []byte(`{"id":"value-1","count":2}`),
		Timestamp: timestamp,
	}

	flowMsg, err := conv.RuntimeToFlow(original)
	assert.NoError(err)

	back, err := conv.FlowToRuntime(flowMsg)
	assert.NoError(err)
	assert.Equal(original, back)
}
