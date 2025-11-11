# Makefile for gcovr-json-util

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Binary name
BINARY_NAME := gcovr-util

.PHONY: all build clean test install help

all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME) $(VERSION)..."
	go install $(LDFLAGS) .

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying modules..."
	go mod tidy

## version: Show version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
