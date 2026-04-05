package stateless

import (
	"context"

	"github.com/hjwalt/platform/commons/structure"
	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/format"
)

var (
	base64Format = format.Base64()
	bytesFormat  = format.Bytes()
)

func Base64PersistenceId(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) (string, error) {
	base64Message, conversionErr := flow.Convert(
		m,
		bytesFormat,
		bytesFormat,
		base64Format,
		bytesFormat,
	)

	if conversionErr != nil {
		return "", conversionErr
	}

	return base64Message.Key, nil
}
