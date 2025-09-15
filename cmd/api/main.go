package main

import (
	"flag"
	"log"
	"os"

	"github.com/Redice1997/http-rest-api/internal/app/api"
	"gopkg.in/yaml.v3"
)

var (
	configPath = "configs/local.yaml"
)

func init() {
	flag.StringVar(&configPath, "config", configPath, "Path to configuration file")
}

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

	// Entry point for the API server
	if err := api.Start(cfg); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
