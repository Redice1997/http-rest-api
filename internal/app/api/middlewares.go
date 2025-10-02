package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Redice1997/http-rest-api/internal/app/model"
	"github.com/Redice1997/http-rest-api/internal/app/service/user"
	"github.com/google/uuid"
)

func (a *API) mwAuth(next http.Handler) http.Handler {
	var service = user.New(a.db, a.ss, a.lg)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ctx, err := service.Authenticate(r); err != nil {
			a.handleAllErrors(w, r, err)
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (a *API) mwSetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set(model.HeaderRequestID, reqID)
		ctx := context.WithValue(r.Context(), model.CtxRequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) mwLogRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		rw := newRW(w)

		next.ServeHTTP(rw, r)

		a.lg.Info("HTTP request",
			slog.String("request_id", r.Context().Value(model.CtxRequestIDKey).(string)),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
			slog.String("user_agent", r.UserAgent()),
			slog.Int("status", rw.status),
			slog.Duration("duration", time.Since(start)),
			slog.Int("size", rw.size),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func newRW(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusInternalServerError, 0}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(data)
	rw.size = n
	return n, err
}
