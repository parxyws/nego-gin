# Makefile for Docker Compose and Golang Migrations

# Variables
DOCKER_COMPOSE = podman-compose
DOCKER = podman
MIGRATE = migrate
DB_HOST = localhost
DB_PORT = 5540
DB_NAME = nego_db
DB_USER = postgres
DB_PASSWORD = postgres
MIGRATIONS_DIR = ./db/migrations
DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
DB_CONTAINER = master_nego_db

# Colors for output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo -e '$(GREEN)Available commands:$(NC)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# Docker Compose Commands
.PHONY: up
up: ## Start all containers in detached mode
	@echo -e "$(GREEN)Starting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.local.yaml up -d

.PHONY: down
down: ## Stop and remove all containers
	@echo -e "$(YELLOW)Stopping Docker containers...$(NC)"
	$(DOCKER_COMPOSE) down

.PHONY: restart
restart: down up ## Restart all containers

.PHONY: build
build: ## Build or rebuild services
	@echo -e "$(GREEN)Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build

.PHONY: logs
logs: ## View logs from all containers
	$(DOCKER_COMPOSE) logs -f

.PHONY: ps
ps: ## List running containers
	$(DOCKER_COMPOSE) ps

.PHONY: clean
clean: ## Remove all containers, volumes, and images
	@echo -e "$(RED)Cleaning up Docker resources...$(NC)"
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans

# Migration Commands
.PHONY: migrate-create
migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo -e "$(RED)Error: name parameter is required$(NC)"; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@echo -e "$(GREEN)Creating migration: $(name)$(NC)"
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Run all pending migrations
	@echo -e "$(GREEN)Running migrations up...$(NC)"
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@echo -e "$(YELLOW)Rolling back last migration...$(NC)"
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

.PHONY: migrate-down-all
migrate-down-all: ## Rollback all migrations
	@echo -e "$(RED)Rolling back all migrations...$(NC)"
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo -e "$(RED)Error: version parameter is required$(NC)"; \
		echo "Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	@echo -e "$(YELLOW)Forcing migration to version $(version)$(NC)"
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@echo -e "$(GREEN)Current migration version:$(NC)"
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

.PHONY: migrate-status
migrate-status: migrate-version ## Alias for migrate-version

# Combined Commands
.PHONY: init
init: up wait-db migrate-up ## Initialize: start containers, wait for DB, and run migrations
	@echo -e "$(GREEN)Initialization complete!$(NC)"

.PHONY: wait-db
wait-db: ## Wait for database to be ready
	@echo -e "$(YELLOW)Waiting for database to be ready...$(NC)"
	@until $(DOCKER) exec $(DB_CONTAINER) pg_isready -U $(DB_USER) 2>/dev/null; do \
		sleep 1; \
	done
	@echo -e "$(GREEN)Database is ready!$(NC)"

.PHONY: reset-db
reset-db: migrate-down-all migrate-up ## Reset database: rollback all and migrate up

.PHONY: fresh
fresh: down clean up wait-db migrate-up ## Fresh start: clean everything and reinitialize

# Development helpers
.PHONY: shell-db
shell-db: ## Open psql shell in database container
	$(DOCKER) exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: dump-db
dump-db: ## Dump database to file
	@echo -e "$(GREEN)Dumping database...$(NC)"
	$(DOCKER) exec $(DB_CONTAINER) pg_dump -U $(DB_USER) $(DB_NAME) > dump_$$(date +%Y%m%d_%H%M%S).sql

# App Commands
.PHONY: run
run: ## Run the Go application locally
	@echo -e "$(GREEN)Running application...$(NC)"
	go run cmd/nego/main.go

.PHONY: build-go
build-go: ## Compile the Go application
	@echo -e "$(GREEN)Compiling application...$(NC)"
	go build -o tmp/main cmd/nego/main.go

.DEFAULT_GOAL := help