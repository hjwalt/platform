package format_test

import (
	"testing"

	"github.com/hjwalt/platform/format"
	"github.com/stretchr/testify/assert"
)

type TestJsonStruct struct {
	Name  string
	Value int64
}

func TestJsonFormat(t *testing.T) {
	assert := assert.New(t)

	f := format.Json[TestJsonStruct]()
	v := TestJsonStruct{
		Name:  "test",
		Value: 1234567890,
	}
	b := []byte(`{"Name":"test","Value":1234567890}`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
}

func TestJsonFormatPointer(t *testing.T) {
	assert := assert.New(t)

	f := format.Json[*TestJsonStruct]()
	v := &TestJsonStruct{
		Name:  "test",
		Value: 1234567890,
	}
	b := []byte(`{"Name":"test","Value":1234567890}`)

	vb, em := f.Marshal(v)
	assert.NoError(em)

	bv, eu := f.Unmarshal(b)
	assert.NoError(eu)

	assert.Equal(b, vb)
	assert.Equal(v, bv)
}

func TestJsonFormatEmptyValue(t *testing.T) {
	assert := assert.New(t)

	f := format.Json[*TestJsonStruct]()

	vb, em := f.Marshal(nil)
	assert.NoError(em)
	assert.Equal([]byte{}, vb)

	bv, eu := f.Unmarshal(nil)
	assert.NoError(eu)
	assert.Equal(&TestJsonStruct{}, bv)
}
