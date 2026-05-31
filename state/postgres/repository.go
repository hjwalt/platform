package postgres_store

import (
	"context"
	"database/sql"

	"github.com/hjwalt/platform/state"
)

// implementation
type Repository struct {
	stateTableName string
	connection     BunConnection
}

func (r Repository) Read(ctx context.Context, persistenceId string) (state.State, error) {
	dbState := &StateTable{}

	readErr := r.connection.Db().
		NewSelect().
		Model(dbState).
		ModelTableExpr(r.stateTableName+" AS state_table").
		Where("id = ?", persistenceId).
		Scan(ctx)

	if readErr != nil && readErr != sql.ErrNoRows {
		return state.State{}, readErr
	}

	resultState, convertErr := TableToState(dbState)
	if convertErr != nil {
		return state.State{}, convertErr
	}

	resultState.Id = persistenceId

	return resultState, nil
}

// func (r Repository) GetAll(ctx context.Context, persistenceId []string) (map[string]stateful.State[structure.Bytes], error) {
// 	dbState := []StateTable{}

// 	readErr := r.connection.Db().
// 		NewSelect().
// 		Model(&dbState).
// 		ModelTableExpr(r.stateTableName+" AS state_table").
// 		Where("id IN (?)", bun.In(persistenceId)).
// 		Scan(ctx)

// 	if readErr != nil && readErr != sql.ErrNoRows {
// 		return map[string]stateful.State[structure.Bytes]{}, readErr
// 	}

// 	stateMap := map[string]stateful.State[structure.Bytes]{}
// 	for _, stateEntry := range dbState {
// 		if stateMapped, mapErr := TableToState(&stateEntry); mapErr != nil {
// 			return map[string]stateful.State[structure.Bytes]{}, readErr
// 		} else {
// 			stateMap[stateEntry.Id] = stateMapped
// 		}
// 	}

// 	for _, persistenceIdEntry := range persistenceId {
// 		if _, idPresent := stateMap[persistenceIdEntry]; !idPresent {
// 			stateMap[persistenceIdEntry] = stateful.NewState[structure.Bytes](persistenceIdEntry, []byte{})
// 		}
// 	}

// 	return stateMap, nil
// }

func (r Repository) Write(ctx context.Context, value state.State) error {

	dbState, err := StateToTable(value)
	if err != nil {
		return err
	}

	dbState.Id = value.Id

	_, upsertErr := r.connection.Db().
		NewInsert().
		Model(dbState).
		ModelTableExpr(r.stateTableName + " AS state_table").
		On("CONFLICT (id) DO UPDATE").
		Set("internal = EXCLUDED.internal").
		Set("results = EXCLUDED.results").
		Set("content = EXCLUDED.content").
		Set("created_timestamp_ms = EXCLUDED.created_timestamp_ms").
		Set("updated_timestamp_ms = EXCLUDED.updated_timestamp_ms").
		Exec(ctx)

	return upsertErr
}

// func (r Repository) UpsertAll(ctx context.Context, stateMap map[string]stateful.State[structure.Bytes]) error {

// 	dbStates := []*StateTable{}

// 	for k, v := range stateMap {
// 		dbState, err := StateToTable(v)
// 		if err != nil {
// 			return err
// 		}
// 		dbState.Id = k

// 		dbStates = append(dbStates, dbState)
// 	}

// 	_, upsertErr := r.connection.Db().
// 		NewInsert().
// 		Model(&dbStates).
// 		ModelTableExpr(r.stateTableName + " AS state_table").
// 		On("CONFLICT (id) DO UPDATE").
// 		Set("internal = EXCLUDED.internal").
// 		Set("results = EXCLUDED.results").
// 		Set("content = EXCLUDED.content").
// 		Set("created_timestamp_ms = EXCLUDED.created_timestamp_ms").
// 		Set("updated_timestamp_ms = EXCLUDED.updated_timestamp_ms").
// 		Exec(ctx)

// 	return upsertErr
// }
