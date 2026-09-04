package either_test

import (
	"testing"

	"github.com/hjwalt/platform/type/either"
	"github.com/stretchr/testify/assert"
)

func TestEitherLeft(t *testing.T) {
	assert := assert.New(t)

	e := either.Left[string, int]("left-value")

	assert.True(e.IsLeft())
	assert.False(e.IsRight())
	assert.Equal("left-value", e.Left())
	// the unset side returns the zero value
	assert.Equal(0, e.Right())
}

func TestEitherRight(t *testing.T) {
	assert := assert.New(t)

	e := either.Right[string, int](42)

	assert.True(e.IsRight())
	assert.False(e.IsLeft())
	assert.Equal(42, e.Right())
	// the unset side returns the zero value
	assert.Equal("", e.Left())
}

func TestEitherConstructors(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		name        string
		value       either.Either[string, int]
		wantIsLeft  bool
		wantIsRight bool
		wantLeft    string
		wantRight   int
	}{
		{
			name:        "left side set",
			value:       either.Left[string, int]("alpha"),
			wantIsLeft:  true,
			wantIsRight: false,
			wantLeft:    "alpha",
			wantRight:   0,
		},
		{
			name:        "right side set",
			value:       either.Right[string, int](17),
			wantIsLeft:  false,
			wantIsRight: true,
			wantLeft:    "",
			wantRight:   17,
		},
		{
			name:        "empty string left value still present",
			value:       either.Left[string, int](""),
			wantIsLeft:  true,
			wantIsRight: false,
			wantLeft:    "",
			wantRight:   0,
		},
	}

	for _, c := range cases {
		assert.Equal(c.wantIsLeft, c.value.IsLeft(), c.name)
		assert.Equal(c.wantIsRight, c.value.IsRight(), c.name)
		assert.Equal(c.wantLeft, c.value.Left(), c.name)
		assert.Equal(c.wantRight, c.value.Right(), c.name)
	}
}

func TestEitherZeroValueIsNotUsable(t *testing.T) {
	assert := assert.New(t)

	// the zero-value Either holds nil Optional interfaces internally,
	// so calling any method on it panics rather than returning zero values
	var e either.Either[string, int]

	assert.Panics(func() { e.IsLeft() })
	assert.Panics(func() { e.IsRight() })
	assert.Panics(func() { e.Left() })
	assert.Panics(func() { e.Right() })
}

func TestEitherPresentZeroValueIsDistinguishable(t *testing.T) {
	assert := assert.New(t)

	// a zero value that is explicitly set on the left side is still present
	e := either.Left[int, int](0)

	assert.True(e.IsLeft())
	assert.False(e.IsRight())
	assert.Equal(0, e.Left())
	assert.Equal(0, e.Right())
}

type eitherSample struct {
	Text  string
	Count int
}

func TestEitherStructValue(t *testing.T) {
	assert := assert.New(t)

	value := eitherSample{Text: "payload", Count: 3}

	e := either.Right[eitherSample, eitherSample](value)

	assert.True(e.IsRight())
	assert.False(e.IsLeft())
	assert.Equal(value, e.Right())
	// the unset side returns the zero value of the struct type
	assert.Equal(eitherSample{}, e.Left())
}
