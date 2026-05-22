BINARY_NAME=superdaw-mcp
VERSION=$(shell cat VERSION.md)
.PHONY: all build clean test fmt help integration-test
all: fmt build test
build:
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) cmd/superdaw/main.go
test:
	go test -v ./...
	$(MAKE) integration-test
integration-test:
	go test -v ./tests/integration/...
fmt:
	go fmt ./...
