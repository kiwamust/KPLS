.PHONY: lint test coverage build clean help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run

test: ## Run tests
	@echo "Running tests..."
	go test -race -v ./...

coverage: ## Generate coverage report
	@echo "Generating coverage report..."
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | tail -1

build: ## Build the application
	@echo "Building application..."
	go build -v -o kpls/kpls ./main.go
	@echo "Build complete: kpls/kpls"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	rm -f kpls/kpls
	@echo "Clean complete"
