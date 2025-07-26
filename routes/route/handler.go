package route

import (
	"context"
	"net/http"

	"github.com/hjwalt/platform/runway/reflect"
)

type Custom[C context.Context] struct {
	Decorators []Decorator[C]
	Handler    Handler[C]
	Error      Error
}

func (p *Custom[C]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := reflect.Construct[C]()

	var ctxErr error
	for _, decorator := range p.Decorators {
		ctx, ctxErr = decorator(ctx, w, r)
		if ctxErr != nil {
			p.Error(w, r, ctxErr)
			return
		}
	}

	p.Handler(ctx, w, r)
}

func Add[C context.Context](c *Configuration[C], r string, m string, p Handler[C], e Error) {
	routeHandler := &Custom[C]{
		Decorators: c.Decorators,
		Handler:    p,
		Error:      e,
	}
	c.AddRoute(r, m, routeHandler)
}
