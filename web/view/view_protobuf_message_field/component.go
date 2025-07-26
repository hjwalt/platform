package view_protobuf_message_field

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/view/view_protobuf_message_basic_field"
	"github.com/hjwalt/platform/web/view/view_protobuf_message_oneof_field"
	"github.com/hjwalt/platform/web/view/view_protobuf_message_repeated_field"
)

//go:embed *
var files embed.FS

var Html = template.Must(template.ParseFS(files, "component.html"))

type Model struct {
	Model *model.ProtobufMessageField
}

func (c *Model) Render(ctx web.Context, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	switch cc := c.Model.GetField().(type) {
	case *model.ProtobufMessageField_BasicField:
		viewModel := view_protobuf_message_basic_field.Model{Model: cc.BasicField}
		return viewModel.Render(ctx, w, r)
	case *model.ProtobufMessageField_RepeatedField:
		viewModel := view_protobuf_message_repeated_field.Model{Model: cc.RepeatedField}
		return viewModel.Render(ctx, w, r)
	case *model.ProtobufMessageField_OneofField:
		viewModel := view_protobuf_message_oneof_field.Model{Model: cc.OneofField, FieldComponent: func(pmf *model.ProtobufMessageField) web.Component {
			return &Model{Model: pmf}
		}}
		return viewModel.Render(ctx, w, r)
	}

	return mvc.ComponentRender[web.Context, *Model](ctx, w, r, Html, c)
}
