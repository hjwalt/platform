package view_protobuf_message_oneof_field

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/view/view_card"
	"github.com/hjwalt/platform/web/view/view_text"
	"github.com/hjwalt/routes/mvc"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Model          *model.ProtobufMessageOneofField
	FieldComponent func(*model.ProtobufMessageField) web.Component
}

func (c *Model) Render(ctx web.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	fields := []web.Component{}
	for _, m := range c.Model.GetFields() {
		fields = append(fields, c.FieldComponent(m))
	}

	card := view_card.Component{
		Components: []web.Component{
			view_text.Model{Label: c.Model.GetName()},
			mvc.ComponentSlice[web.Context, *model.ProtobufMessageOneofField]{
				Template:   Html,
				Model:      c.Model,
				Components: fields,
			},
		},
	}

	return card.Render(ctx, w, r)
}
