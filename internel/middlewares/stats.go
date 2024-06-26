package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"tg.ai/internel/session"
	"tg.ai/internel/views"
)

func Stats(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("INFO -- : Started %s '%s'\n", r.Method, r.URL)
		defer func() {
			log.Printf("INFO -- : Completed %s in %fms\n", r.Method, time.Now().Sub(start).Seconds())
		}()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
			return
		}
		if len(body) > 0 {
			log.Printf("INFO -- : Paremeters %s\n", string(body))
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		r = r.WithContext(session.WithRequestBody(r.Context(), string(body)))
		w.Header().Set("X-Build-Info", runtime.Version())
		handler.ServeHTTP(w, r)
	})
}
