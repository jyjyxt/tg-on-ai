package views

import "net/http"

type ResponseView struct {
	Data    any    `json:"data,omitempty"`
	Session any    `json:"session,omitempty"`
	Error   error  `json:"error,omitempty"`
	Prev    string `json:"prev,omitempty"`
	Next    string `json:"next,omitempty"`
}

func RenderDataResponse(w http.ResponseWriter, r *http.Request, view any) {
	rv := ResponseView{
		Data: view,
	}
	session.Render(r.Context()).JSON(w, http.StatusOK, rv)
}

func RenderErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sessionError, ok := err.(*session.Error)
	if !ok {
		sessionError = session.ServerError(r.Context(), err)
	}
	if sessionError.Code == 10001 {
		sessionError.Code = 500
	}
	session.Render(r.Context()).JSON(w, sessionError.Status, ResponseView{Error: sessionError})
}

func RenderBlankResponse(w http.ResponseWriter, r *http.Request) {
	session.Render(r.Context()).JSON(w, http.StatusOK, ResponseView{})
}

func RenderOriginalResponse(w http.ResponseWriter, r *http.Request, view any) {
	session.Render(r.Context()).JSON(w, http.StatusOK, view)
}
