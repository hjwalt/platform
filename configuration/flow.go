package configuration

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/flow/converter"
	"github.com/hjwalt/platform/flow/flow_runtime_memory"
	"github.com/hjwalt/platform/flow/metadata"
	"github.com/hjwalt/platform/flow/stateless"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/message/memory"
	"github.com/hjwalt/platform/runtime"
)

func RegisterInMemoryAgentMessageProducer(
	holder runtime.Holder,
	conf Configuration,
	agentMessageChannel memory.MemoryConfiguration,
) flow.Producer[agent.Message] {
	messageProducer := converter.RuntimeToFlowProducer(
		memory.NewProducer(agentMessageChannel),
		converter.NewConverter(
			flow_runtime_memory.New(),
			format.Json[agent.Message](),
		),
	)
	holder.Add(messageProducer)
	return messageProducer
}

func RegisterInMemoryAgentMessageConsumer(
	holder runtime.Holder,
	conf Configuration,
	agentMessageChannel memory.MemoryConfiguration,
	messageProducer flow.Producer[agent.Message],
	model agent.LanguageModel,
	tools []agent.Tool,
) {
	toolMap := map[string]agent.Tool{}
	for _, tool := range tools {
		toolMap[tool.Name()] = tool
	}
	agentFlow := harness.OpenAiFlow[context.Context]{
		Tools: toolMap,
		Model: model,
	}
	chatConsumer := memory.NewConsumer(
		agentMessageChannel,
		converter.FlowToRuntimeHandler(
			stateless.NewExploder(
				"agent_handle",
				agentFlow.Handle,
				metadata.MessageUpdate(),
				messageProducer,
				messageProducer,
			),
			converter.NewConverter(
				flow_runtime_memory.New(),
				format.Json[agent.Message](),
			),
		),
	)
	holder.Add(chatConsumer)
}
