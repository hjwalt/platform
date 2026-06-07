package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	linux_shell_tool "github.com/hjwalt/platform/agent/tool/linux_shell"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	web_search_tool "github.com/hjwalt/platform/agent/tool/web_search"
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
	Configurations map[string]llm.ModelConfig
	Parser         string
	Agent          string
}

type ToolConfiguration struct {
	Shell     linux_shell_tool.Configuration
	WebFetch  web_fetch_tool.Configuration
	WebSearch web_search_tool.Configuration
	FxPrice   finance_fx_price_tool.Configuration
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
