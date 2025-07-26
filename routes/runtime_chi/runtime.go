package runtime_chi

import (
	"context"
	"net/http"

	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/routes/page"
	"github.com/hjwalt/platform/routes/route"
	"github.com/hjwalt/platform/runway/runtime"
	"github.com/hjwalt/platform/runway/structure"
)

var ports = structure.NewSet[int]()

func New[C context.Context](configurations ...runtime.Configuration[*Runtime[C]]) runtime.Runtime {
	r := &Runtime[C]{
		middlewares:   []route.Middleware{},
		configuration: route.NewConfiguration[C](),
	}

	for _, configuration := range configurations {
		r = configuration(r)
	}

	runtimeConfiguration := append(r.httpConfiguration, runtime.HttpWithHandler(HttpHandler(r.middlewares, r.configuration)))

	return runtime.NewHttp(runtimeConfiguration...)
}

func WithPort[C context.Context](port int) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		for ports.Contain(port) {
			port = port + 1
		}

		ports.Add(port)
		r.httpConfiguration = append(r.httpConfiguration, runtime.HttpWithPort(port))
		return r
	}
}

func WithTls[C context.Context](certPath string, keyPath string) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		r.httpConfiguration = append(r.httpConfiguration, runtime.HttpWithTls(certPath, keyPath))
		return r
	}
}

func WithHttpConfiguration[C context.Context](config runtime.Configuration[*runtime.HttpRunnable]) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		r.httpConfiguration = append(r.httpConfiguration, config)
		return r
	}
}

func WithDecorator[C context.Context](decorator ...route.Decorator[C]) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		r.configuration.AddDecorators(decorator...)
		return r
	}
}

func WithMiddleware[C context.Context](middleware ...route.Middleware) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		r.middlewares = append(r.middlewares, middleware...)
		return r
	}
}

func WithPage[C context.Context, M any](path string, method string, pageHandler page.Handler[C, M], errorHandler page.Error[C]) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		page.Add(r.configuration, path, method, pageHandler, errorHandler)
		return r
	}
}

func WithController[C context.Context](path string, method string, controller mvc.Controller[C], errorHandler mvc.Error[C]) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		mvc.Add(r.configuration, path, method, controller, errorHandler)
		return r
	}
}

func WithCustom[C context.Context](path string, method string, handler route.Handler[C], errorHandler route.Error) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		route.Add(r.configuration, path, method, handler, errorHandler)
		return r
	}
}

func WithHttpHandler[C context.Context](path string, method string, handler http.Handler) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		r.configuration.AddRoute(path, method, handler)
		return r
	}
}

func WithStatic[C context.Context](prefix string, dir string) runtime.Configuration[*Runtime[C]] {
	return func(r *Runtime[C]) *Runtime[C] {
		r.configuration.AddRoute(prefix+"*", http.MethodGet, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
		return r
	}
}

type Runtime[C context.Context] struct {
	httpConfiguration []runtime.Configuration[*runtime.HttpRunnable]
	middlewares       []route.Middleware
	configuration     *route.Configuration[C]
}
