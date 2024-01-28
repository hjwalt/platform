package main

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/example/route/page_billing"
	"github.com/hjwalt/platform/example/route/page_home"
	"github.com/hjwalt/platform/logger"
	"github.com/hjwalt/platform/runtime"
	"github.com/hjwalt/platform/web/route"
	"github.com/joho/godotenv"
)

func main() {
	logger.Default()
	godotenv.Load()

	httpBuilder := route.NewConfiguration[example.Context]()

	httpBuilder.AddMiddlewares(
		middleware.RequestID,
		middleware.RealIP,
		middleware.CleanPath,
		middleware.Recoverer,
	)
	httpBuilder.HandleStatic("/static/", "./example/route/static")
	page_home.Add(httpBuilder)
	page_billing.Add(httpBuilder)

	httpHandler := httpBuilder.Build()

	httpRuntime := runtime.NewHttp(
		runtime.HttpWithPort(3001),
		runtime.HttpWithHandler(httpHandler),
	)

	startErr := runtime.Start(
		[]runtime.Runtime{
			httpRuntime,
		},
		time.Second,
	)

	if startErr != nil {
		panic(startErr)
	}

	slog.Info("started")

	runtime.Wait()

	slog.Info("stopped")
}
