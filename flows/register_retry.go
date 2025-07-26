package flows

import (
	"context"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/runtime"
)

const (
	QualifierRetry = "QualifierRetry"
)

// Retry
func RegisterRetry(
	container inverse.Container,
	configs []runtime.Configuration[*runtime.Retry],
) {

	resolver := runtime.NewResolver[*runtime.Retry, *runtime.Retry](
		QualifierRetry,
		container,
		false,
		runtime.NewRetry,
	)

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRetry, container)
}

func GetRetry(ctx context.Context, ci inverse.Container) (*runtime.Retry, error) {
	return inverse.GenericGetLast[*runtime.Retry](ci, ctx, QualifierRetry)
}
