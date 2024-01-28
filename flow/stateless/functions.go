package stateless

import (
	"context"

	"github.com/hjwalt/platform/type/optional"
)

type Consume[IV any, ERR any] func(context.Context, IV) optional.Optional[ERR]

type Operate[IV any, OV any, ERR any] func(context.Context, IV) (optional.Optional[OV], optional.Optional[ERR])

type Explode[IV any, OV any, ERR any] func(context.Context, IV) (optional.Optional[[]OV], optional.Optional[ERR])
