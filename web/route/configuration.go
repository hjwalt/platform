package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func NewConfiguration[C context.Context]() Builder[C] {
	return &Configuration[C]{
		Routes:     map[string]*Configuration[C]{},
		Decorators: []Decorator[C]{},
	}
}

func newConfiguration[C context.Context]() *Configuration[C] {
	return &Configuration[C]{
		Routes:     map[string]*Configuration[C]{},
		Decorators: []Decorator[C]{},
	}
}

type Configuration[C context.Context] struct {
	Delete      http.Handler
	Get         http.Handler
	Post        http.Handler
	Put         http.Handler
	Routes      map[string]*Configuration[C]
	Decorators  []Decorator[C]
	Middlewares []Middleware
}

func (config *Configuration[C]) AddRoute(fullPath string, method string, handler http.Handler) {
	path := strings.TrimSpace(fullPath)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")

	config.addParts(parts, method, handler)
}

func (config *Configuration[C]) addParts(parts []string, method string, handler http.Handler) {
	if len(parts) == 0 {
		switch method {
		case http.MethodDelete:
			config.Delete = handler
		case http.MethodGet:
			config.Get = handler
		case http.MethodPost:
			config.Post = handler
		case http.MethodPut:
			config.Put = handler
		default:
			config.Get = handler
		}
	} else {
		if curr, currExists := config.Routes[parts[0]]; currExists {
			curr.addParts(parts[1:], method, handler)
		} else {
			curr := newConfiguration[C]()
			config.Routes[parts[0]] = curr
			curr.addParts(parts[1:], method, handler)
		}
	}
}

func (config *Configuration[C]) AddDecorators(decorator ...Decorator[C]) {
	config.Decorators = append(config.Decorators, decorator...)
}

func (config *Configuration[C]) AddMiddlewares(middleware ...Middleware) {
	config.Middlewares = append(config.Middlewares, middleware...)
}

func (c *Configuration[C]) Handle(fullPath string, method string, handler Handler[C]) {
	routeHandler := &RouteHandler[C]{
		Decorators: c.Decorators,
		Handler:    handler,
	}
	c.AddRoute(fullPath, method, routeHandler)
}

func (config *Configuration[C]) HandleStatic(prefix string, dir string) {
	config.AddRoute(prefix+"*", http.MethodGet, http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
}

func (config *Configuration[C]) Build() http.Handler {

	r := chi.NewRouter()

	for _, m := range config.Middlewares {
		r.Use(m)
	}

	handleCurrent(config, "", r)
	handleSubPath(config, r)

	return r
}

func handleSubPath[C context.Context](config *Configuration[C], router chi.Router) {
	for subPath, subConfig := range config.Routes {
		if len(subConfig.Routes) > 0 {
			router.Route(
				"/"+subPath,
				func(r chi.Router) {
					handleCurrent(subConfig, "", r)
					handleSubPath(subConfig, r)
				},
			)
		} else {
			handleCurrent(subConfig, subPath, router)
		}
	}
}

func handleCurrent[C context.Context](config *Configuration[C], path string, router chi.Router) {
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
