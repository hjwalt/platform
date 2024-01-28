package page_error_500

import (
	"embed"
	"net/http"

	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/web/render"
)

//go:embed *
var files embed.FS

var Html = render.Embedded(files, "page.html")

func Error(c example.Context, w http.ResponseWriter, r *http.Request, err error) render.View {
	return render.Component(
		Html,
		"",
		make(map[string]render.View),
		[]render.View{},
	)
}
