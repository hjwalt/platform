package agent_skill

import (
	"fmt"

	"github.com/hjwalt/platform/type/optional"
)

type Skill struct {
	Name          string
	Description   string
	License       optional.Optional[string]
	Compatibility optional.Optional[string]
	AllowedTools  []string
	Metadata      map[string]string
	Body          string
}

type SkillError struct {
	Message string
}

// Error implements the error interface.
func (e SkillError) Error() string {
	return e.Message
}

// NewParseError creates a new ParseError with the given message.
func NewParseError(message string) SkillError {
	return SkillError{Message: message}
}

// NewParseErrorf creates a new ParseError with a formatted message.
func NewParseErrorf(format string, args ...interface{}) SkillError {
	return SkillError{Message: fmt.Sprintf(format, args...)}
}

// ValidationError indicates that skill properties failed validation.
// The Errors field contains all validation error messages.
type ValidationError struct {
	SkillError
	// Errors contains all validation error messages found.
	Errors []string
}

// NewValidationError creates a new ValidationError with a single error message.
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		SkillError: SkillError{Message: message},
		Errors:     []string{message},
	}
}

// NewValidationErrors creates a new ValidationError with multiple error messages.
// The message parameter is used as the primary error message, while errors
// contains the complete list of validation failures.
func NewValidationErrors(message string, errors []string) *ValidationError {
	return &ValidationError{
		SkillError: SkillError{Message: message},
		Errors:     errors,
	}
}
