.PHONY: help build run test test-docker clean stop logs

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Build the application
build: ## Build the Docker image
	docker-compose build

# Run the application
run: ## Run the application with Docker Compose
	docker-compose up --build

# Run in background
run-bg: ## Run the application in background
	docker-compose up --build -d

# Run tests locally
test: ## Run tests locally
	go test -v ./...

# Run tests with Docker
test-docker: ## Run tests using Docker Compose
	docker-compose -f compose.test.yml up --build --abort-on-container-exit

# Clean up
clean: ## Stop and remove containers, networks, and volumes
	docker-compose down -v
	docker-compose -f compose.test.yml down -v

# Stop services
stop: ## Stop running containers
	docker-compose down
	docker-compose -f compose.test.yml down

# Show logs
logs: ## Show application logs
	docker-compose logs -f app

# Show test logs
logs-test: ## Show test logs
	docker-compose -f compose.test.yml logs -f test

# Development setup
dev-setup: ## Setup development environment
	go mod download
	cp .env.example .env

# Lint (if golangci-lint is available)
lint: ## Run linter
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt: ## Format Go code
	go fmt ./...

# Vet code
vet: ## Run go vet
	go vet ./...