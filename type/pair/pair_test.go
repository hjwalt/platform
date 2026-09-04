package pair_test

import (
	"testing"

	"github.com/hjwalt/platform/type/pair"
	"github.com/stretchr/testify/assert"
)

func TestPairOf(t *testing.T) {
	assert := assert.New(t)

	p := pair.Of("left", 42)

	assert.Equal("left", p.Left)
	assert.Equal(42, p.Right)
}

func TestPairOfVariousTypes(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		name      string
		value     pair.Pair[string, int]
		wantLeft  string
		wantRight int
	}{
		{
			name:      "words and numbers",
			value:     pair.Of[string, int]("alpha", 1),
			wantLeft:  "alpha",
			wantRight: 1,
		},
		{
			name:      "empty left and negative right",
			value:     pair.Of[string, int]("", -7),
			wantLeft:  "",
			wantRight: -7,
		},
	}

	for _, c := range cases {
		assert.Equal(c.wantLeft, c.value.Left, c.name)
		assert.Equal(c.wantRight, c.value.Right, c.name)
	}
}

func TestPairFieldTypes(t *testing.T) {
	assert := assert.New(t)

	p := pair.Of[string, bool]("flag", true)

	assert.Equal("flag", p.Left)
	assert.Equal(true, p.Right)
}

func TestPairFieldsAreAssignable(t *testing.T) {
	assert := assert.New(t)

	p := pair.Of("before", 1)

	p.Left = "after"
	p.Right = 2

	assert.Equal("after", p.Left)
	assert.Equal(2, p.Right)
}

func TestPairZeroValue(t *testing.T) {
	assert := assert.New(t)

	var p pair.Pair[string, int]

	assert.Equal("", p.Left)
	assert.Equal(0, p.Right)
}
