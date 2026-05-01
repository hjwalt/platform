package page_chat

import (
	"context"
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
	Messages []agent.Message
}

func get(c example.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			html,
			model{
				Messages: []agent.Message{
					{
						Type:    agent.MessageType_User,
						Message: "hello world",
						Raw:     "",
					},
					{
						Type:    agent.MessageType_Agent,
						Message: "hello world back",
						Raw:     "",
					},
					{
						Type:    agent.MessageType_ToolRequest,
						Message: "hello world back",
						Raw:     "",
					},
				},
			},
			map[string]render.View{
				"chats": chat_list.View([]render.View{
					chat_item.View(agent.Message{
						Type:    agent.MessageType_User,
						Message: "hello world",
						Raw:     "",
					}),
					chat_item.View(agent.Message{
						Type:    agent.MessageType_User,
						Message: "hello world",
						Raw:     "",
					}),
				}),
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

		agentMessage := []agent.Message{
			{
				Type:    agent.MessageType_User,
				Message: message[0],
			},
		}

		chatResult, err := c.Chat.Chat(context.Background(), agentMessage)
		if err != nil {
			return nil, err
		}

		results := []render.View{
			chat_item.View(agent.Message{
				Type:    agent.MessageType_User,
				Message: message[0],
				Raw:     "",
			}),
		}

		for _, agentMessage := range chatResult {
			results = append(results, chat_item.View(agentMessage))
		}

		return chat_list.View(results), nil
	} else {
		return chat_list.View([]render.View{
			chat_item.View(agent.Message{
				Type:    agent.MessageType_User,
				Message: "no message received",
				Raw:     "",
			}),
		}), nil
	}
}

func Add(builder route.Builder[example.Context]) {
	builder.Handle(path, http.MethodGet, render.Handle(get, page_error_500.Error))
	builder.Handle(path, http.MethodPost, render.Handle(post, page_error_500.Error))
}
