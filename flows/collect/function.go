package collect

import (
	"context"

	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/flows/stateful"
)

type Aggregator[S any, K any, V any] func(context.Context, flow.Message[K, V], stateful.State[S]) (stateful.State[S], error)

type Collector func(ctx context.Context, persistenceId string, s stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error)

type OneToOneCollector[S any, OK any, OV any] func(ctx context.Context, persistenceId string, s stateful.State[S]) (*flow.Message[OK, OV], error)

type OneToTwoCollector[S any, OK1 any, OV1 any, OK2 any, OV2 any] func(ctx context.Context, persistenceId string, s stateful.State[S]) (*flow.Message[OK1, OV1], *flow.Message[OK2, OV2], error)
