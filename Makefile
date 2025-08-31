.PHONY: help build run test clean docker-up docker-down dev run-examples run-debug run-json run-file

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/app ./cmd

run: ## Run the application
	go run ./cmd

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

dev: ## Run in development mode with hot reload (requires air)
	air

install-dev-tools: ## Install development tools
	go install github.com/air-verse/air@latest

# Logger Examples
run-examples: ## Run logger examples demo
	go run ./cmd/examples

run-debug: ## Run logger examples with debug level
	APP_LOGGER_LEVEL=debug go run ./cmd/examples

run-json: ## Run logger examples with JSON format
	APP_LOGGER_FORMAT=json go run ./cmd/examples

run-file: ## Run logger examples with file output
	APP_LOGGER_PROVIDER=file go run ./cmd/examples

# Docker
docker-up: ## Start services with Docker Compose
	docker-compose up -d

docker-down: ## Stop services with Docker Compose
	docker-compose down

docker-build: ## Build Docker image
	docker build -t boilerplate-go .

docker-run: ## Run Docker container
	docker run -p 8080:8080 boilerplate-go

# Database
migrate-up: ## Run database migrations up
	@echo "Add your migration command here"

migrate-down: ## Run database migrations down
	@echo "Add your migration command here"

seed: ## Seed the database
	@echo "Add your seed command here"

# Code Quality
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

# Dependencies
deps: ## Download dependencies
	go mod download
	go mod tidy
