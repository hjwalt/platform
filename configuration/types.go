package configuration

import (
	brave_search_web_tool "github.com/hjwalt/platform/agent/tool/brave_search_web"
	shell_tool "github.com/hjwalt/platform/agent/tool/shell"
	"github.com/hjwalt/platform/message/kafka"
)

type Configuration struct {
	Tool   ToolConfiguration
	OpenAi OpenAiConfiguration
	Server WebServerConfiguration
	Flow   FlowConfiguration
}

type OpenAiConfiguration struct {
	Model    string
	Endpoint string
	Secret   string
}

type ToolConfiguration struct {
	BraveSearch brave_search_web_tool.Configuration
	Shell       shell_tool.Configuration
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
