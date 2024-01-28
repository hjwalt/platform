package route

import (
	"context"
	"net/http"

	"github.com/hjwalt/platform/reflect"
)

type RouteHandler[C context.Context] struct {
	Decorators []Decorator[C]
	Handler    Handler[C]
}

func (p *RouteHandler[C]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := reflect.Construct[C]()

	var ctxErr error
	for _, decorator := range p.Decorators {
		ctx, ctxErr = decorator(ctx, w, r)
		if ctxErr != nil {
			p.Handler.Error(ctx, w, r, ctxErr)
			return
		}
	}

	p.Handler.Handle(ctx, w, r)
}
