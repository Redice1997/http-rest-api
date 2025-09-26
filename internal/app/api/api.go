package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Redice1997/http-rest-api/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"golang.org/x/sync/errgroup"
)

type api struct {
	sess sessions.Store
	db   storage.Storage
	srv  *http.Server
	lg   *slog.Logger
}

func New(cfg *Config, db storage.Storage) *api {

	a := new(api)

	a.db = db
	a.sess = sessions.NewCookieStore([]byte(cfg.SessionKey))
	a.configureServer(cfg.ServerAddress)
	a.configureLogger(cfg.LogLevel)

	return a
}

// Start initializes and starts the API server
func (api *api) Start(ctx context.Context) error {

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		api.lg.Info("Starting API server", "address", api.srv.Addr)

		err := api.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()
		api.lg.Info("Shutting down API server")

		toCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := api.srv.Shutdown(toCtx); err != nil {
			api.lg.Error("Error shutting down server", "error", err)
			return err
		}

		api.lg.Info("API server stopped gracefully")
		return nil
	})

	return eg.Wait()
}

func (a *api) configureLogger(logLever string) {

	var level slog.Level
	switch logLever {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	a.lg = logger
}

func (a *api) configureServer(addr string) {
	h := mux.NewRouter()

	h.Use(mux.CORSMethodMiddleware(h))

	v1 := h.PathPrefix("/api/v1").Subrouter()
	{
		auth := v1.PathPrefix("/auth").Subrouter()
		{
			auth.HandleFunc("/users", a.handleUserCreate()).Methods("POST")
			auth.HandleFunc("/sessions", a.handleSessionCreate()).Methods("POST")
		}
	}

	a.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
}

func (a *api) respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (a *api) error(w http.ResponseWriter, r *http.Request, status int, err error) {
	a.lg.Error("Error occurred", "error", err)
	a.errorNoLog(w, r, status, err)
}

func (a *api) errorNoLog(w http.ResponseWriter, r *http.Request, status int, err error) {
	a.respond(w, r, status, map[string]string{"error": err.Error()})
}
