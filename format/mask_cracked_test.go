package format_test

import (
	"testing"

	"github.com/hjwalt/platform/format"
	"github.com/stretchr/testify/assert"
)

func TestCracked(t *testing.T) {
	assert := assert.New(t)

	f := format.Cracked()

	v := []byte("test")

	vb, em := f.Mask(v)
	assert.NoError(em)
	assert.Equal("tset", string(vb))

	bv, eu := f.Unmask(vb)
	assert.NoError(eu)

	assert.Equal(v, bv)
}

func TestCrackedError(t *testing.T) {
	assert := assert.New(t)

	f := format.Cracked()

	var err error

	_, err = f.Mask([]byte("mask"))
	assert.ErrorIs(err, format.ErrMask)
	_, err = f.Mask([]byte("error"))
	assert.ErrorIs(err, format.ErrBasic)

	_, err = f.Unmask([]byte("ksamnu"))
	assert.ErrorIs(err, format.ErrUnmask)
	_, err = f.Unmask([]byte("rorre"))
	assert.ErrorIs(err, format.ErrBasic)
}
