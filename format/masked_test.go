package format_test

import (
	"testing"

	"github.com/hjwalt/platform/format"
	"github.com/stretchr/testify/assert"
)

func TestMaskedEncryption(t *testing.T) {
	assert := assert.New(t)

	mask := format.Cracked()
	actual := format.Broken()

	masked := format.Masked(mask, actual)

	maskedBytes, marshalErr := masked.Marshal("test")

	assert.NoError(marshalErr)

	gengarBytes, _ := actual.Marshal("test")

	assert.Equal("test", string(gengarBytes))
	assert.NotEqual("test", string(maskedBytes))

	unmaskedBytes, unmarshalErr := masked.Unmarshal(maskedBytes)

	assert.NoError(unmarshalErr)
	assert.Equal("test", unmaskedBytes)
}

func TestMaskedEncryptionMarshalErr(t *testing.T) {
	assert := assert.New(t)

	mask := format.Cracked()
	actual := format.Broken()

	masked := format.Masked(mask, actual)

	_, marshalErr := masked.Marshal("marshal")

	assert.ErrorIs(marshalErr, format.ErrMarshal)
	assert.ErrorIs(marshalErr, format.ErrMaskActualMarshal)

	_, maskErr := masked.Marshal("mask")

	assert.ErrorIs(maskErr, format.ErrMask)
	assert.ErrorIs(maskErr, format.ErrMaskMarshal)
}

func TestMaskedEncryptionUnmarshalErr(t *testing.T) {
	assert := assert.New(t)

	mask := format.Cracked()
	actual := format.Broken()

	masked := format.Masked(mask, actual)

	errInducingInput1, _ := mask.Mask([]byte("unmarshal"))
	errInducingInput2, _ := mask.Mask([]byte("unmask"))

	_, marshalErr := masked.Unmarshal(errInducingInput1)

	assert.ErrorIs(marshalErr, format.ErrUnmarshal)
	assert.ErrorIs(marshalErr, format.ErrMaskActualUnmarshal)

	_, maskErr := masked.Unmarshal(errInducingInput2)

	assert.ErrorIs(maskErr, format.ErrUnmask)
	assert.ErrorIs(maskErr, format.ErrMaskUnmarshal)
}
