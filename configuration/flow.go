package configuration

import (
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_kafka"
	"github.com/hjwalt/platform/flow/stateful"
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

	resultProducer := converter.RuntimeToFlowProducer(
		kafkaProducer,
		converter.NewConverter(
			flow_runtime_kafka.New(conf.Flow.Result.Topic),
			format.Json[agent.Result](),
		),
	)
	holder.Add(resultProducer)

	// Consumer

	agentFlow := harness.Flow{
		Tools: holder.GetTool(),
		Model: holder.GetLanguageModel(),
	}

	chatConsumer := kafka.NewConsumer(
		conf.Flow.Agent.Consumer,
		converter.FlowToRuntimeHandler(
			stateful.NewOperator(
				"agent_handle",
				agentFlow.Key,
				agentFlow.Update,
				agentFlow.Next,
				agentFlow.ResultMetadata,
				agentFlow.MessageMetadata,
				resultProducer,
				holder.GetAgentMessageProducer(),
				converter.RuntimeToFlowStore(
					holder.GetAgentHarnessStore(),
					format.Json[harness.ExecutionState](),
				),
			),
			converter.NewConverter(
				flow_runtime_kafka.New(conf.Flow.Result.Topic), // irrelevant
				format.Json[agent.Message](),
			),
		),
	)
	holder.Add(chatConsumer)

	resultConsumer := kafka.NewConsumer(
		conf.Flow.Result.Consumer,
		converter.FlowToRuntimeHandler(
			stateless.NewExploder(
				"agent_explode",
				agentFlow.Explode,
				agentFlow.MessageMetadata,
				agentFlow.MessageMetadata,
				holder.GetAgentMessageProducer(),
				holder.GetAgentMessageProducer(),
			),
			converter.NewConverter(
				flow_runtime_kafka.New(conf.Flow.Agent.Topic), // irrelevant
				format.Json[agent.Result](),
			),
		),
	)
	holder.Add(resultConsumer)
}
