package models

import (
	"context"

	"tg.ai/internel/session"
)

const pathTest = "/tmp/test.sqlite3"

func setup() context.Context {
	store, _ := session.OpenDataSQLite3Store(pathTest)
	ctx := context.Background()
	return session.WithSqliteDB(ctx, store)
}
