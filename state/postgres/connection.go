package postgres_store

import "github.com/uptrace/bun"

type BunConnection interface {
	Start() error
	Stop()
	Db() bun.IDB
}
