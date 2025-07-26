package mvc

import (
	"context"
	"html/template"
	"net/http"
)

type ComponentMapModel[M any] struct {
	Components map[string]template.HTML
	Model      M
}

func ComponentMapRender[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, t *template.Template, m M, cs map[string]Component[C]) (template.HTML, error) {
	model, err := ComponentMapCreateModel(ctx, w, r, m, cs)
	if err != nil {
		return template.HTML(""), err
	}
	return ComponentRender[C, ComponentMapModel[M]](ctx, w, r, t, model)
}

func ComponentMapWrite[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, t *template.Template, m M, cs map[string]Component[C]) error {
	model, err := ComponentMapCreateModel(ctx, w, r, m, cs)
	if err != nil {
		return err
	}
	return ComponentWrite[C, ComponentMapModel[M]](ctx, w, r, t, model)
}

func ComponentMapCreateModel[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, m M, cs map[string]Component[C]) (ComponentMapModel[M], error) {
	components := map[string]template.HTML{}
	for i, com := range cs {
		var renderErr error
		components[i], renderErr = com.Render(ctx, w, r)
		if renderErr != nil {
			return ComponentMapModel[M]{}, renderErr
		}
	}

	return ComponentMapModel[M]{
		Model:      m,
		Components: components,
	}, nil
}

// Example of a basic map component with page write

type ComponentMap[C context.Context, M any] struct {
	Template   *template.Template
	Components map[string]Component[C]
	Model      M
}

func (c ComponentMap[C, M]) Render(ctx C, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return ComponentMapRender[C, M](ctx, w, r, c.Template, c.Model, c.Components)
}

func (c ComponentMap[C, M]) Write(ctx C, w http.ResponseWriter, r *http.Request) error {
	return ComponentMapWrite[C, M](ctx, w, r, c.Template, c.Model, c.Components)
}
