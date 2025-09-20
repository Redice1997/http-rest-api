package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	mw "github.com/Redice1997/http-rest-api/internal/app/api/middleware"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

type server struct {
	config  *Config
	logger  *slog.Logger
	router  http.Handler
	storage storage.Storage
}

func newServer(config *Config, storage storage.Storage) *server {

	s := new(server)
	s.config = config
	s.storage = storage
	s.configureLogger()
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	r := http.NewServeMux()

	r.Handle("/hello", mw.Wrap(s.handleHello(), mw.Logging(s.logger)))

	s.router = r
}

func (s *server) configureLogger() {

	var level slog.Level
	switch s.config.LogLevel {
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

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	s.logger = log
}

func (s *server) response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Error("failed to encode response", "error", err)
		}
	}
}

func (s *server) error(w http.ResponseWriter, status int, err error) {
	s.response(w, status, map[string]string{"error": err.Error()})
}

func (s *server) handleHello() http.HandlerFunc {

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.response(w, http.StatusOK, &response{Message: "Hello, World!"})
	}
}
