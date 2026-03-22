.PHONY: build clean test install release

VERSION := 1.0.0
BINARY := cli-battery
GO := go
GOFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

# Detect OS and Architecture
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
	OS := darwin
else ifeq ($(UNAME_S),Linux)
	OS := linux
else
	OS := windows
endif

ifeq ($(UNAME_M),arm64)
	ARCH := arm64
else
	ARCH := amd64
endif

# Build the binary
build:
	$(GO) build $(GOFLAGS) -o $(BINARY) .

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BINARY)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BATTERY)-linux-arm64 .

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/

# Run tests
test:
	$(GO) test -v ./...

# Install locally
install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)

# Create release
release: clean build-all
	@echo "Release binaries created in dist/"

# Run the binary
run: build
	./$(BINARY)

# Watch mode
watch: build
	./$(BINARY) --watch

# Show help
help:
	@echo "cli-battery Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the binary for current platform"
	@echo "  make build-all   Build for all platforms"
	@echo "  make clean       Remove build artifacts"
	@echo "  make test        Run tests"
	@echo "  make install     Install to /usr/local/bin"
	@echo "  make release     Create release binaries"
	@echo "  make run         Build and run"
	@echo "  make watch       Build and run in watch mode"
