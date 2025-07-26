package task_executor_retry

import (
	"context"
	"errors"

	"github.com/hjwalt/platform/commons/logger"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/task"
	"go.uber.org/zap"
)

// constructor
func New(configurations ...runtime.Configuration[*TaskRetry]) task.Executor[structure.Bytes] {
	function := &TaskRetry{}
	for _, configuration := range configurations {
		function = configuration(function)
	}
	return function.Apply
}

// configurations
func WithRetry(r *runtime.Retry) runtime.Configuration[*TaskRetry] {
	return func(p *TaskRetry) *TaskRetry {
		p.retry = r
		return p
	}
}

func WithExecutor(e task.Executor[structure.Bytes]) runtime.Configuration[*TaskRetry] {
	return func(p *TaskRetry) *TaskRetry {
		p.executor = e
		return p
	}
}

// implementation
type TaskRetry struct {
	retry    *runtime.Retry
	executor task.Executor[structure.Bytes]
}

func (r *TaskRetry) Apply(c context.Context, t task.Message[structure.Bytes]) error {
	retryErr := r.retry.Do(func(tryCount int64) error {
		err := r.executor(runtime.SetRetryCount(c, tryCount), t)
		if err != nil {
			logger.Warn("retrying", zap.Int64("try", tryCount), zap.Error(err))
			return err
		}
		return nil
	})

	if retryErr != nil {
		return errors.Join(ErrorRetryAttempt, retryErr)
	} else {
		return nil
	}
}

var ErrorRetryAttempt = errors.New("retry all attempts failed")
