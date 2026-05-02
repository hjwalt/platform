package configuration

import (
	"log/slog"
	"time"

	"github.com/hjwalt/platform/runtime"
)

func Holder() runtime.Holder {
	return &holder{}
}

type holder struct {
	Runtimes []runtime.Runtime
}

func (r *holder) Add(runtimes ...runtime.Runtime) {
	r.Runtimes = append(r.Runtimes, runtimes...)
}

func (r *holder) Block() {
	startErr := runtime.Start(
		r.Runtimes,
		time.Second,
	)

	if startErr != nil {
		panic(startErr)
	}

	slog.Info("started")

	runtime.Wait()

	slog.Info("stopped")
}
