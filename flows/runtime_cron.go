package flows

import (
	"context"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/flows/runtime_cron"
	"github.com/hjwalt/platform/flows/runtime_rabbit"
	"github.com/hjwalt/platform/flows/task"
	"github.com/hjwalt/platform/routes/runtime_chi"
)

type CronConfiguration[T any] struct {
	Name                        string
	TaskChannel                 task.Channel[T]
	Scheduler                   task.Scheduler[T]
	Schedules                   []string
	TaskConnectionString        string
	HttpPort                    int
	RabbitProducerConfiguration []runtime.Configuration[*runtime_rabbit.Producer]
	RetryConfiguration          []runtime.Configuration[*runtime.Retry]
	RouteConfiguration          []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
}

func (c CronConfiguration[T]) Register(ci inverse.Container) {
	RegisterRetry(
		ci,
		c.RetryConfiguration,
	)
	RegisterRabbitProducer(
		ci,
		c.Name,
		c.TaskConnectionString,
		c.RabbitProducerConfiguration,
	)
	RegisterRoute(
		ci,
		c.HttpPort,
		c.RouteConfiguration,
	)
	RegisterCron(
		ci,
		[]runtime.Configuration[*runtime_cron.Cron]{},
	)

	// ADDING CRON AFTER CONFIG
	// Moving this above the cron config will result in NPE

	for _, schedule := range c.Schedules {
		RegisterCronConfig(ci, runtime_cron.WithCronJob(schedule, c.Scheduler, c.TaskChannel))
	}
}
