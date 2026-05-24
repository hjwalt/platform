package agent_tool

import (
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
)

type Definition[REQ any, RES any] interface {
	Name() string
	Description() string
	RequestFormat() format.Format[REQ]
	RequestSchema() *jsonschema.Schema
	DescribeRequest(REQ) string
	ResultFormat() format.Format[RES]
	ResultSchema() *jsonschema.Schema
	DescribeResult(RES) string
	Auto() bool
}

type Sync[REQ any, RES any] interface {
	Definition[REQ, RES]
	Apply(REQ) (RES, error)
}

type SyncWrapper Sync[string, string]

type Container interface {
	AddSync(SyncWrapper)
	GetSync() map[string]SyncWrapper

	// TO DEPRECATE
	AsToolMap() map[string]agent.Tool
}
