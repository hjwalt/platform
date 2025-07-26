package stateful

import (
	"context"

	"github.com/hjwalt/platform/commons/format"
	"github.com/hjwalt/platform/flows/flow"
)

func ConvertPersistenceId[IK any, IV any](
	source PersistenceIdFunction[IK, IV],
	ik format.Format[IK],
	iv format.Format[IV],
) PersistenceIdFunction[[]byte, []byte] {
	return func(ctx context.Context, m flow.Message[[]byte, []byte]) (string, error) {
		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return "", unmarshalError
		}

		return source(ctx, formattedMessage)
	}
}
