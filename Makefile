.PHONY: help build run dev test clean migrate migrate-down migrate-status seed templ css deps install-tools lint format

# Variables
BINARY_NAME=go-rates
MAIN_PATH=./cmd/api
BUILD_DIR=./bin
MIGRATIONS_DIR=./migrations

ifneq (,$(wildcard .env))
    include .env
    export
endif

# Construct MySQL DSN
DB_DSN := $(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&charset=utf8mb4&loc=Local

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo '${GREEN}Available commands:${NC}'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## build: Build the application
build:
	@echo '${GREEN}Building application...${NC}'
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo '${GREEN}Build complete: $(BUILD_DIR)/$(BINARY_NAME)${NC}'

## run: Build and run the application
run: build
	@echo '${GREEN}Starting application...${NC}'
	@$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run with hot reload (air)
dev:
	@echo '${GREEN}Starting development server with hot reload...${NC}'
	@air

## test: Run tests
test:
	@echo '${GREEN}Running tests...${NC}'
	@go test -v -race -timeout 30s ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo '${GREEN}Running tests with coverage...${NC}'
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo '${GREEN}Coverage report generated: coverage.html${NC}'

## bench: Run benchmarks
bench:
	@echo '${GREEN}Running benchmarks...${NC}'
	@go test -bench=. -benchmem ./...

## migrate: Run database migrations
migrate:
	@echo "Running migrations on $$DB_DRIVER..."
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" up

seed:
	@echo "Running Seeder"
	@go run cmd/seeder/main.go

## migrate-down: Rollback last migration
migrate-down:
	@echo '${YELLOW}Rolling back last migration...${NC}'
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" down
	@echo '${GREEN}Rollback completed${NC}'

## migrate-status: Check migration status
migrate-status:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" status

## migrate-create: Create new migration (usage: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo '${YELLOW}Error: NAME is required. Usage: make migrate-create NAME=migration_name${NC}'; \
		exit 1; \
	fi
	@goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
	@echo '${GREEN}Migration created in $(MIGRATIONS_DIR)${NC}'

## seed: Seed database with sample data
seed:
	@echo '${GREEN}Seeding database...${NC}'
	@go run cmd/seeder/main.go

## templ: Generate templ files
templ:
	@echo '${GREEN}Generating templ files...${NC}'
	@templ generate
	@echo '${GREEN}Templ generation complete${NC}'

## css: Build Tailwind CSS
css:
	@echo '${GREEN}Building Tailwind CSS...${NC}'
	@npx tailwindcss -i ./web/css/input.css -o ./static/css/output.css --minify
	@echo '${GREEN}CSS build complete${NC}'

## css-watch: Watch and rebuild Tailwind CSS
css-watch:
	@echo '${GREEN}Watching Tailwind CSS...${NC}'
	@npx tailwindcss -i ./web/css/input.css -o ./static/css/output.css --watch

## deps: Download Go dependencies
deps:
	@echo '${GREEN}Downloading Go dependencies...${NC}'
	@go mod download
	@go mod tidy
	@echo '${GREEN}Dependencies downloaded${NC}'

## deps-frontend: Install frontend dependencies
deps-frontend:
	@echo '${GREEN}Installing frontend dependencies...${NC}'
	@pnpm install
	@echo '${GREEN}Frontend dependencies installed${NC}'

## install-tools: Install development tools
install-tools:
	@echo '${GREEN}Installing development tools...${NC}'
	@go install github.com/air-verse/air@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo '${GREEN}Tools installed${NC}'

## lint: Run linters
lint:
	@echo '${GREEN}Running linters...${NC}'
	@golangci-lint run --timeout 5m

## format: Format code
format:
	@echo '${GREEN}Formatting code...${NC}'
	@go fmt ./...
	@goimports -w .
	@echo '${GREEN}Code formatted${NC}'

## clean: Clean build artifacts
clean:
	@echo '${GREEN}Cleaning build artifacts...${NC}'
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/
	@rm -rf logs/
	@rm -f coverage.out coverage.html
	@echo '${GREEN}Clean complete${NC}'

## clean-all: Clean everything including database
clean-all: clean
	@echo '${YELLOW}Cleaning database...${NC}'
	@rm -rf data/
	@echo '${GREEN}All cleaned${NC}'

## docker-build: Build Docker image
docker-build:
	@echo '${GREEN}Building Docker image...${NC}'
	@docker build -t $(BINARY_NAME):latest .
	@echo '${GREEN}Docker image built${NC}'

## docker-run: Run Docker container
docker-run:
	@echo '${GREEN}Running Docker container...${NC}'
	@docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

## setup: Initial project setup
setup: deps deps-frontend install-tools migrate templ css
	@echo '${GREEN}Project setup complete!${NC}'
	@echo '${GREEN}Run "make dev" to start development server${NC}'

## all: Build everything
all: clean deps templ css build
	@echo '${GREEN}Build complete!${NC}'

## dev-all: Start all development processes
dev-all:
	@echo '${GREEN}Starting all development processes...${NC}'
	@make -j3 dev css-watch

# Default target
.DEFAULT_GOAL := help