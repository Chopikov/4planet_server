.PHONY: help build run test clean docker-build docker-run docker-stop migrate seed

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Start services with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo "  migrate      - Run database migrations"
	@echo "  seed         - Seed database with initial data"
	@echo "  deps         - Download Go dependencies"
	@echo "  fmt          - Format Go code"
	@echo "  lint         - Lint Go code"
	@echo ""
	@echo "Database Configuration:"
	@echo "  Set DB_DSN environment variable to override default connection:"
	@echo "    export DB_DSN='postgres://user:pass@host:5432/db?sslmode=require'"
	@echo "    make migrate     # Uses custom DB_DSN"
	@echo "    make migrate     # Uses default: $(DB_DSN)"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api ./cmd/api

# Run the application locally
run: build
	@echo "Running application..."
	./bin/api

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Download Go dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	goimports -w .

# Lint Go code
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t 4planet-backend .

# Start services with Docker Compose
docker-run:
	@echo "Starting services..."
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	@echo "Stopping services..."
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f

# Run database migrations
# Database connection string - defaults to local development
DB_DSN ?= postgres://postgres:postgres@localhost:5432/planet?sslmode=disable

migrate:
	@echo "Running migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		if [ -f "migrations/000001_init.up.sql" ]; then \
			migrate -path migrations -database "$(DB_DSN)" up; \
		else \
			echo "No migration files found. Please check migrations directory."; \
		fi; \
	else \
		echo "golang-migrate not found. Installing..."; \
		go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		echo "Please run 'make migrate' again."; \
	fi

migrate-status:
	@echo "Migration status:"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_DSN)" version; \
	else \
		echo "golang-migrate not found. Run 'make deps' first."; \
	fi

migrate-down:
	@echo "Rolling back last migration..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_DSN)" down 1; \
	else \
		echo "golang-migrate not found. Run 'make deps' first."; \
	fi

migrate-reset:
	@echo "Resetting all migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DB_DSN)" down; \
	else \
		echo "golang-migrate not found. Run 'make deps' first."; \
	fi

migrate-create:
	@echo "Creating new migration..."
	@if command -v migrate >/dev/null 2>&1; then \
		read -p "Enter migration name: " name; \
		migrate create -ext sql -dir migrations -seq $$name; \
	else \
		echo "golang-migrate not found. Run 'make deps' first."; \
	fi

migrate-validate:
	@echo "Validating migration files..."
	@./scripts/validate-migrations.sh

# Seed database with initial data
seed:
	@echo "Seeding database..."
	@if [ -f "cmd/seed/main.go" ]; then \
		if command -v go >/dev/null 2>&1; then \
			go run ./cmd/seed; \
		else \
			echo "Go not found. Please install Go."; \
		fi; \
	else \
		echo "No Go seeder found at cmd/seed/main.go"; \
		echo "Create cmd/seed/main.go with your seeding logic or run migrations only."; \
	fi

# Development setup
dev-setup: deps docker-run
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Running migrations..."
	@make migrate
	@echo "Seeding database..."
	@make seed
	@echo "Development environment is ready!"
	@echo "API: http://localhost:8080"
	@echo "Admin: http://localhost:8080/admin (admin/admin)"
	@echo "MailHog: http://localhost:8025"

# Production build
prod-build:
	@echo "Building production binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/api ./cmd/api

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
