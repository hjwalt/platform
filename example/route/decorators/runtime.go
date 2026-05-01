package decorators

import (
	"net/http"

	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/example"
)

type RuntimeDecorator struct {
	Chat llm.LanguageModel
}

func (d *RuntimeDecorator) Decorate(c example.Context, w http.ResponseWriter, r *http.Request) (example.Context, error) {
	return example.Context{
		Context: c.Context,
		Chat:    d.Chat,
	}, nil
}
