package page

import (
	"context"

	"github.com/hjwalt/platform/routes/route"
)

func Add[C context.Context, M any](c *route.Configuration[C], r string, m string, p Handler[C, M], e Error[C]) {
	routeHandler := &Page[C, M]{
		Decorators:   c.Decorators,
		PageHandler:  p,
		ErrorHandler: e,
	}
	c.AddRoute(r, m, routeHandler)
}
