package runtime

import (
	"errors"
	"log/slog"
	"sync"
)

type Runnable interface {
	Start() error
	Stop()
	Run() error
}

func NewRunner(runnable Runnable) Runtime {
	return &Runner{
		runnable: runnable,
	}
}

type Runner struct {
	wait     sync.WaitGroup
	mu       sync.Mutex
	started  bool
	runnable Runnable
}

func (r *Runner) Start() error {
	defer r.mu.Unlock()
	r.mu.Lock()

	if r.runnable == nil {
		return ErrRunnerRuntimeNoRunnable
	}
	if err := r.runnable.Start(); err != nil {
		return errors.Join(ErrRunnerRuntimeRunnableStart, err)
	}
	go r.Run()
	r.wait.Add(1)
	r.started = true
	return nil
}

func (r *Runner) Stop() {
	defer r.mu.Unlock()
	r.mu.Lock()

	if !r.started {
		return
	}

	r.started = false
	r.runnable.Stop()
	r.wait.Wait()
}

func (r *Runner) Run() {
	defer r.wait.Done()
	if err := r.runnable.Run(); err != nil {
		slog.Error("runner runtime run error", "error", err)
		Error(err)
	}
}

// Errors
var (
	ErrRunnerRuntimeNoRunnable    = errors.New("runner runtime no runnable provided")
	ErrRunnerRuntimeRunnableStart = errors.New("runner runtime runnable failed to start")
)
