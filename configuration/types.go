package configuration

import (
	"github.com/hjwalt/platform/agent/llm"
	finance_fx_price_tool "github.com/hjwalt/platform/agent/tool/finance_fx_price"
	finance_stock_price_tool "github.com/hjwalt/platform/agent/tool/finance_stock_price"
	linux_shell_tool "github.com/hjwalt/platform/agent/tool/linux_shell"
	memory_tool "github.com/hjwalt/platform/agent/tool/memory"
	web_fetch_tool "github.com/hjwalt/platform/agent/tool/web_fetch"
	web_search_tool "github.com/hjwalt/platform/agent/tool/web_search"
	kafka_integration "github.com/hjwalt/platform/integration/kafka"
	file_store "github.com/hjwalt/platform/state/file"
)

type Configuration struct {
	Tool   ToolConfiguration
	Model  ModelConfiguration
	Server WebServerConfiguration
	Flow   FlowConfiguration
	Store  StoreConfiguration
}

type ModelConfiguration struct {
	Configurations map[string]llm.ModelConfig
	Parser         string
	Agent          string
}

type ToolConfiguration struct {
	Shell      linux_shell_tool.Configuration
	WebFetch   web_fetch_tool.Configuration
	WebSearch  web_search_tool.Configuration
	FxPrice    finance_fx_price_tool.Configuration
	StockPrice finance_stock_price_tool.Configuration
	Memory     []memory_tool.Configuration
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
	Producer kafka_integration.Configuration
	Consumer kafka_integration.Configuration
}

type StoreConfiguration struct {
	Agent  file_store.Configuration
	Memory file_store.Configuration
}
