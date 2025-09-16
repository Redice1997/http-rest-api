package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	mw "github.com/Redice1997/http-rest-api/internal/app/api/middleware"
	"github.com/Redice1997/http-rest-api/internal/app/storage"
)

type server struct {
	logger  *slog.Logger
	router  *http.ServeMux
	storage storage.Storage
}

func newServer(logger *slog.Logger, storage storage.Storage) *server {

	s := &server{
		logger:  logger,
		storage: storage,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router = http.NewServeMux()

	s.router.Handle("/hello", mw.Wrap(s.handleHello(), mw.Logging(s.logger)))
}

func (s *server) handleHello() http.HandlerFunc {

	type request struct {
		// Add request fields if necessary
	}

	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		res := response{Message: "Hello, World!"}
		if err := writeJSON(w, http.StatusOK, res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
