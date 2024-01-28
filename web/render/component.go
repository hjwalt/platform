package render

import (
	"bufio"
	"bytes"
	"html/template"
	"log/slog"
)

func Component[M any](
	template *template.Template,
	model M,
	components map[string]View,
	elements []View,
) View {
	return &component[M]{
		Template:   template,
		Model:      model,
		Components: components,
		Elements:   elements,
	}
}

type component[M any] struct {
	Template   *template.Template
	Model      M
	Components map[string]View
	Elements   []View
}

func (c *component[M]) Render() (string, error) {
	renderedComponents := make(map[string]template.HTML)
	for k, v := range c.Components {
		if rendered, renderErr := v.Render(); renderErr != nil {
			return "", renderErr
		} else {
			renderedComponents[k] = template.HTML(rendered)
		}
	}

	renderedElements := make([]template.HTML, len(c.Elements))
	for i, v := range c.Elements {
		if rendered, renderErr := v.Render(); renderErr != nil {
			return "", renderErr
		} else {
			renderedElements[i] = template.HTML(rendered)
		}
	}

	renderModel := Model[M]{
		Model:      c.Model,
		Components: renderedComponents,
		Elements:   renderedElements,
	}

	byteBuffer := bytes.NewBuffer(make([]byte, 0))
	byteWriter := bufio.NewWriter(byteBuffer)

	if executeErr := c.Template.Execute(byteWriter, renderModel); executeErr != nil {
		slog.Error("err", "err", executeErr)
		return "", executeErr
	}

	if flushErr := byteWriter.Flush(); flushErr != nil {
		slog.Error("err", "err", flushErr)
		return "", flushErr
	}

	return byteBuffer.String(), nil
}

type Model[M any] struct {
	Model      M
	Components map[string]template.HTML
	Elements   []template.HTML
}
