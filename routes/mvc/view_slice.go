package mvc

import (
	"context"
	"html/template"
	"net/http"
)

type ComponentSliceModel[M any] struct {
	Components []template.HTML
	Model      M
}

func ComponentSliceRender[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, t *template.Template, m M, cs []Component[C]) (template.HTML, error) {
	model, err := ComponentSliceCreateModel(ctx, w, r, m, cs)
	if err != nil {
		return template.HTML(""), err
	}
	return ComponentRender[C, ComponentSliceModel[M]](ctx, w, r, t, model)
}

func ComponentSliceWrite[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, t *template.Template, m M, cs []Component[C]) error {
	model, err := ComponentSliceCreateModel(ctx, w, r, m, cs)
	if err != nil {
		return err
	}
	return ComponentWrite[C, ComponentSliceModel[M]](ctx, w, r, t, model)
}

func ComponentSliceCreateModel[C context.Context, M any](ctx C, w http.ResponseWriter, r *http.Request, m M, cs []Component[C]) (ComponentSliceModel[M], error) {
	components := make([]template.HTML, len(cs))
	for i, com := range cs {
		var renderErr error
		components[i], renderErr = com.Render(ctx, w, r)
		if renderErr != nil {
			return ComponentSliceModel[M]{}, renderErr
		}
	}

	return ComponentSliceModel[M]{
		Model:      m,
		Components: components,
	}, nil
}

// Example of a basic slice component with page and component

type ComponentSlice[C context.Context, M any] struct {
	Template   *template.Template
	Components []Component[C]
	Model      M
}

func (c ComponentSlice[C, M]) Render(ctx C, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	return ComponentSliceRender[C, M](ctx, w, r, c.Template, c.Model, c.Components)
}

func (c ComponentSlice[C, M]) Write(ctx C, w http.ResponseWriter, r *http.Request) error {
	return ComponentSliceWrite[C, M](ctx, w, r, c.Template, c.Model, c.Components)
}
