.PHONY: dev build test clean

# Development mode (Live Reload)
dev:
	air

# Production build
build:
	go build -o bin/api ./cmd/api

# Run tests
test:
	go test ./...