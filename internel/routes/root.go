package routes

import (
	"net/http"

	"tg.ai/internel/session"
	"tg.ai/internel/views"
)

type apiHandler struct{}

func RegisterRoutes(mux *http.ServeMux) {
	h := &apiHandler{}
	mux.HandleFunc("/perpetuals", h.perpetuals)
	mux.HandleFunc("/", h.root)
}

func (apiHandler) perpetuals(w http.ResponseWriter, r *http.Request) {
	views.RenderDataResponse(w, r, map[string]any{"source": "perpetuals"})
}

func (apiHandler) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		views.RenderDataResponse(w, r, session.NotFoundError(r.Context()))
		return
	}

	views.RenderDataResponse(w, r, map[string]any{"source": "https://github.com/jyjyxt/tg-on-ai"})
}
