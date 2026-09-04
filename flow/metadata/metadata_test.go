package metadata_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/stretchr/testify/assert"
)

func TestMetadataHeaderKey(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("FLOW_METADATA", metadata.MetadataHeaderKey)
}

func TestDefault(t *testing.T) {
	assert := assert.New(t)

	d := metadata.Default()

	assert.Equal(int32(0), d.Attempt)
	assert.Equal(int64(-1), d.Sequence)
	assert.Equal("UNKNOWN", d.Source)
	assert.Equal("", d.Group)

	assert.NotEmpty(d.Id)
	_, err := uuid.Parse(d.Id)
	assert.NoError(err, "Default() Id should be a valid UUID")
}

func TestDefaultGeneratesFreshId(t *testing.T) {
	assert := assert.New(t)

	first := metadata.Default()
	second := metadata.Default()

	assert.NotEqual(first.Id, second.Id, "each Default() call should generate a fresh Id")
}

func TestIdUpdate(t *testing.T) {
	assert := assert.New(t)

	pref := flow.Metadata{
		Id:       "previous-id",
		Group:    "group-a",
		Attempt:  7,
		Sequence: 42,
		Source:   "source-x",
	}

	got := metadata.IdUpdate(context.Background(), pref, "value-payload")

	// Id is replaced by a fresh non-empty UUID
	assert.NotEmpty(got.Id)
	assert.NotEqual(pref.Id, got.Id)
	_, err := uuid.Parse(got.Id)
	assert.NoError(err, "IdUpdate() Id should be a valid UUID")

	// Attempt is reset to zero
	assert.Equal(int32(0), got.Attempt)

	// Group, Sequence and Source are preserved
	assert.Equal(pref.Group, got.Group)
	assert.Equal(pref.Sequence, got.Sequence)
	assert.Equal(pref.Source, got.Source)
}

func TestIdUpdateGeneratesFreshIdPerCall(t *testing.T) {
	assert := assert.New(t)

	pref := flow.Metadata{Id: "previous-id", Attempt: 3, Sequence: 1, Source: "src"}

	first := metadata.IdUpdate(context.Background(), pref, "value-payload")
	second := metadata.IdUpdate(context.Background(), pref, "value-payload")

	assert.NotEqual(first.Id, second.Id, "each IdUpdate() call should generate a fresh Id")
}

func TestAttemptIncrement(t *testing.T) {
	assert := assert.New(t)

	pref := flow.Metadata{
		Id:       "previous-id",
		Group:    "group-b",
		Attempt:  7,
		Sequence: -5,
		Source:   "source-y",
	}

	got := metadata.AttemptIncrement(context.Background(), pref, "value-payload")

	// Id is replaced by a fresh non-empty UUID
	assert.NotEmpty(got.Id)
	assert.NotEqual(pref.Id, got.Id)
	_, err := uuid.Parse(got.Id)
	assert.NoError(err, "AttemptIncrement() Id should be a valid UUID")

	// Attempt is incremented by one
	assert.Equal(int32(8), got.Attempt)

	// Group, Sequence and Source are preserved
	assert.Equal(pref.Group, got.Group)
	assert.Equal(pref.Sequence, got.Sequence)
	assert.Equal(pref.Source, got.Source)
}

func TestAttemptIncrementFromDefault(t *testing.T) {
	assert := assert.New(t)

	got := metadata.AttemptIncrement(context.Background(), metadata.Default(), "value-payload")

	assert.Equal(int32(1), got.Attempt)
}

func TestMetadataFormatRoundTrip(t *testing.T) {
	assert := assert.New(t)

	f := metadata.Format

	m := flow.Metadata{
		Id:       uuid.New().String(),
		Group:    "group-c",
		Attempt:  3,
		Sequence: 99,
		Source:   "source-z",
	}

	bytes, err := f.Marshal(m)
	assert.NoError(err)
	assert.NotEmpty(bytes)

	roundTripped, err := f.Unmarshal(bytes)
	assert.NoError(err)

	assert.Equal(m, roundTripped)
}

func TestMetadataFormatRoundTripDefault(t *testing.T) {
	assert := assert.New(t)

	f := metadata.Format

	d := metadata.Default()

	bytes, err := f.Marshal(d)
	assert.NoError(err)

	roundTripped, err := f.Unmarshal(bytes)
	assert.NoError(err)

	assert.Equal(d, roundTripped)
}

func TestMetadataFormatUnmarshalEmptyGivesDefault(t *testing.T) {
	assert := assert.New(t)

	f := metadata.Format

	// empty bytes unmarshal to an empty (non-defaulted) metadata
	got, err := f.Unmarshal([]byte{})
	assert.NoError(err)

	assert.Equal(flow.Metadata{}, got)
}
