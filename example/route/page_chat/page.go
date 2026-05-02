package page_chat

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/example"
	"github.com/hjwalt/platform/example/route/component_sidebar"
	"github.com/hjwalt/platform/example/route/components/chat_item"
	"github.com/hjwalt/platform/example/route/components/chat_list"
	"github.com/hjwalt/platform/example/route/layout"
	"github.com/hjwalt/platform/example/route/page_error_500"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/web/render"
	"github.com/hjwalt/platform/web/route"
)

//go:embed *
var files embed.FS

var html = render.Embedded(files, "page.html")

const (
	path = "/chat"
)

type model struct {
}

func get(c example.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	messages, _ := c.RagStore.GetAll("web")

	messageViews := make([]render.View, 0)
	for _, message := range messages {
		messageViews = append(messageViews, chat_item.View(message))
	}

	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			html,
			model{},
			map[string]render.View{
				"chats": chat_list.View(messageViews),
			},
			[]render.View{},
		)), nil
}

func post(c example.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return nil, err
	}

	slog.Info("", "form", r.PostForm)
	if message, exists := r.PostForm["message"]; exists && len(message) > 0 {

		c.AgentMessageProducer.Produce(c, []flow.Message[agent.Message]{
			{
				Value: agent.Message{
					Context: "web",
					Type:    agent.MessageType_User,
					Message: message[0],
				},
			},
		})

		return chat_list.View([]render.View{}), nil
	} else {
		return chat_list.View([]render.View{
			chat_item.View(agent.Message{
				Context: "web",
				Type:    agent.MessageType_User,
				Message: "no message received",
			}),
		}), nil
	}
}

func Add(builder route.Builder[example.Context]) {
	builder.Handle(path, http.MethodGet, render.Handle(get, page_error_500.Error))
	builder.Handle(path, http.MethodPost, render.Handle(post, page_error_500.Error))
}
