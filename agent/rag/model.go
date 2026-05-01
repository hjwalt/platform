package rag

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/reflect"
)

const (
	IdContextKey = "RAG_ID_KEY"
)

func Rag(model agent.LanguageModel, store Store) agent.LanguageModel {
	return &RagModel{
		Store:    store,
		Delegate: model,
	}
}

type RagModel struct {
	Store    Store
	Delegate agent.LanguageModel
}

func (r *RagModel) Start() error {
	return nil
}

func (r *RagModel) Stop() {
}

func (r *RagModel) Chat(ctx context.Context, messages []agent.Message) ([]agent.Message, error) {
	id := reflect.GetString(ctx.Value(IdContextKey))
	if id == "" {
		id = "DEFAULT"
	}

	storedMessages, err := r.Store.GetAll(id)
	if err != nil {
		return []agent.Message{}, err
	}

	if storeErr := r.Store.Add(id, messages); storeErr != nil {
		return []agent.Message{}, storeErr
	}

	allmessages := make([]agent.Message, 0)
	allmessages = append(allmessages, storedMessages...)
	allmessages = append(allmessages, messages...)

	resultMessages, err := r.Delegate.Chat(context.Background(), allmessages)
	if err != nil {
		return []agent.Message{}, err
	}

	if storeErr := r.Store.Add(id, resultMessages); storeErr != nil {
		return []agent.Message{}, storeErr
	}

	return resultMessages, nil
}
