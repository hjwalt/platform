package route

import (
	"context"
	"net/http"
	"strings"
)

func NewConfiguration[C context.Context]() *Configuration[C] {
	return &Configuration[C]{
		Routes:     map[string]*Configuration[C]{},
		Decorators: []Decorator[C]{},
	}
}

type Configuration[C context.Context] struct {
	Delete     http.Handler
	Get        http.Handler
	Post       http.Handler
	Put        http.Handler
	Routes     map[string]*Configuration[C]
	Decorators []Decorator[C]
}

func (config *Configuration[C]) Set(method string, handler http.Handler) {
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
}

func (config *Configuration[C]) AddRoute(fullPath string, method string, handler http.Handler) {
	path := strings.TrimSpace(fullPath)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")

	config.AddRouteParts(parts, method, handler)
}

func (config *Configuration[C]) AddRouteParts(parts []string, method string, handler http.Handler) {
	if len(parts) == 0 {
		config.Set(method, handler)
	} else {
		if curr, currExists := config.Routes[parts[0]]; currExists {
			curr.AddRouteParts(parts[1:], method, handler)
		} else {
			curr := NewConfiguration[C]()
			config.Routes[parts[0]] = curr
			curr.AddRouteParts(parts[1:], method, handler)
		}
	}
}

func (config *Configuration[C]) AddDecorators(decorator ...Decorator[C]) {
	config.Decorators = append(config.Decorators, decorator...)
}
