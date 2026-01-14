.PHONY: help deps build run dev test test-coverage lint docker-build docker-up docker-down migrate-up migrate-down clean

# Variables
GO=C:\Program Files\Go\bin\go.exe
BINARY_NAME=gin-collection-api
DOCKER_COMPOSE=docker-compose -f docker/docker-compose.yml

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Install Go dependencies
	"$(GO)" mod download
	"$(GO)" mod tidy

build: ## Build the application
	"$(GO)" build -o bin/$(BINARY_NAME) cmd/api/main.go

run: ## Run the application locally
	"$(GO)" run cmd/api/main.go

dev: ## Start Docker development environment
	$(DOCKER_COMPOSE) up -d
	@echo "Development environment started!"
	@echo "MySQL: localhost:3306"
	@echo "Redis: localhost:6379"
	@echo "API: Run 'make run' to start"

test: ## Run tests
	"$(GO)" test -v ./...

test-coverage: ## Run tests with coverage
	"$(GO)" test -cover -coverprofile=coverage.out ./...
	"$(GO)" tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	"$(GO)" vet ./...
	"$(GO)" fmt ./...

docker-build: ## Build Docker image
	docker build -f docker/Dockerfile.api -t $(BINARY_NAME):latest .

docker-up: ## Start all Docker services
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all Docker services
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	docker exec -i gin-mysql mysql -ugin_app -pgin_password gin_collection < internal/infrastructure/database/migrations/001_initial_schema.up.sql
	@echo "Migrations completed!"

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	docker exec -i gin-mysql mysql -ugin_app -pgin_password gin_collection < internal/infrastructure/database/migrations/001_initial_schema.down.sql
	@echo "Rollback completed!"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	"$(GO)" clean

install-tools: ## Install development tools
	"$(GO)" install github.com/golang/mock/mockgen@latest
	"$(GO)" install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

setup: deps ## Initial project setup
	@echo "Setting up project..."
	@cp .env.example .env
	@echo ".env file created - please configure it!"
	@echo "Run 'make dev' to start development environment"

db-shell: ## Open MySQL shell
	docker exec -it gin-mysql mysql -ugin_app -pgin_password gin_collection

redis-cli: ## Open Redis CLI
	docker exec -it gin-redis redis-cli
