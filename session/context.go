package session

import (
	"context"
)

type contextValueKey int

const (
	keyDatabase contextValueKey = 10
)

func SqliteDB(ctx context.Context) *SQLite3Store {
	v, _ := ctx.Value(keyDatabase).(*SQLite3Store)
	return v
}

func WithSqliteDB(ctx context.Context, sqlite *SQLite3Store) context.Context {
	return context.WithValue(ctx, keyDatabase, sqlite)
}
