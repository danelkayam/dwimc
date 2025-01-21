APP_NAME := dwimc
BUILD_DIR := bin
SRC := cmd/dwimc/main.go

default: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)/*

test:
	@echo "Running tests..."
	@go test -v ./... -cover

lint:
	@echo "Running linting and vetting..."
	@golangci-lint run || true # Use golangci-lint if installed
	@go vet ./...

run: build
	@echo "Running $(APP_NAME)..."
	@./$(BUILD_DIR)/$(APP_NAME)

update:
	@echo "Updating modules..."
	@go get -u ./...
	@go mod tidy

help:
	@echo "Available commands:"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Clean the build directory"
	@echo "  make test     - Run tests with coverage"
	@echo "  make lint     - Run linting and vetting"
	@echo "  make run      - Build and run the application"
	@echo "  make update   - Update Go modules"
