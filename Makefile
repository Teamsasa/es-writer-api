.PHONY: all up down seed prune fmt migrate help test test-setup test-repository test-usecase test-handler test-all test-cleanup

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := $(shell tput setaf 2)
RESET  := $(shell tput sgr0)

all: up migrate seed ## Start the application, run migrations and seed the database

up: ## Start the application
	@docker-compose up -d --build

down: ## Stop the application
	@docker-compose down -v

migrate: ## Run database migrations
	@docker-compose run --rm api go run ../migrate/main.go

prune: ## Remove dangling images
	@docker image prune -f

fmt: ## Format all Go code files
	@go fmt ./...

test-setup: ## Setup test environment
	@docker-compose up -d test-db
	@go run app/cmd/migrate_test/main.go

test-repository: test-setup ## Run repository tests
	@go test -v ./app/internal/repository/...

test-usecase: ## Run usecase tests
	@go test -v ./app/internal/usecase/...

test-handler: ## Run handler tests
	@go test -v ./app/internal/handler/...

test-all: ## Run all tests
	@go test -v ./app/internal/repository/... ./app/internal/usecase/... ./app/internal/handler/...

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'
