package flows

import (
	"context"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/runtime"
)

const (
	QualifierRuntime = "QualifierRuntime"
)

func RegisterRuntime(qualifier string, ci inverse.Container) {
	ci.Add(QualifierRuntime, func(ctx context.Context, ci2 inverse.Container) (any, error) {
		return ci2.Get(ctx, qualifier)
	})
}

func InjectedRuntimes(ci inverse.Container) []runtime.Runtime {
	allRuntimes, err := inverse.GenericGetAll[runtime.Runtime](ci, context.Background(), QualifierRuntime)
	if err != nil {
		panic(err)
	}
	return allRuntimes
}
