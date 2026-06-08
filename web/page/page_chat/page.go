package page_chat

import (
	"embed"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow"
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
	path           = "/chat"
	idPath         = "/chat/{chat_id}"
	toolPath       = "/chat/{chat_id}/accept"
	rejectToolPath = "/chat/{chat_id}/reject"
	resetPath      = "/chat/{chat_id}/reset"
	chatViewPath   = "/chat-view"
)

type chat struct {
	Id      string
	Current bool
}

type model struct {
	Context string
	Chats   []chat
}

func view(c web.Context, w http.ResponseWriter, r *http.Request, chatId string) (render.View, error) {
	state, _ := c.AgentHarnessStore.Read(c, chatId)

	messageViews := make([]render.View, 0)
	for _, message := range state.Value.Messages {
		messageViews = append(messageViews, component_chat_item.ViewWithState(message, state.Value))
	}

	chats := make([]chat, 0)

	keys, _ := c.AgentHarnessStore.Keys(c)
	for _, key := range keys {
		chats = append(chats, chat{
			Id:      key,
			Current: key == chatId,
		})
	}
	if len(chats) == 0 {
		chats = append(chats, chat{
			Id:      "web",
			Current: true,
		})
	}

	return layout.Dashboard(
		component_sidebar.View(),
		render.Component(
			html,
			model{
				Context: chatId,
				Chats:   chats,
			},
			map[string]render.View{
				"chats": component_chat_list.View(messageViews),
			},
			[]render.View{},
		)), nil
}

func get(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	return view(c, w, r, "web")
}

func getWithId(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	chatId := chi.URLParam(r, "chat_id")
	return view(c, w, r, chatId)
}

func postChatView(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return nil, err
	}
	if id, exists := r.PostForm["chat_id"]; exists && len(id) > 0 {
		http.Redirect(w, r, "/chat/"+id[0], http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/chat/", http.StatusSeeOther)
	}
	return nil, nil
}

func postReset(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
	chatId := chi.URLParam(r, "chat_id")

	if err := c.AgentHarnessStore.Write(c, flow.State[harness.ExecutionState]{
		Id: chatId,
		Value: harness.ExecutionState{
			Context:    chatId,
			Messages:   make([]agent.Message, 0),
			ToolStates: make(map[string]harness.ToolState),
		},
		Timestamp: time.Now(),
	}); err != nil {
		return nil, err
	}

	http.Redirect(w, r, "/chat/"+chatId, http.StatusSeeOther)
	return nil, nil
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
				"",
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
				"",
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
	toolReasoningContent, _ := r.PostForm["tool_reasoning_content"]

	if contextExists && toolIdExists && toolArgumentsExists && toolNameExists && toolMessageExists {
		c.AgentMessageProducer.Produce(c, []agent.Message{
			agent.NewMessage(
				contextValue[0],
				agent.MessageType_ToolExecute,
				"execution approved to "+toolMessage[0],
				getValue(toolReasoningContent),
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
				getValue(toolReasoningContent),
				agent.ToolCall{},
			)),
		}), nil
	}
}

func rejectTool(c web.Context, w http.ResponseWriter, r *http.Request) (render.View, error) {
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
	toolReasoningContent, _ := r.PostForm["tool_reasoning_content"]

	if contextExists && toolIdExists && toolArgumentsExists && toolNameExists && toolMessageExists {
		c.AgentMessageProducer.Produce(c, []agent.Message{
			agent.NewMessage(
				contextValue[0],
				agent.MessageType_ToolResult,
				"execution rejected to "+toolMessage[0],
				getValue(toolReasoningContent),
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
				getValue(toolReasoningContent),
				agent.ToolCall{},
			)),
		}), nil
	}
}

func getValue(contents []string) string {
	if len(contents) == 0 {
		return ""
	}
	return contents[0]
}

func Add(builder route.Builder[web.Context]) {
	builder.Handle(path, http.MethodGet, render.Handle(get, page_error_500.Error))
	builder.Handle(idPath, http.MethodGet, render.Handle(getWithId, page_error_500.Error))
	builder.Handle(idPath, http.MethodPost, render.Handle(post, page_error_500.Error))
	builder.Handle(toolPath, http.MethodPost, render.Handle(postTool, page_error_500.Error))
	builder.Handle(rejectToolPath, http.MethodPost, render.Handle(rejectTool, page_error_500.Error))
	builder.Handle(resetPath, http.MethodPost, render.Handle(postReset, page_error_500.Error))
	builder.Handle(chatViewPath, http.MethodPost, render.Handle(postChatView, page_error_500.Error))
}
