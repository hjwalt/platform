package agent_skill_test

import (
	"errors"
	"testing"

	"github.com/hjwalt/platform/agent/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillErrorReportsMessage(t *testing.T) {
	assert := assert.New(t)

	err := agent_skill.NewParseError("something broke")

	assert.Equal("something broke", err.Error())
	var _ error = err
}

func TestSkillErrorSatisfiesErrorInterface(t *testing.T) {
	var err error = agent_skill.SkillError{Message: "boom"}
	require.NotNil(t, err)
	assert.Equal(t, "boom", err.Error())
}

func TestNewParseError(t *testing.T) {
	assert := assert.New(t)

	err := agent_skill.NewParseError("bad frontmatter")

	assert.Equal("bad frontmatter", err.Message)
	assert.Equal("bad frontmatter", err.Error())
}

func TestNewParseErrorfFormatsArguments(t *testing.T) {
	assert := assert.New(t)

	err := agent_skill.NewParseErrorf("failed at %s on line %d", "parse", 42)

	assert.Equal("failed at parse on line 42", err.Error())
}

func TestNewValidationErrorCreatesSingleErrorList(t *testing.T) {
	assert := assert.New(t)

	err := agent_skill.NewValidationError("name is required")

	require.NotNil(t, err)
	assert.Equal("name is required", err.Error())
	assert.Equal([]string{"name is required"}, err.Errors)

	var asErr error = err
	require.NotNil(t, asErr)

	var target *agent_skill.ValidationError
	assert.True(errors.As(asErr, &target))
}

func TestNewValidationErrorsKeepsFullList(t *testing.T) {
	assert := assert.New(t)

	err := agent_skill.NewValidationErrors(
		"skill validation failed: 2 problems found",
		[]string{"name is required", "description is required"},
	)

	require.NotNil(t, err)
	assert.Equal("skill validation failed: 2 problems found", err.Error())
	assert.Equal([]string{"name is required", "description is required"}, err.Errors)
}

func TestValidationErrorEmbedsSkillError(t *testing.T) {
	assert := assert.New(t)

	err := agent_skill.NewValidationError("bad")

	// the embedded SkillError exposes Message and Error()
	assert.Equal("bad", err.Message)
	assert.Equal("bad", err.SkillError.Message)
}
