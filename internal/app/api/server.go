package api

import (
	"log/slog"
	"net/http"

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
	// Here you would typically set up your routes and handlers
	// For simplicity, we'll just set up a simple health check endpoint
}
