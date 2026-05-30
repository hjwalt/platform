package postgres_store

import (
	"time"

	"github.com/hjwalt/platform/state"
)

type StateTable struct {
	Id                 string `bun:",pk"`
	Content            []byte
	CreatedTimestampMs int64
	UpdatedTimestampMs int64
}

func TableToState(dbState *StateTable) (state.State, error) {
	return state.State{
		Id:        dbState.Id,
		Value:     dbState.Content,
		Timestamp: time.UnixMilli(dbState.UpdatedTimestampMs),
	}, nil
}

func StateToTable(nextState state.State) (*StateTable, error) {
	return &StateTable{
		Id:                 nextState.Id,
		Content:            nextState.Value,
		UpdatedTimestampMs: nextState.Timestamp.UnixMilli(),
	}, nil
}
