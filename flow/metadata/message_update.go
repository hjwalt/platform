package metadata

import (
	"github.com/google/uuid"
	"github.com/hjwalt/platform/flow"
)

func MessageUpdate() flow.MetadataOperation {
	return MessageUpdateOperation{}
}

type MessageUpdateOperation struct {
}

func (r MessageUpdateOperation) OnSuccess(meta flow.Metadata) flow.Metadata {
	return flow.Metadata{
		Id:       uuid.New().String(),
		Group:    meta.Group,
		Attempt:  0,
		Sequence: meta.Sequence,
		Source:   meta.Source,
	}
}

func (r MessageUpdateOperation) OnError(meta flow.Metadata) flow.Metadata {
	return flow.Metadata{
		Id:       uuid.New().String(),
		Group:    meta.Group,
		Attempt:  1,
		Sequence: meta.Sequence,
		Source:   meta.Source,
	}
}
