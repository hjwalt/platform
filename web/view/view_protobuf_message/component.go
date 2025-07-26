package view_protobuf_message

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/view/view_card"
	"github.com/hjwalt/platform/web/view/view_protobuf_message_field"
	"github.com/hjwalt/platform/web/view/view_text"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Model *model.ProtobufMessage
}

func (c *Model) Render(ctx web.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	fields := []web.Component{}
	for _, m := range c.Model.GetFields() {
		fields = append(fields, &view_protobuf_message_field.Model{Model: m})
	}

	card := view_card.Component{
		Components: []web.Component{
			view_text.Model{Label: c.Model.GetName()},
			mvc.ComponentSlice[web.Context, *model.ProtobufMessage]{
				Template:   Html,
				Model:      c.Model,
				Components: fields,
			},
		},
	}

	return card.Render(ctx, w, r)
}
