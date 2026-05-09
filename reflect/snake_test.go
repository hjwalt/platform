package reflect_test

import (
	"testing"

	"github.com/hjwalt/platform/reflect"
	"github.com/stretchr/testify/assert"
)

func TestSnake(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("snake_case_this", reflect.ToLowerSnakeCase("snakeCaseThis"))
	assert.Equal("snake_case_this", reflect.ToLowerSnakeCase("SnakeCaseThis"))

	assert.Equal("SNAKE_CASE_THIS", reflect.ToUpperSnakeCase("snakeCaseThis"))
	assert.Equal("SNAKE_CASE_THIS", reflect.ToUpperSnakeCase("SnakeCaseThis"))

	// odd cases

	assert.Equal("s_snake_case_this", reflect.ToLowerSnakeCase("SSnakeCaseThis"))
}
