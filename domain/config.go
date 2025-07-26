package domain

import (
	"strings"

	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/routes/page"
	"github.com/hjwalt/platform/routes/runtime_chi"
	"github.com/hjwalt/platform/web"
)

func Page[M any](path string, method string, pageHandler page.Handler[web.Context, M], errorHandler page.Error[web.Context]) runtime.Configuration[*runtime_chi.Runtime[web.Context]] {
	return runtime_chi.WithPage(strings.TrimSpace(path), method, pageHandler, errorHandler)
}

func Component[M any](path string, method string, pageHandler page.Handler[web.Context, M], errorHandler page.Error[web.Context]) runtime.Configuration[*runtime_chi.Runtime[web.Context]] {
	return runtime_chi.WithPage(strings.TrimSpace(path), method, pageHandler, errorHandler)
}

func Controller(path string, method string, pageHandler mvc.Controller[web.Context], errorHandler mvc.Error[web.Context]) runtime.Configuration[*runtime_chi.Runtime[web.Context]] {
	return runtime_chi.WithController(strings.TrimSpace(path), method, pageHandler, errorHandler)
}
