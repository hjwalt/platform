package page_base

import (
	"html/template"
	"net/http"

	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/page"
)

const (
	Directory = "page_base"
)

var Html = page.Page(Directory + "/instance.html")

type model struct {
	ComponentUrl string
}

func Page(c web.Context, w http.ResponseWriter, r *http.Request) (*template.Template, model, error) {
	if r.URL.Path == "/" {
		return Html, model{ComponentUrl: "/component/home"}, nil
	}
	return Html, model{ComponentUrl: "/component" + r.URL.Path}, nil
}
