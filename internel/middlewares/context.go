package middlewares

import (
	"net/http"

	"github.com/unrolled/render"
	"tg.ai/internel/session"
)

func Context(handler http.Handler, db *session.SQLite3Store, render *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := session.WithSqliteDB(r.Context(), db)
		ctx = session.WithRender(ctx, render)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
