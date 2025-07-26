package page_error_500

import (
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/routes/mvc"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/page"
)

const (
	directory = "page_error_500"
)

var Html = page.Page(directory + "/page.html")

func Error(c web.Context, w http.ResponseWriter, r *http.Request, err error) *template.Template {
	return Html
}

func Controller(c web.Context, w http.ResponseWriter, r *http.Request, err error) mvc.View[web.Context] {
	return mvc.ComponentBasic[web.Context, error]{Template: Html, Model: err}
}
