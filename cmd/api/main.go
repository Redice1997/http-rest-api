package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Redice1997/http-rest-api/internal/app/api"
	"github.com/Redice1997/http-rest-api/internal/app/storage/sqlstorage"
	"gopkg.in/yaml.v3"
)

var (
	configPath = "configs/local.yaml"
)

func init() {
	flag.StringVar(&configPath, "config", configPath, "Path to configuration file")
}

// @title Auth API
// @version 1.0
// @description API with authentication and user management

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func main() {
	flag.Parse()

	cfg := api.NewConfig()
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	db, err := sqlstorage.New(cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if err := api.New(cfg, db).Start(ctx); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
