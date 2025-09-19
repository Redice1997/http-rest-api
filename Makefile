all: build run
		
setup:
	@go get -u github.com/alecthomas/gometalinter
	@go get -u golang.org/x/tools/cmd/cover
	@go get -u go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@gometalinter --install --update

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
	@migrate -path ./migrations -database "postgres://user:password@localhost:5432/api_db?sslmode=disable" up

test:
	@echo "🧪 Running tests..."
	@go test -v -cover -race -timeout 30s ./...

.PHONY: all setup init lint build migrate test run