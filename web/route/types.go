package route

import (
	"context"
	"net/http"
)

type Builder[C context.Context] interface {
	AddDecorators(decorator ...Decorator[C])
	AddMiddlewares(middlewares ...Middleware)
	HandleStatic(prefix string, dir string)
	Handle(fullPath string, method string, handler Handler[C])
	Build() http.Handler
}

type Handler[C context.Context] interface {
	Handle(c C, w http.ResponseWriter, r *http.Request)
	Error(c C, w http.ResponseWriter, r *http.Request, err error)
}

type Error[C context.Context] func(c C, w http.ResponseWriter, r *http.Request, err error)

type Decorator[C context.Context] func(c C, w http.ResponseWriter, r *http.Request) (C, error)

type Middleware func(http.Handler) http.Handler

type Adder[C context.Context] func(Builder[C])
