package session

import (
	"context"

	"github.com/unrolled/render"
)

type contextValueKey int

const (
	keyDatabase    contextValueKey = 10
	keyRender      contextValueKey = 11
	keyRequestBody contextValueKey = 12
)

func SqliteDB(ctx context.Context) *SQLite3Store {
	v, _ := ctx.Value(keyDatabase).(*SQLite3Store)
	return v
}

func Render(ctx context.Context) *render.Render {
	v, _ := ctx.Value(keyRender).(*render.Render)
	return v
}

func RequestBody(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestBody).(string)
	return v
}

func WithSqliteDB(ctx context.Context, sqlite *SQLite3Store) context.Context {
	return context.WithValue(ctx, keyDatabase, sqlite)
}

func WithRender(ctx context.Context, render *render.Render) context.Context {
	return context.WithValue(ctx, keyRender, render)
}

func WithRequestBody(ctx context.Context, body string) context.Context {
	return context.WithValue(ctx, keyRequestBody, body)
}
