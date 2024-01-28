package metadata

import (
	"github.com/google/uuid"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/format"
)

var MetadataHeaderKey = "FLOW_METADATA"

var Format = format.Json[flow.Metadata]()

func Default() flow.Metadata {
	return flow.Metadata{
		Id:       uuid.New().String(),
		Attempt:  0,
		Sequence: -1,
		Source:   "UNKNOWN",
	}
}
