package example_task_executor

import (
	"context"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/logger"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/flows"
	"github.com/hjwalt/platform/flows/runtime_rabbit"
	"github.com/hjwalt/platform/flows/task"
	"go.uber.org/zap"
)

func fn(c context.Context, t task.Message[string]) error {
	logger.Info("handling", zap.String("body", t.Value))
	return nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.ExecutorConfiguration[string]{
		Name:                 Instance,
		TaskChannel:          task.StringChannel("tasks"),
		TaskExecutor:         fn,
		TaskConnectionString: "amqp://rabbit:rabbit@localhost:5672/",
		HttpPort:             8081,
		RabbitConsumerConfiguration: []runtime.Configuration[*runtime_rabbit.Consumer]{
			runtime_rabbit.WithConsumerQueueDurable(false),
		},
	}
}

const (
	Instance = "tasks-example-executor"
)

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
