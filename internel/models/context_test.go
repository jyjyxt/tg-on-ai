package models

import (
	"context"
	"os"

	"tg.ai/internel/session"
)

const pathTest = "/tmp/test.sqlite3"

func teardownTest(ctx context.Context) {
	err := os.Remove(pathTest)
	if err != nil {
		panic(err)
	}
}

func setup() context.Context {
	store, _ := session.OpenDataSQLite3Store(pathTest)
	ctx := context.Background()
	return session.WithSqliteDB(ctx, store)
}
