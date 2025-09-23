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

	"golang.org/x/sync/errgroup"
)

type api struct {
	db  storage.Storage
	srv *http.Server
	lg  *slog.Logger
}

func New(srvAddress, logLevel string, db storage.Storage) *api {

	a := new(api)

	a.db = db
	a.configureServer(srvAddress)
	a.configureLogger(logLevel)

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
	h := http.NewServeMux()

	h.Handle("/hello", a.handleHello())

	a.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
}

func (a *api) response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			a.error(w, http.StatusInternalServerError, err)
		}
	}
}

func (a *api) error(w http.ResponseWriter, status int, err error) {
	a.lg.Error("failed to encode response", "error", err)
	a.response(w, status, map[string]string{"error": err.Error()})
}
