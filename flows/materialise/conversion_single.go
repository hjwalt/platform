package materialise

import (
	"context"

	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/format"
)

func ConvertOneToOne[IK any, IV any, T any](
	source MapFunction[IK, IV, T],
	ik format.Format[IK],
	iv format.Format[IV],
) MapFunction[structure.Bytes, structure.Bytes, T] {
	return func(ctx context.Context, m flow.Message[[]byte, []byte]) ([]T, error) {
		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]T, 0), unmarshalError
		}
		return source(ctx, formattedMessage)
	}
}
