package managed_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/hjwalt/platform/commons/cron"
	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/managed"
	"github.com/stretchr/testify/assert"
)

func TestCron(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	ic := inverse.NewContainer()

	cj := &mockCronJob{}

	managed.AddCron(ic)
	managed.AddCronJob(ic, cj)

	managed, err := managed.New(ic, ctx)
	assert.NoError(err)

	starterr := managed.Start()
	assert.NoError(starterr)

	time.Sleep(2 * time.Second)

	managed.Stop()

	assert.GreaterOrEqual(cj.count, 2)
}

type mockCronJob struct {
	count int
}

func (m *mockCronJob) Schedule() cron.Schedule {
	return cron.Every(time.Second)
}

func (m *mockCronJob) Run(t time.Time) {
	slog.Info("cron time", "time", t)
	m.count += 1
}
