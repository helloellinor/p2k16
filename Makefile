# Go + HTMX P2K16 Makefile

.PHONY: build run test clean dev help

# Build the application
build:
	go build -o p2k16-server ./cmd/server

# Run the application
run: build
	./p2k16-server

# Run in development mode (with auto-restart would require additional tools)
dev:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f p2k16-server

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
	@echo "  build  - Build the application"
	@echo "  run    - Build and run the application"
	@echo "  dev    - Run in development mode"
	@echo "  test   - Run tests"
	@echo "  clean  - Clean build artifacts"
	@echo "  deps   - Download and tidy dependencies"
	@echo "  fmt    - Format code"
	@echo "  lint   - Lint code"
	@echo "  help   - Show this help"