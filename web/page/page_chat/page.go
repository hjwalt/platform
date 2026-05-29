package page_chat

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/web"
	"github.com/hjwalt/platform/web/component/component_chat_item"
	"github.com/hjwalt/platform/web/component/component_chat_list"
	"github.com/hjwalt/platform/web/component/component_sidebar"
	"github.com/hjwalt/platform/web/layout"
	"github.com/hjwalt/platform/web/page/page_error_500"
	"github.com/hjwalt/platform/web/render"
	"github.com/hjwalt/platform/web/route"
)

//go:embed *
var files embed.FS

var html = render.Embedded(files, "page.html")

const (
	path     = "/chat"
	idPath   = "/chat/{chat_id}"
	toolPath = "/chat/tool"
)

type model struct {
	Context string
}

func get(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	state, _ := c.AgentHarnessStore.Read(c, "web")

	messageViews := make([]render.View, 0)
	for _, message := range state.Value.Messages {
		messageViews = append(messageViews, component_chat_item.View(message))
	}

	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			html,
			model{
				Context: "web",
			},
			map[string]render.View{
				"chats": component_chat_list.View(messageViews),
			},
			[]render.View{},
		)), nil
}

func getWithId(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	chatId := chi.URLParam(r, "chat_id")

	state, _ := c.AgentHarnessStore.Read(c, chatId)

	messageViews := make([]render.View, 0)
	for _, message := range state.Value.Messages {
		messageViews = append(messageViews, component_chat_item.View(message))
	}

	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			html,
			model{
				Context: chatId,
			},
			map[string]render.View{
				"chats": component_chat_list.View(messageViews),
			},
			[]render.View{},
		)), nil
}

func post(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	chatId := chi.URLParam(r, "chat_id")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return nil, err
	}

	slog.Info("", "form", r.PostForm)
	if message, exists := r.PostForm["message"]; exists && len(message) > 0 {

		c.AgentMessageProducer.Produce(c, []agent.Message{
			agent.NewMessage(
				chatId,
				agent.MessageType_User,
				message[0],
				agent.ToolCall{},
			),
		})

		return component_chat_list.View([]render.View{}), nil
	} else {
		return component_chat_list.View([]render.View{
			component_chat_item.View(agent.NewMessage(
				chatId,
				agent.MessageType_User,
				"no message received",
				agent.ToolCall{},
			)),
		}), nil
	}
}

func postTool(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return nil, err
	}

	slog.Info("", "form", r.PostForm)

	contextValue, contextExists := r.PostForm["context"]
	toolId, toolIdExists := r.PostForm["tool_id"]
	toolArguments, toolArgumentsExists := r.PostForm["tool_arguments"]
	toolName, toolNameExists := r.PostForm["tool_name"]
	toolMessage, toolMessageExists := r.PostForm["tool_message"]

	if contextExists && toolIdExists && toolArgumentsExists && toolNameExists && toolMessageExists {
		c.AgentMessageProducer.Produce(c, []agent.Message{
			agent.NewMessage(
				contextValue[0],
				agent.MessageType_ToolExecute,
				"execution approved to "+toolMessage[0],
				agent.ToolCall{
					Id:        toolId[0],
					Arguments: toolArguments[0],
					Name:      toolName[0],
				},
			),
		})

		return component_chat_list.View([]render.View{}), nil
	} else {
		return component_chat_list.View([]render.View{
			component_chat_item.View(agent.NewMessage(
				"web",
				agent.MessageType_User,
				"no message received",
				agent.ToolCall{},
			)),
		}), nil
	}
}

func Add(builder route.Builder[web.Context]) {
	builder.Handle(path, http.MethodGet, render.Handle(get, page_error_500.Error))
	builder.Handle(idPath, http.MethodGet, render.Handle(getWithId, page_error_500.Error))
	builder.Handle(idPath, http.MethodPost, render.Handle(post, page_error_500.Error))
	builder.Handle(toolPath, http.MethodPost, render.Handle(postTool, page_error_500.Error))
}
