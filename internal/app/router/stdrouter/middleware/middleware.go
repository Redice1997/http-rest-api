package middleware

import (
	"log/slog"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Wrap(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func Logging(logger *slog.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.Log(r.Context(), slog.LevelInfo, "Request received",
				"method", r.Method,
				"url", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			next(w, r)
		}
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
