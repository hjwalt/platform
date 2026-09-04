package configuration

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hjwalt/platform/agent/llm"
	"github.com/stretchr/testify/assert"
)

// The Configuration structs carry no yaml tags, so format.Yaml[Configuration]()
// (gopkg.in/yaml.v3) keys fields by lowercased field name, e.g. WebFetch ->
// "webfetch", StaticResourcePath -> "staticresourcepath", Configurations ->
// "configurations".
const validYamlFixture = `tool:
  webfetch: {}
model:
  parser: main
  agent: main
  configurations:
    main:
      type: 0
      model: gpt-4o-mini
      endpoint: https://api.openai.com/v1
      secret: sk-test-123
server:
  port: 8080
  staticresourcepath: /static
store:
  agent:
    path: /tmp/agent-store
  memory:
    path: /tmp/memory-store
`

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config fixture: %v", err)
	}
	return path
}

func TestReadValidYaml(t *testing.T) {
	assert := assert.New(t)

	path := writeConfigFile(t, validYamlFixture)

	conf, err := Read(path)

	assert.NoError(err)
	assert.Equal(8080, conf.Server.Port)
	assert.Equal("/static", conf.Server.StaticResourcePath)
	assert.Equal("/tmp/agent-store", conf.Store.Agent.Path)
	assert.Equal("/tmp/memory-store", conf.Store.Memory.Path)
	assert.Equal("main", conf.Model.Agent)
	assert.Equal("main", conf.Model.Parser)

	mainModel, ok := conf.Model.Configurations["main"]
	assert.True(ok)
	assert.Equal(llm.OpenAi, mainModel.Type)
	assert.Equal("gpt-4o-mini", mainModel.Model)
	assert.Equal("https://api.openai.com/v1", mainModel.Endpoint)
	assert.Equal("sk-test-123", mainModel.Secret)

	// re-marshalling a parsed document must reproduce an identical document
	bytes, err := Format.Marshal(conf)
	assert.NoError(err)

	roundTripped, err := Format.Unmarshal(bytes)
	assert.NoError(err)

	roundTripBytes, err := Format.Marshal(roundTripped)
	assert.NoError(err)
	assert.Equal(string(bytes), string(roundTripBytes))
	assert.Equal(conf.Server.Port, roundTripped.Server.Port)
	assert.Equal(conf.Store.Agent.Path, roundTripped.Store.Agent.Path)
	assert.Equal(conf.Model.Configurations["main"], roundTripped.Model.Configurations["main"])
}

func TestReadMissingFileReturnsDefaultAndErrReadFail(t *testing.T) {
	assert := assert.New(t)

	path := filepath.Join(t.TempDir(), "does-not-exist.yaml")

	conf, err := Read(path)

	assert.Error(err)
	assert.True(errors.Is(err, ErrReadFail))
	assert.Equal(Configuration{}, conf)
}

func TestReadInvalidYamlReturnsError(t *testing.T) {
	assert := assert.New(t)

	// unterminated quoted scalar is a yaml syntax error
	path := writeConfigFile(t, "server: \"unterminated")

	conf, err := Read(path)

	assert.Error(err)
	assert.Equal(0, conf.Server.Port)
}
