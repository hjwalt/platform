package mvc

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
)

type View[C context.Context] interface {
	Write(C, http.ResponseWriter, *http.Request) error
}

type Component[C context.Context] interface {
	Render(C, http.ResponseWriter, *http.Request) (template.HTML, error)
}

func ComponentRender[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, t *template.Template, m M) (template.HTML, error) {
	buf := new(bytes.Buffer)
	templateErr := t.Execute(buf, m)

	if templateErr != nil {
		return template.HTML(""), templateErr
	}

	return template.HTML(buf.String()), nil
}

func ComponentWrite[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, t *template.Template, m M) error {
	return t.Execute(w, m)
}

// Example of a basic component with view

type ComponentBasic[C context.Context, M any] struct {
	Template *template.Template
	Model    M
}

func (c ComponentBasic[C, M]) Render(ctx C, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return ComponentRender[C, M](ctx, w, r, c.Template, c.Model)
}

func (c ComponentBasic[C, M]) Write(ctx C, w http.ResponseWriter, r *http.Request) error {
	return ComponentWrite[C, M](ctx, w, r, c.Template, c.Model)
}
