package flows

import (
	"context"

	"github.com/hjwalt/platform/commons/format"
	"github.com/hjwalt/platform/commons/inverse"
	"github.com/hjwalt/platform/commons/runtime"
	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/adapter"
	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/flows/runtime_sarama"
	"github.com/hjwalt/platform/flows/stateless"
	"github.com/hjwalt/platform/routes/runtime_chi"
)

// Wiring configuration
type RouterAdapterConfiguration[Request any, InputKey any, InputValue any] struct {
	Name                       string
	Path                       string
	ProduceTopic               flow.Topic[InputKey, InputValue]
	ProduceBroker              string
	RequestBodyFormat          format.Format[Request]
	RequestMapFunction         stateless.OneToOneFunction[structure.Bytes, Request, InputKey, InputValue]
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	RouteConfiguration         []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) Register(ci inverse.Container) {
	RegisterRouteConfig(
		ci,
		c.RouteConfiguration...,
	)
	RegisterProducerRoute(
		ci,
		"POST",
		c.Path,
		adapter.RouteProduceTopicBodyMapConvert(
			c.RequestMapFunction,
			c.RequestBodyFormat,
			c.ProduceTopic,
		),
	)

	// RUNTIME

	RegisterKafkaProducer(
		ci,
		c.ProduceBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterRoute(
		ci,
		c.HttpPort,
		c.RouteConfiguration,
	)
}
