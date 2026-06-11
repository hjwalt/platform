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

func RegisterKafkaProducer(holder Context, conf Configuration) {
	kafkaProducer := kafka.NewProducer(conf.Flow.Agent.Producer)
	holder.Add(kafkaProducer)
	holder.SetKafkaProducer(kafkaProducer)
}

func RegisterKafkaAgentMessageProducer(holder Context, conf Configuration) {
	flowMetadata := harness.FlowMetadata{}

	messageProducer := converter.RuntimeToFlowProducer(
		holder.GetKafkaProducer(),
		converter.NewConverter(
			flow_runtime_kafka.New(conf.Flow.Agent.Topic),
			format.Json[agent.Message](),
		),
		flowMetadata.MessageMetadata,
	)
	holder.Add(messageProducer)
	holder.SetAgentMessageProducer(messageProducer)
}

func RegisterKafkaAgentFlow(holder Context, conf Configuration) {
	flowMetadata := harness.FlowMetadata{}
	agentFlow := harness.Flow{
		Tools: holder.GetToolContainer(),
		Model: holder.GetAgentModel(),
	}

	// Producer

	resultProducer := converter.RuntimeToFlowProducer(
		holder.GetKafkaProducer(),
		converter.NewConverter(
			flow_runtime_kafka.New(conf.Flow.Result.Topic),
			format.Json[agent.Result](),
		),
		flowMetadata.ResultMetadata,
	)
	holder.Add(resultProducer)

	// Consumer

	chatConsumer := kafka.NewConsumer(
		conf.Flow.Agent.Consumer,
		conf.Flow.Agent.Topic,
		converter.FlowToRuntimeHandler(
			stateful.NewOperator(
				"agent_handle",
				flowMetadata.Key,
				agentFlow.Update,
				agentFlow.Next,
				flowMetadata.ResultMetadata,
				flowMetadata.MessageMetadata,
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
		conf.Flow.Result.Topic,
		converter.FlowToRuntimeHandler(
			stateless.NewExploder(
				"agent_explode",
				agentFlow.Explode,
				flowMetadata.MessageMetadata,
				flowMetadata.MessageMetadata,
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
