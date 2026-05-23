BINARY_NAME=superdaw-mcp
VERSION=$(shell cat VERSION.md)
all: build test
build:
	mkdir -p bin; go build -o bin/$(BINARY_NAME) cmd/superdaw/main.go
test:
	go test -v ./...
	$(MAKE) integration-test
integration-test:
	go test -v ./tests/integration/...
