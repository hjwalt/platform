package runtime_chi

import (
	"context"
	"net/http"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/managed"
	"github.com/hjwalt/platform/routes/route"
)

const (
	QualifierChiMiddleware = "chi_middleware"
	QualifierChiRoute      = "chi_route_builder"
)

func AddHttpHandler[C context.Context](ic inverse.Container) {
	managed.AddComponent(ic, &chiComponent[C]{})
}

// implementation
type chiComponent[C context.Context] struct {
}

func (c *chiComponent[C]) Name() string {
	return "chi"
}

func (r *chiComponent[C]) Register(rctx context.Context, ic inverse.Container) error {
	managed.AddHttpHandler(ic, func(ctx context.Context, c inverse.Container) (http.Handler, error) {
		middlewares, middlewareErr := inverse.GenericGetAll[route.Middleware](ic, ctx, QualifierChiMiddleware)
		if middlewareErr != nil {
			return nil, middlewareErr
		}

		routeConfig, routeConfigErr := inverse.GenericGetAll[route.Adder[C]](ic, ctx, QualifierChiRoute)
		if routeConfigErr != nil {
			return nil, routeConfigErr
		}

		configuration := route.NewConfiguration[C]()
		for _, config := range routeConfig {
			config(configuration)
		}

		return HttpHandler(middlewares, configuration), nil
	})
	return nil
}

func (r *chiComponent[C]) Resolve(ctx context.Context, ic inverse.Container) error {
	return nil
}

func (r *chiComponent[C]) Clean() error {
	return nil
}
