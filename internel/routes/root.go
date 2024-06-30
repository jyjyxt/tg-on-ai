package routes

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"tg.ai/internel/models"
	"tg.ai/internel/session"
	"tg.ai/internel/views"
)

func RegisterRoutes(router *httptreemux.TreeMux) {
	router.GET("/perpetuals", perpetuals)
	router.GET("/", root)
}

func perpetuals(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	perpetuals, err := models.ReadPerpetuals(r.Context(), "")
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderDataResponse(w, r, views.RenderPerpetuals(perpetuals))
	}
}

func root(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if r.URL.Path != "/" {
		views.RenderDataResponse(w, r, session.NotFoundError(r.Context()))
		return
	}

	views.RenderDataResponse(w, r, map[string]any{"source": "https://github.com/jyjyxt/tg-on-ai"})
}
