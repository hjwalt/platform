package page_error_500

import (
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/routes/example"
	"github.com/hjwalt/platform/routes/mvc"
)

const (
	directory = "page_error_500"
)

var Html = example.Page(directory + "/page.html")

func Error(c example.Context, w http.ResponseWriter, r *http.Request, err error) *template.Template {
	return Html
}

func Controller(c example.Context, w http.ResponseWriter, r *http.Request, err error) mvc.View[example.Context] {
	return mvc.ComponentBasic[example.Context, error]{Template: Html, Model: err}
}
