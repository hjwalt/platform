package web

import (
	"context"
	"net/http"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/routes/htmx"
	"github.com/hjwalt/platform/routes/mvc"
)

type Context struct {
	context.Context
	Schema  *model.ProtobufSchema
	TypeMap map[string]*model.ProtobufType
	Htmx    htmx.HxRequestHeader
}

type DecoratorContext struct {
	Schema  *model.ProtobufSchema
	TypeMap map[string]*model.ProtobufType
}

func (d *DecoratorContext) Decorate(c Context, w http.ResponseWriter, r *http.Request) (Context, error) {
	c.Schema = d.Schema
	c.TypeMap = d.TypeMap
	return c, nil
}

type Component = mvc.Component[Context]

type View = mvc.View[Context]

func DecoratorHtmx(c Context, w http.ResponseWriter, r *http.Request) (Context, error) {
	c.Htmx = htmx.Extract(r)
	return c, nil
}
