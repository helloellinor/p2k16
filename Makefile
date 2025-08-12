# Go + HTMX P2K16 Makefile

.PHONY: build run test clean dev demo help

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

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  build-demo - Build demo application"
	@echo "  run        - Build and run the application"
	@echo "  demo       - Build and run demo mode (no database)"
	@echo "  dev        - Run in development mode"
	@echo "  test       - Run tests"
	@echo "  test-auth  - Run authentication test utility"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  help       - Show this help"