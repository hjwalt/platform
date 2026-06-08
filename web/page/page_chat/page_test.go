package page_chat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/harness"
	"github.com/hjwalt/platform/flow"
	"github.com/hjwalt/platform/web"
)

type testStore struct {
	states   map[string]flow.State[harness.ExecutionState]
	readErr  error
	writeErr error
	keysErr  error
}

func (s *testStore) Read(ctx context.Context, id string) (flow.State[harness.ExecutionState], error) {
	if s.readErr != nil {
		return flow.State[harness.ExecutionState]{}, s.readErr
	}
	if st, ok := s.states[id]; ok {
		return st, nil
	}
	return flow.State[harness.ExecutionState]{
		Id: id,
		Value: harness.ExecutionState{
			Context:    id,
			Messages:   make([]agent.Message, 0),
			ToolStates: make(map[string]harness.ToolState),
		},
		Timestamp: time.Now(),
	}, nil
}

func (s *testStore) Write(ctx context.Context, st flow.State[harness.ExecutionState]) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	if s.states == nil {
		s.states = make(map[string]flow.State[harness.ExecutionState])
	}
	s.states[st.Id] = st
	return nil
}

func (s *testStore) Keys(ctx context.Context) ([]string, error) {
	if s.keysErr != nil {
		return nil, s.keysErr
	}
	keys := make([]string, 0, len(s.states))
	for key := range s.states {
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *testStore) Start() error { return nil }

func (s *testStore) Stop() {}

type testProducer struct {
	produced []agent.Message
}

func (p *testProducer) ProduceMessage(ctx context.Context, _ []flow.Message[agent.Message]) error {
	return nil
}

func (p *testProducer) Produce(ctx context.Context, values []agent.Message) error {
	p.produced = append(p.produced, values...)
	return nil
}

func (p *testProducer) Start() error { return nil }

func (p *testProducer) Stop() {}

func newRequestWithChatID(method, target, body, chatID string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("chat_id", chatID)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestPostChatViewRedirectsToSelectedChat(t *testing.T) {
	c := web.Context{Context: context.Background()}
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Set("chat_id", "demo")
	req := httptest.NewRequest(http.MethodPost, "/chat-view", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	view, err := postChatView(c, rr, req)
	if err != nil {
		t.Fatalf("postChatView returned error: %v", err)
	}
	if view != nil {
		t.Fatalf("expected nil view, got non-nil")
	}
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/chat/demo" {
		t.Fatalf("expected redirect location /chat/demo, got %s", got)
	}
}

func TestPostChatViewRedirectsToDefaultWhenMissingChatID(t *testing.T) {
	c := web.Context{Context: context.Background()}
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/chat-view", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := postChatView(c, rr, req)
	if err != nil {
		t.Fatalf("postChatView returned error: %v", err)
	}
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/chat/" {
		t.Fatalf("expected redirect location /chat/, got %s", got)
	}
}

func TestPostPublishesUserMessage(t *testing.T) {
	store := &testStore{states: map[string]flow.State[harness.ExecutionState]{}}
	producer := &testProducer{}
	c := web.Context{
		Context:              context.Background(),
		AgentHarnessStore:    store,
		AgentMessageProducer: producer,
	}
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Set("message", "hello from test")
	req := newRequestWithChatID(http.MethodPost, "/chat/test-chat", form.Encode(), "test-chat")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	view, err := post(c, rr, req)
	if err != nil {
		t.Fatalf("post returned error: %v", err)
	}
	if view == nil {
		t.Fatalf("expected non-nil view")
	}

	if len(producer.produced) != 1 {
		t.Fatalf("expected 1 produced message, got %d", len(producer.produced))
	}
	msg := producer.produced[0]
	if msg.Context != "test-chat" {
		t.Fatalf("expected context test-chat, got %s", msg.Context)
	}
	if msg.Type != agent.MessageType_User {
		t.Fatalf("expected message type %s, got %s", agent.MessageType_User, msg.Type)
	}
	if msg.Message != "hello from test" {
		t.Fatalf("expected message body hello from test, got %s", msg.Message)
	}
}

func TestPostResetClearsChatHistoryAndRedirects(t *testing.T) {
	chatID := "chat-to-reset"
	store := &testStore{states: map[string]flow.State[harness.ExecutionState]{
		chatID: {
			Id: chatID,
			Value: harness.ExecutionState{
				Context: chatID,
				Messages: []agent.Message{
					agent.NewMessage(chatID, agent.MessageType_User, "stale", "", agent.ToolCall{}),
				},
				ToolStates: map[string]harness.ToolState{"tool": harness.ToolState_Executed},
			},
			Timestamp: time.Now(),
		},
	}}
	c := web.Context{
		Context:           context.Background(),
		AgentHarnessStore: store,
	}
	rr := httptest.NewRecorder()

	req := newRequestWithChatID(http.MethodPost, "/chat/chat-to-reset/reset", "", chatID)

	view, err := postReset(c, rr, req)
	if err != nil {
		t.Fatalf("postReset returned error: %v", err)
	}
	if view != nil {
		t.Fatalf("expected nil view, got non-nil")
	}
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/chat/chat-to-reset" {
		t.Fatalf("expected redirect location /chat/chat-to-reset, got %s", got)
	}

	st, ok := store.states[chatID]
	if !ok {
		t.Fatalf("expected state for %s to exist after reset", chatID)
	}
	if st.Value.Context != chatID {
		t.Fatalf("expected context %s, got %s", chatID, st.Value.Context)
	}
	if len(st.Value.Messages) != 0 {
		t.Fatalf("expected 0 messages after reset, got %d", len(st.Value.Messages))
	}
	if len(st.Value.ToolStates) != 0 {
		t.Fatalf("expected 0 tool states after reset, got %d", len(st.Value.ToolStates))
	}
}
