package trusted_test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	"github.com/hjwalt/platform/trusted"
	"github.com/stretchr/testify/assert"
)

func TestEndian(t *testing.T) {
	assert := assert.New(t)

	order := trusted.Endian()

	// native endianness is machine-dependent, so accept either ordering
	assert.True(
		order == binary.LittleEndian || order == binary.BigEndian,
		"Endian() returned %T, expected binary.LittleEndian or binary.BigEndian",
		order,
	)
}

func TestEndianIsStable(t *testing.T) {
	assert := assert.New(t)

	// repeated calls should always report the same native ordering
	assert.Equal(trusted.Endian(), trusted.Endian())
}

func TestMustReturnsValueOnNilError(t *testing.T) {
	assert := assert.New(t)

	got := trusted.Must("value", nil)
	assert.Equal("value", got)

	gotInt := trusted.Must(42, nil)
	assert.Equal(42, gotInt)

	gotErr := trusted.Must(errors.New("carried"), nil)
	assert.EqualError(gotErr, "carried")
}

func TestMustPanicsOnError(t *testing.T) {
	assert := assert.New(t)

	assert.PanicsWithError(
		"boom",
		func() {
			trusted.Must("value", errors.New("boom"))
		},
	)
}

// Exit(err) with ErrPrimaryTesting is the only error that does not trigger
// os.Exit(1), so it is the only path that can be exercised inside a unit test.
func TestExitReturnsOnPrimaryTestingError(t *testing.T) {
	assert := assert.New(t)

	assert.NotPanics(func() {
		trusted.Exit(trusted.ErrPrimaryTesting)
	})
}

func TestExitReturnsOnWrappedPrimaryTestingError(t *testing.T) {
	assert := assert.New(t)

	// Exit uses errors.Is, so wrapped ErrPrimaryTesting errors also return normally
	wrapped := fmt.Errorf("wrapped: %w", trusted.ErrPrimaryTesting)

	assert.NotPanics(func() {
		trusted.Exit(wrapped)
	})
}
