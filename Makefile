BINARY_NAME=superdaw-mcp
VERSION=$(shell cat VERSION.md)

.PHONY: all build clean test fmt lint help integration-test package

all: fmt build test

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
	go fmt ./...

lint:
	go vet ./...

package: build
	@echo "Packaging distribution..."
	mkdir -p dist/superdaw-$(VERSION)
	cp bin/$(BINARY_NAME) dist/superdaw-$(VERSION)/
	cp scripts/install_adapters.sh dist/superdaw-$(VERSION)/
	cp -r pkg/agents dist/superdaw-$(VERSION)/
	tar -czf superdaw-$(VERSION).tar.gz -C dist superdaw-$(VERSION)
	@echo "Package created: superdaw-$(VERSION).tar.gz"

clean:
	rm -rf bin/
	rm -rf dist/
	rm -f *.tar.gz
	rm -f vst_cache.json
