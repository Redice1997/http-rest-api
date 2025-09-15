package api

import (
	"net/http"

	"github.com/Redice1997/http-rest-api/internal/app/storage/memorystorage"
)

// Start initializes and starts the API server
func Start(cfg *Config) error {
	storage := memorystorage.New()
	logger := newLogger(cfg.LogLevel)
	srv := newServer(logger, storage)

	srv.logger.Info("Starting API server", "address", cfg.ServerAddress)

	return http.ListenAndServe(cfg.ServerAddress, srv)
}
