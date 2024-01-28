package view_protobuf_message_repeated_field

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/routes/mvc"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Model *model.ProtobufMessageRepeatedField
}

func (c *Model) Render(ctx web.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return mvc.ComponentRender[web.Context, *Model](ctx, w, r, Html, c)
}
