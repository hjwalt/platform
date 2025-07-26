package page

import (
	"context"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/commons/reflect"
	"github.com/hjwalt/platform/routes/route"
)

type Handler[C context.Context, M any] func(c C, w http.ResponseWriter, r *http.Request) (*template.Template, M, error)

type Page[C context.Context, M any] struct {
	Decorators   []route.Decorator[C]
	PageHandler  Handler[C, M]
	ErrorHandler Error[C]
}

func (p *Page[C, M]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := reflect.Construct[C]()

	var ctxErr error
	for _, decorator := range p.Decorators {
		ctx, ctxErr = decorator(ctx, w, r)
		if ctxErr != nil {
			handleError(ctx, w, r, ctxErr, p.ErrorHandler)
			return
		}
	}

	pageTemplate, pageModel, pageErr := p.PageHandler(ctx, w, r)
	if pageErr != nil {
		handleError(ctx, w, r, pageErr, p.ErrorHandler)
		return
	}

	executeErr := pageTemplate.Execute(w, pageModel)
	if executeErr != nil {
		handleError(ctx, w, r, executeErr, p.ErrorHandler)
		return
	}
}
