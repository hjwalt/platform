package render

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/hjwalt/platform/web/route"
)

func Handle[C context.Context](h Handler[C], errorHandler ErrorHandler[C]) route.Handler[C] {
	routeHandler := handler[C]{
		Handler:      h,
		ErrorHandler: errorHandler,
	}
	return &routeHandler
}

type handler[C context.Context] struct {
	Handler      Handler[C]
	ErrorHandler ErrorHandler[C]
}

func (h *handler[C]) Handle(ctx C, w http.ResponseWriter, r *http.Request) {
	view, pageErr := h.Handler(ctx, w, r)
	if pageErr != nil {
		h.Error(ctx, w, r, pageErr)
		return
	}

	if view == nil {
		// TODO: handle redirects better
		return
	}

	rendered, renderErr := view.Render()
	if renderErr != nil {
		h.Error(ctx, w, r, pageErr)
		return
	}

	_, writeErr := w.Write([]byte(rendered))
	if writeErr != nil {
		h.Error(ctx, w, r, writeErr)
		return
	}
}

func (h *handler[C]) Error(ctx C, w http.ResponseWriter, r *http.Request, err error) {
	errView := h.ErrorHandler(ctx, w, r, err)

	rendered, renderErr := errView.Render()
	if renderErr != nil {
		slog.Error("error handling error page", "error", renderErr)
		return
	}

	_, writeErr := w.Write([]byte(rendered))
	if writeErr != nil {
		slog.Error("error writing error page", "error", renderErr)
		return
	}
}
