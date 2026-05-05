package configuration

import (
	"github.com/hjwalt/platform/agent/mcp/mcp_brave_search_web"
	"github.com/hjwalt/platform/message/kafka"
)

type Configuration struct {
	OpenAi      OpenAiConfiguration
	BraveSearch mcp_brave_search_web.BraveSearchConfiguration
	Server      WebServerConfiguration
	Flow        FlowConfiguration
}

type OpenAiConfiguration struct {
	Model    string
	Endpoint string
	Secret   string
}

type WebServerConfiguration struct {
	Port               int
	StaticResourcePath string
}

type FlowConfiguration struct {
	Agent  AgentFlowConfiguration
	Result AgentFlowConfiguration
}

type AgentFlowConfiguration struct {
	Topic    string
	Producer kafka.KafkaProducerConfiguration
	Consumer kafka.KafkaConsumerConfiguration
}
