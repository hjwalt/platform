package render_test

import (
	"embed"
	"log/slog"
	"testing"

	"github.com/hjwalt/platform/web/render"
	"github.com/stretchr/testify/assert"
)

//go:embed *
var files embed.FS

func TestComponent_Render(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		component render.View
		want      string
		wantErr   bool
	}{
		{
			name: "basic",
			component: render.Component(
				render.Embedded(files, "html/component_test_sub.html"),
				"hello",
				make(map[string]render.View),
				[]render.View{},
			),
			want:    "<span>hello</span>",
			wantErr: false,
		},
		{
			name: "subcomponent",
			component: render.Component(
				render.Embedded(files, "html/component_test_page.html"),
				"hello",
				map[string]render.View{
					"sub": render.Component(
						render.Embedded(files, "html/component_test_sub.html"),
						"world",
						make(map[string]render.View),
						[]render.View{},
					),
				},
				[]render.View{},
			),
			want:    "hello\n\n<span>world</span>",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			got, gotErr := tt.component.Render()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Render() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Render() succeeded unexpectedly")
			}

			slog.Info("want", "want", got)

			assert.Equal(tt.want, got)
		})
	}
}
