package main

import (
	"github.com/hjwalt/platform/commons/environment"
	"github.com/hjwalt/platform/flows"
	"github.com/hjwalt/platform/flows/example/example_task_cron"
	"github.com/hjwalt/platform/flows/example/example_task_executor"
	"github.com/hjwalt/platform/flows/example/example_word_adapter"
	"github.com/hjwalt/platform/flows/example/example_word_collect"
	"github.com/hjwalt/platform/flows/example/example_word_count"
	"github.com/hjwalt/platform/flows/example/example_word_join"
	"github.com/hjwalt/platform/flows/example/example_word_materialise"
	"github.com/hjwalt/platform/flows/example/example_word_remap"
)

func main() {
	m := flows.NewMain()

	example_word_adapter.Register(m)
	example_word_collect.Register(m)
	example_word_count.Register(m)
	example_word_join.Register(m)
	example_word_materialise.Register(m)
	example_word_remap.Register(m)
	example_task_executor.Register(m)
	example_task_cron.Register(m)

	err := m.Start(environment.GetString("INSTANCE", flows.AllInstances))

	if err != nil {
		panic(err)
	}
}
