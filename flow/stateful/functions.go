package stateful

import (
	"context"

	"github.com/hjwalt/platform/type/either"
	"github.com/hjwalt/platform/type/optional"
)

type StateKey[IV any] func(context.Context, IV) (string, error)

type StateUpdate[IV any, ST any, ERR any] func(context.Context, IV, ST) either.Either[ST, ERR]

type Operate[IV any, OV any, ST any, ERR any] func(context.Context, IV, ST) (optional.Optional[OV], optional.Optional[ERR])
