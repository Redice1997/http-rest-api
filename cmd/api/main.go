package main

import (
	"log"

	"github.com/Redice1997/http-rest-api/internal/app/api"
)

func main() {
	// Entry point for the API server
	if err := api.Run(); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
