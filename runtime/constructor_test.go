package runtime_test

import (
	"strconv"
	"testing"

	"github.com/hjwalt/platform/runtime"
	"github.com/stretchr/testify/assert"
)

func TestConstructorForNoConfigs(t *testing.T) {
	assert := assert.New(t)

	defaultCalls := 0
	constructor := runtime.ConstructorFor(
		func() int {
			defaultCalls++
			return 42
		},
		func(v int) string { return strconv.Itoa(v) },
	)

	assert.Equal("42", constructor())
	assert.Equal("42", constructor())
	assert.Equal(2, defaultCalls, "defaultValue should be invoked once per construction")
}

func TestConstructorForConfigsAppliedInOrder(t *testing.T) {
	assert := assert.New(t)

	addOne := func(v int) int { return v + 1 }
	multiplyByTen := func(v int) int { return v * 10 }

	constructor := runtime.ConstructorFor(
		func() int { return 0 },
		func(v int) string { return strconv.Itoa(v) },
	)

	// (0 + 1) * 10 = 10
	assert.Equal("10", constructor(addOne, multiplyByTen))

	// Reversing the order changes the result: (0 * 10) + 1 = 1,
	// proving configurations apply in the order they are passed.
	assert.Equal("1", constructor(multiplyByTen, addOne))
}

func TestConstructorForCastingIsApplied(t *testing.T) {
	assert := assert.New(t)

	constructor := runtime.ConstructorFor(
		func() int { return 7 },
		func(v int) string { return "value:" + strconv.Itoa(v) },
	)

	assert.Equal("value:7", constructor())
	assert.Equal("value:17", constructor(func(v int) int { return v + 10 }))
}
