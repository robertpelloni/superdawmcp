# Makefile for SuperDAW-Universal-MCP

BINARY_NAME=superdaw-mcp
VERSION=$(shell cat VERSION.md)

.PHONY: all build clean test fmt lint help

all: fmt lint build test

build:
	@echo "Building version $(VERSION)..."
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) cmd/superdaw/main.go

test:
	@echo "Running tests..."
	go test -v ./...
	$(MAKE) integration-test

integration-test:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	go vet ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f pkg/vst/cache.json

help:
	@echo "SuperDAW-MCP Makefile"
	@echo "Targets:"
	@echo "  all     - Format, lint, build, and test"
	@echo "  build   - Compile the binary"
	@echo "  test    - Run unit tests"
	@echo "  fmt     - Format Go source files"
	@echo "  lint    - Run go vet"
	@echo "  clean   - Remove build artifacts"
