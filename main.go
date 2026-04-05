package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/domain"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/routes/runtime_chi"
	"github.com/hjwalt/platform/state"
	"github.com/hjwalt/platform/state/file"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/component/component_home"
	"github.com/hjwalt/platform/web/component/component_navbar"
	"github.com/hjwalt/platform/web/component/component_sidebar"
	"github.com/hjwalt/platform/web/page/page_base"
	"github.com/hjwalt/platform/web/page/page_error_500"
)

func main() {
	storageFormat := format.Protojson[*model.ProtobufSchema]()

	fileStore := &file.FileStore{
		Path: "/home/hjwalt/Projects/platform/tmp",
	}

	formattedStore := &state.FormattedStore[*model.ProtobufSchema]{
		Format: storageFormat,
		Store:  fileStore,
	}

	conf, err := formattedStore.Read(context.Background(), "protobuf.json")
	if err != nil {
		panic(err)
	}
	typeMap := domain.Parse(conf.Value)

	runtimeDecorator := web.DecoratorContext{
		Schema:  conf.Value,
		TypeMap: typeMap,
	}

	httpRuntime := runtime_chi.New(
		runtime_chi.WithPort[web.Context](3000),

		runtime_chi.WithMiddleware[web.Context](middleware.RequestID),
		runtime_chi.WithMiddleware[web.Context](middleware.RealIP),
		// runtime_chi.WithMiddleware[web.Context](middleware.CleanPath),
		runtime_chi.WithMiddleware[web.Context](middleware.Recoverer),

		runtime_chi.WithDecorator(runtimeDecorator.Decorate),
		runtime_chi.WithDecorator(web.DecoratorHtmx),

		runtime_chi.WithStatic[web.Context]("/static/", "./web/static"),

		domain.Page("/", http.MethodGet, page_base.Page, page_error_500.Error),

		component_home.Get(),
		component_sidebar.Get(),
		component_navbar.Get(),
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

	runtime.Wait()

	// logger.Info("conf", zap.Any("conf", conf))

	// logger.Info("test", zap.Bool("test", domain.InUse("Test", conf)))
	// logger.Info("test", zap.Bool("test", domain.InUse("ProtobufMessage", conf)))

	// domain.Flatten(conf, typeMap)
	// store.Write("model/protobuf.json", storageFormat, conf)

	// protoFile, _ := os.Create("model/protobuf.proto")
	// defer protoFile.Close()

	// w := bufio.NewWriter(protoFile)
	// defer w.Flush()

	// write.WriteProtobufSchema(w, conf)
}
