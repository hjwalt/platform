package runtime_chi

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hjwalt/platform/routes/route"
)

func HttpHandler[C context.Context](middlewares []route.Middleware, config *route.Configuration[C]) http.Handler {

	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	HandleCurrent(config, "", r)
	HandleSubpath(config, r)

	return r
}

func HandleSubpath[C context.Context](config *route.Configuration[C], router chi.Router) {
	for subPath, subConfig := range config.Routes {
		if len(subConfig.Routes) > 0 {
			router.Route(
				"/"+subPath,
				func(r chi.Router) {
					HandleCurrent(subConfig, "", r)
					HandleSubpath(subConfig, r)
				},
			)
		} else {
			HandleCurrent(subConfig, subPath, router)
		}
	}
}

func HandleCurrent[C context.Context](config *route.Configuration[C], path string, router chi.Router) {
	nextPath := "/" + path
	if config.Delete != nil {
		router.Method(http.MethodDelete, nextPath, config.Delete)
	}
	if config.Get != nil {
		router.Method(http.MethodGet, nextPath, config.Get)
	}
	if config.Post != nil {
		router.Method(http.MethodPost, nextPath, config.Post)
	}
	if config.Put != nil {
		router.Method(http.MethodPut, nextPath, config.Put)
	}
}
