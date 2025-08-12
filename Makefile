# Go + HTMX P2K16 Makefile

.PHONY: build run test clean dev demo dev-python dev-migration help

# Defaults for runtime configuration (override on the command line)
DB_HOST ?= localhost
DB_PORT ?= 2016
DB_USER ?= p2k16-web
DB_PASSWORD ?= p2k16-web
DB_NAME ?= p2k16
PORT ?= 8080

# Build the application
build:
	go build -o p2k16-server ./cmd/server

# Build demo application
build-demo:
	go build -o p2k16-demo ./cmd/demo

# Run the application
run: build
	DB_HOST=$(DB_HOST) \
	DB_PORT=$(DB_PORT) \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_NAME=$(DB_NAME) \
	PORT=$(PORT) \
	./p2k16-server

# Run demo mode (no database required)
demo: build-demo
	./p2k16-demo

# Run in development mode (with auto-restart would require additional tools)
dev:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Run test utility
test-auth:
	go run ./cmd/test

# Clean build artifacts
clean:
	rm -f p2k16-server p2k16-demo

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Migration support targets

# Run Python Flask development server
dev-python:
	@echo "Starting Python Flask server on :5000..."
	@echo "Make sure you've run: source .settings.fish && p2k16-run-web"
	@echo "This target is for documentation - run the command above manually"

# Run both systems for migration testing
dev-migration:
	@echo "=== Migration Development Setup ==="
	@echo "This will help you run both systems in parallel for testing"
	@echo ""
	@echo "Terminal 1 (Python Flask on :5000):"
	@echo "  source .settings.fish"
	@echo "  p2k16-run-web"
	@echo ""
	@echo "Terminal 2 (Go server on :8081):"
	@echo "  make run PORT=8081"
	@echo ""
	@echo "Then test both systems:"
	@echo "  Python: http://localhost:5000"
	@echo "  Go:     http://localhost:8081"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-demo    - Build demo application"
	@echo "  run           - Build and run the application"
	@echo "  demo          - Build and run demo mode (no database)"
	@echo "  dev           - Run in development mode"
	@echo "  dev-python    - Instructions for Python development server"
	@echo "  dev-migration - Instructions for parallel development setup"
	@echo "  test          - Run tests"
	@echo "  test-auth     - Run authentication test utility"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help"