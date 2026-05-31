package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	brave_search_web_tool "github.com/hjwalt/platform/agent/tool/brave_search_web"
	shell_tool "github.com/hjwalt/platform/agent/tool/shell"
	agent_skill "github.com/hjwalt/platform/agent/tool/skill"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	"github.com/hjwalt/platform/message/kafka"
	file_store "github.com/hjwalt/platform/state/file"
)

type Configuration struct {
	Tool   ToolConfiguration
	Model  ModelConfiguration
	Server WebServerConfiguration
	Flow   FlowConfiguration
}

type ModelConfiguration struct {
	Parser llm.ModelConfig
	Agent  llm.ModelConfig
}

type ToolConfiguration struct {
	BraveSearch   brave_search_web_tool.Configuration
	Shell         shell_tool.Configuration
	ResearchAgent agent_skill.Configuration
	WebFetch      web_fetch_tool.Configuration
}

type WebServerConfiguration struct {
	Port               int
	StaticResourcePath string
}

type FlowConfiguration struct {
	Store  file_store.Configuration
	Agent  AgentFlowConfiguration
	Result AgentFlowConfiguration
}

type AgentFlowConfiguration struct {
	Topic    string
	Producer kafka.KafkaProducerConfiguration
	Consumer kafka.KafkaConsumerConfiguration
}
