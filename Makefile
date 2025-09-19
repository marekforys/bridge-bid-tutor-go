.PHONY: test test-cover lint lint-fix build clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=bridge-bid-tutor-go

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-cover:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run --fix=false

# Run linter and fix issues
lint-fix:
	golangci-lint run --fix

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server
	./$(BINARY_NAME)
