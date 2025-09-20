package api

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"

	_ "golang.org/x/sync/errgroup"
)

// Start initializes and starts the API server
func Start(cfg *Config) error {

	db, err := newDB(cfg.DbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	srv := newServer(cfg, sqlstorage.New(db))

	srv.logger.Info("Starting API server", "address", cfg.ServerAddress)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		srv.logger.Info("Shutting down API server")
		if err := srv.shutdown(context.Background()); err != nil {
			srv.logger.Error("Error shutting down server", "error", err)
		}
	}()

	return srv.listenAndServe()
}

func (s *server) shutdown(ctx context.Context) error {
	return nil
}

func newDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
