# Makefile for Go REST API project
# Usage:
#   make setup      # Install dependencies
#   make init       # Initialize the project

all: build run
		
setup: ## Install all the build and lint dependencies
	go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/dep/cmd/dep
	gometalinter --install --update
	@$(MAKE) init

init: ## Initialize the project
	@go mod tidy

build:
	@echo "Building the API server..."
	@go build -o bin/api ./cmd/api

run:
	@echo "Starting the API server..."
	@./bin/api

migrate:
	@echo "Running database migrations..."
	# Add your migration commands here


.PHONY: all setup init lint build migrate test run