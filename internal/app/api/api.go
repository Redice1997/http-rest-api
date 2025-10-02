package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Redice1997/http-rest-api/internal/app/service/user"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"golang.org/x/sync/errgroup"

	_ "net/http/pprof"

	_ "github.com/Redice1997/http-rest-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	cfg  *Config
	user *user.Service
	ss   sessions.Store
	db   storage.Storage
	srv  *http.Server
	lg   *slog.Logger
}

func New(cfg *Config, db storage.Storage) *API {
	a := new(API)

	a.cfg = cfg
	a.db = db
	a.configureLogger()
	a.configureSessionStore()
	a.configureServices()
	a.configureServer()

	return a
}

// Start initializes and starts the API server
func (a *API) Start(ctx context.Context) error {

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		a.lg.Info("Starting API server", "address", a.srv.Addr)
		a.lg.Debug("Swagger URL", "address", fmt.Sprintf("http://localhost%s/swagger/", a.srv.Addr))

		err := a.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		<-egCtx.Done()
		a.lg.Info("Shutting down API server")

		toCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.srv.Shutdown(toCtx); err != nil {
			a.lg.Error("Error shutting down server", "error", err)
			return err
		}

		a.lg.Info("API server stopped gracefully")
		return nil
	})

	return eg.Wait()
}

func (a *API) configureSessionStore() {
	a.ss = sessions.NewCookieStore([]byte(a.cfg.SessionKey))
}

func (a *API) configureServices() {
	a.user = user.New(a.db, a.ss, a.lg)
}

func (a *API) configureLogger() {

	var level slog.Level
	switch a.cfg.LogLevel {
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

func (a *API) configureServer() {
	h := mux.NewRouter()

	// MIDDLEWARES
	h.Use(a.mwSetRequestID)
	h.Use(a.mwLogRequest)
	h.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Cookie"}),
	))

	if a.cfg.LogLevel == "debug" {
		// PPROF
		h.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
		// Swagger UI
		h.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)).Methods(http.MethodGet)
		h.PathPrefix("/swagger/doc.json").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)
	}

	// OPTIONS
	h.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodOptions)

	// API V1
	v1 := h.PathPrefix("/API/v1").Subrouter()
	{
		users := v1.PathPrefix("/users").Subrouter()
		{
			users.HandleFunc("/signup", a.HandleSignUp()).Methods(http.MethodPost)
			users.HandleFunc("/signin", a.HandleSignIn()).Methods(http.MethodPost)

			private := users.NewRoute().Subrouter()
			{
				private.Use(a.mwAuth)
				private.HandleFunc("/me", a.HandleMe()).Methods(http.MethodGet)
			}
		}
	}

	a.srv = &http.Server{
		Addr:    a.cfg.ServerAddress,
		Handler: h,
	}
}

func (a *API) respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (a *API) error(w http.ResponseWriter, r *http.Request, status int, err error) {
	a.lg.Error("Error occurred", "error", err)
	a.errorNoLog(w, r, status, err)
}

func (a *API) errorNoLog(w http.ResponseWriter, r *http.Request, status int, err error) {
	a.respond(w, r, status, &ErrorResponse{
		Error: err.Error(),
	})
}

func (a *API) handleAllErrors(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, user.ErrUserValidation):
		a.errorNoLog(w, r, http.StatusBadRequest, err)
	case errors.Is(err, user.ErrUserUnauthorized):
		a.errorNoLog(w, r, http.StatusUnauthorized, err)
	case errors.Is(err, user.ErrUserExists):
		a.errorNoLog(w, r, http.StatusConflict, err)
	case errors.Is(err, user.ErrInvalidEmailOrPassword):
		a.errorNoLog(w, r, http.StatusUnauthorized, err)
	default:
		a.error(w, r, http.StatusInternalServerError, err)
	}
}
