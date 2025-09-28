package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Redice1997/http-rest-api/internal/app/storage"
	"github.com/google/uuid"
)

func (a *api) mwAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.sess.Get(r, SessionName)
		if err != nil {
			a.error(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		}
		id, ok := session.Values["user_id"].(int64)
		if !ok {
			a.error(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		ctx := r.Context()

		u, err := a.db.User().GetByID(ctx, id)
		if errors.Is(err, storage.ErrRecordNotFound) {
			a.error(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		} else if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, CtxUserKey, u)))
	})
}

func (a *api) mwSetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), CtxRequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *api) mwLogRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		rw := newRW(w)

		next.ServeHTTP(rw, r)

		a.lg.Info("HTTP request",
			slog.String("request_id", r.Context().Value(CtxRequestIDKey).(string)),
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
