package render

import (
	"context"
	"net/http"
)

type View interface {
	Render() (string, error)
}

type Handler[C context.Context] func(c C, w http.ResponseWriter, r *http.Request) (View, error)

type ErrorHandler[C context.Context] func(c C, w http.ResponseWriter, r *http.Request, err error) View
