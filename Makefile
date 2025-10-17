.PHONY: help build clean test run lint fmt install-deps

# Variables
BINARY_NAME=gar
CMD_PATH=./cmd/gar
BUILD_DIR=./bin
VERSION=1.0.0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)

# Default target
help:
	@echo "GoArchive (gar) - Build and Development Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make build              Build the application for current platform"
	@echo "  make build-all          Build for all platforms (Windows, Linux, macOS)"
	@echo "  make run                Build and run the application"
	@echo "  make test               Run all tests with coverage"
	@echo "  make clean              Remove build artifacts"
	@echo "  make lint               Run golangci-lint"
	@echo "  make fmt                Format code with gofmt"
	@echo "  make install-deps       Install development dependencies"
	@echo ""

# Install Go dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully"

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	@cd $(CMD_PATH) && go build -o ../../$(BINARY_PATH) -ldflags="-X github.com/cubetiqlabs/gar/pkg/version.Version=$(VERSION)" .
	@echo "Build complete: $(BINARY_PATH)"

# Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "All builds complete in $(BUILD_DIR)/"

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -ldflags="-X github.com/cubetiqlabs/gar/pkg/version.Version=$(VERSION)" $(CMD_PATH)

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -ldflags="-X github.com/cubetiqlabs/gar/pkg/version.Version=$(VERSION)" $(CMD_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -ldflags="-X github.com/cubetiqlabs/gar/pkg/version.Version=$(VERSION)" $(CMD_PATH)

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -ldflags="-X github.com/cubetiqlabs/gar/pkg/version.Version=$(VERSION)" $(CMD_PATH)

# Run the application
run: build
	@$(BINARY_PATH)

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1

# Run tests with HTML coverage report
test-html: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint code
lint:
	@echo "Running linter..."
	@go fmt ./...
	@go vet ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install locally (adds to GOPATH/bin)
install: build
	@mkdir -p $(HOME)/go/bin
	@cp $(BINARY_PATH) $(HOME)/go/bin/$(BINARY_NAME)
	@echo "Installed to $(HOME)/go/bin/$(BINARY_NAME)"

# Show version
version:
	@echo "GoArchive version $(VERSION)"

# Watch mode (requires entr or similar)
watch:
	@find . -type f -name "*.go" | entr make test build

# Development setup
dev-setup: install-deps lint test
	@echo "Development environment ready"
