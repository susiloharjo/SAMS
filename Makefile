# SAMS Backend Makefile

.PHONY: help build run test clean deps docker-up docker-down

# Default target
help:
	@echo "SAMS Backend - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  build     - Build the Go application"
	@echo "  run       - Run the application"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps      - Download Go dependencies"
	@echo "  deps-tidy - Tidy Go dependencies"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up   - Start all Docker services"
	@echo "  docker-down - Stop all Docker services"
	@echo "  docker-logs - View Docker logs"
	@echo ""
	@echo "Frontend:"
	@echo "  frontend-install - Install frontend dependencies"
	@echo "  frontend-build   - Build frontend for production"
	@echo "  frontend-dev     - Start frontend development server"
	@echo ""
	@echo "Database:"
	@echo "  db-reset   - Reset database (drop and recreate)"
	@echo "  db-migrate - Run database migrations"
	@echo ""
	@echo "Setup:"
	@echo "  setup      - Complete setup (deps + build + docker-up)"

# Development commands
build:
	@echo "Building SAMS Backend..."
	cd backend && go build -o bin/sams-backend ./cmd/main.go

run: build
	@echo "Running SAMS Backend..."
	cd backend && ./bin/sams-backend

test:
	@echo "Running tests..."
	cd backend && go test ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin/
	rm -rf backend/tmp/

# Dependency management
deps:
	@echo "Downloading Go dependencies..."
	cd backend && go mod download

deps-tidy:
	@echo "Tidying Go dependencies..."
	cd backend && go mod tidy

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs:
	@echo "Viewing Docker logs..."
	docker-compose logs -f

# Frontend commands
frontend-build:
	@echo "Building frontend..."
	cd frontend && npm run build

frontend-dev:
	@echo "Starting frontend development server..."
	cd frontend && npm run dev

frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Database commands
db-reset:
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	sleep 10
	docker-compose exec postgres psql -U $$DB_USER -d $$DB_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	docker-compose exec postgres psql -U $$DB_USER -d $$DB_NAME -f /docker-entrypoint-initdb.d/init-db.sql

db-migrate:
	@echo "Running database migrations..."
	cd backend && go run cmd/main.go

# Complete setup
setup: deps build docker-up
	@echo "Setup complete! SAMS is ready."
	@echo ""
	@echo "Services running:"
	@echo "  - Frontend: http://localhost:3000"
	@echo "  - Backend: http://localhost:8080"
	@echo "  - Database: localhost:5433"
	@echo "  - pgAdmin: http://localhost:5051"
	@echo "  - Redis: localhost:6380"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Copy env.example to .env and configure API keys"
	@echo "  2. Run 'make docker-up' to start all services"
	@echo "  3. Access the frontend at http://localhost:3000"

# Development with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	cd backend && air

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Linting
lint:
	@echo "Running linter..."
	cd backend && golangci-lint run

# Format code
fmt:
	@echo "Formatting Go code..."
	cd backend && go fmt ./...
	cd backend && go vet ./...
