package runtime_chi

import (
	"context"
	"net/http"

	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/routes/page"
	"github.com/hjwalt/platform/routes/route"
)

func AddMiddleware[C context.Context](ic inverse.Container, middleware ...route.Middleware) {
	for _, m := range middleware {
		ic.AddVal(QualifierChiMiddleware, m)
	}
}

func AddRoute[C context.Context](ic inverse.Container, routes ...route.Adder[C]) {
	for _, m := range routes {
		ic.AddVal(QualifierChiRoute, m)
	}
}

func AddDecorator[C context.Context](decorator ...route.Decorator[C]) route.Adder[C] {
	return func(config *route.Configuration[C]) {
		config.AddDecorators(decorator...)
	}
}

func AddPage[C context.Context, M any](path string, method string, pageHandler page.Handler[C, M], errorHandler page.Error[C]) route.Adder[C] {
	return func(config *route.Configuration[C]) {
		page.Add(config, path, method, pageHandler, errorHandler)
	}
}

func AddController[C context.Context](path string, method string, controller mvc.Controller[C], errorHandler mvc.Error[C]) route.Adder[C] {
	return func(config *route.Configuration[C]) {
		mvc.Add(config, path, method, controller, errorHandler)
	}
}

func AddCustom[C context.Context](path string, method string, handler route.Handler[C], errorHandler route.Error) route.Adder[C] {
	return func(config *route.Configuration[C]) {
		route.Add(config, path, method, handler, errorHandler)
	}
}

func AddHttpHandlerRoute[C context.Context](path string, method string, handler http.Handler) route.Adder[C] {
	return func(config *route.Configuration[C]) {
		config.AddRoute(path, method, handler)
	}
}

func AddStatic[C context.Context](prefix string, dir string) route.Adder[C] {
	return func(config *route.Configuration[C]) {
		config.AddRoute(prefix+"*", http.MethodGet, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
	}
}
