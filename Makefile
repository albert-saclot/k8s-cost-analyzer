.PHONY: build test clean install help

BINARY_NAME=kcost
GO=go

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  test     - Run all tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install the binary to GOPATH/bin"
	@echo "  help     - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GO) test ./... -v

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)

install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install
