package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"tg-on-ai/session"
	"tg-on-ai/views"
)

func parseRemoteAddr(remoteAddress string) (string, error) {
	host, _, err := net.SplitHostPort(remoteAddress)
	if err == nil {
		remoteAddress = host
	}
	ip := net.ParseIP(remoteAddress)
	if ip == nil {
		return "", fmt.Errorf("invalid remote address %s", remoteAddress)
	}
	return ip.String(), nil
}

func Constraint(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 0 && !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
			return
		}

		if fwd := r.Header.Get(http.CanonicalHeaderKey("X-Forwarded-For")); fwd != "" {
			s := strings.Index(fwd, ", ")
			if s == -1 {
				s = len(fwd)
			}
			r.RemoteAddr = fwd[:s]
		}

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST,DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "600")
		if r.Method == "OPTIONS" {
			views.RenderBlankResponse(w, r)
		} else {
			handler.ServeHTTP(w, r.WithContext(r.Context()))
		}
	})
}
