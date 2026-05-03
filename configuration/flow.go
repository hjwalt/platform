package configuration

import (
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/flow/stateless"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message/kafka"
)

func RegisterKafkaAgentFlow(holder Context, conf Configuration) {
	// Producer

	kafkaProducer := kafka.NewProducer(conf.Flow.Agent.Producer)
	holder.Add(kafkaProducer)

	messageProducer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New(conf.Flow.Agent.Topic),
			format.Json[agent.Message](),
		),
	)
	holder.Add(messageProducer)
	holder.SetAgentMessageProducer(messageProducer)

	// Consumer

	agentFlow := harness.OpenAiFlow{
		Store: holder.GetRagStore(),
		Tools: holder.GetTool(),
		Model: holder.GetLanguageModel(),
	}
	chatConsumer := kafka.NewConsumer(
		conf.Flow.Agent.Consumer,
		converter.FlowToRuntimeHandler(
			stateless.NewExploder(
				"agent_handle",
				agentFlow.Handle,
				metadata.MessageUpdate(),
				holder.GetAgentMessageProducer(),
				holder.GetAgentMessageProducer(),
			),
			converter.NewConverter(
				flow_runtime_kafka.New(conf.Flow.Agent.Topic),
				format.Json[agent.Message](),
			),
		),
	)
	holder.Add(chatConsumer)
}
