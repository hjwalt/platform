package stateful_error_handler_skip_list

import (
	"context"
	"errors"

	"github.com/hjwalt/platform/commons/format"
	"github.com/hjwalt/platform/commons/logger"
	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/flows/stateful"
)

func New(skipList []error) stateful.ErrorHandlerFunction {
	f := fn{skipList: skipList}
	return f.handle
}

func Default() stateful.ErrorHandlerFunction {
	return New(
		[]error{
			format.ErrFormatConversionMarshal,
			format.ErrFormatConversionUnmarshal,
		},
	)
}

type fn struct {
	skipList []error
}

func (r *fn) handle(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], s stateful.State[structure.Bytes], e error) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {
	for _, skipErr := range r.skipList {
		if errors.Is(e, skipErr) {
			logger.WarnErr("stateful skipped error", e)
			return flow.EmptySlice(), s, nil
		}
	}
	return flow.EmptySlice(), s, e
}
