package metadata

import "github.com/hjwalt/platform/flow"

func NoOp() flow.MetadataOperation {
	return NoOpOperation{}
}

type NoOpOperation struct {
}

func (r NoOpOperation) OnSuccess(meta flow.Metadata) flow.Metadata {
	return meta
}

func (r NoOpOperation) OnError(meta flow.Metadata) flow.Metadata {
	return meta
}
