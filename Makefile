all: docs build run

setup:
	@go install golang.org/x/tools/cmd/cover
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

init:
	@go mod tidy

build:
	@echo "🔨 Building the API server..."
	@go build -o bin/api ./cmd/api

run:
	@echo "🚀 Starting the API server..."
	@./bin/api

migrate:
	@echo "🗄️ Running database migrations..."
	@migrate -path ./migrations -database "postgres://api:password@localhost:5432/api_db?sslmode=disable" up

migrate_test:
	@echo "🗄️ Running test database migrations..."
	@migrate -path ./migrations -database "postgres://api:password@localhost:5432/test_api_db?sslmode=disable" up

test:
	@echo "🧪 Running tests..."
	@go test -v -cover -race -timeout 30s ./...

compose:
	@echo "🐳 Starting the API server with Docker Compose..."
	@docker-compose up -d --scale api=2

lint:
	@echo "Linting the code..."
	@golangci-lint run

docs:
	@echo "Generating API documentation..."
	@swag init -g ./cmd/api/main.go --output docs/

.PHONY: all setup init lint build migrate test run compose lint docs
