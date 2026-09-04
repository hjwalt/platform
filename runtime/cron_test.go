package runtime

import (
	"testing"
	"time"

	cronlib "github.com/hjwalt/platform/cron"
	"github.com/stretchr/testify/assert"
)

// The runtime cron struct is unexported and has no constructor in the repo,
// so it is built directly in this same-package test.
func TestCronStartStopLifecycle(t *testing.T) {
	assert := assert.New(t)

	c := &cron{cron: cronlib.New()}

	// Stopping a cron that was never started must not panic or hang.
	c.Stop()

	assert.NoError(c.Start())

	// Starting twice is a no-op but must not error.
	assert.NoError(c.Start())

	c.Stop()

	// A stopped cron can be started again.
	assert.NoError(c.Start())
	c.Stop()
}

func TestCronStopReturnsPromptly(t *testing.T) {
	assert := assert.New(t)

	c := &cron{cron: cronlib.New()}
	assert.NoError(c.Start())

	// Stop() handshakes with the scheduler goroutine (sends on its stop
	// channel); with no jobs it must return promptly, not hang.
	start := time.Now()
	c.Stop()
	elapsed := time.Since(start)

	assert.Less(elapsed, 5*time.Second, "Stop should not block for a long time")
}
