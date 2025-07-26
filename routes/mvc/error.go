package mvc

import (
	"context"
	"log/slog"
	"net/http"
)

type Error[C context.Context] func(c C, w http.ResponseWriter, r *http.Request, err error) View[C]

func handleError[C context.Context](ctx C, w http.ResponseWriter, r *http.Request, err error, errHandler Error[C]) {
	errView := errHandler(ctx, w, r, err)
	errTemplateErr := errView.Write(ctx, w, r)
	if errTemplateErr != nil {
		slog.Error("error handling error", "error", errTemplateErr)
	}
}
