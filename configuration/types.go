package configuration

type Configuration struct {
	OpenAi      OpenAiConfiguration
	BraveSearch BraveSearchConfiguration
	Server      WebServerConfiguration
}

type OpenAiConfiguration struct {
	Model    string
	Endpoint string
	Secret   string
}

type BraveSearchConfiguration struct {
	BaseUrl string
	Secret  string
}

type WebServerConfiguration struct {
	Port               int
	StaticResourcePath string
}
