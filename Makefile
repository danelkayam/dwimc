APP_NAME := dwimc
BIN_DIR := $(CURDIR)/bin
BUILD_DIR := $(CURDIR)/build
BIN_TOOLS_DIR := $(CURDIR)/bin-tools
SRC := $(CURDIR)/cmd/dwimc/main.go

COVERAGE_PROFILE := $(CURDIR)/coverage.out

DEFAULT_DB := $(CURDIR)/dwimc.db

.PHONY: \
	help build clean test lint run update \
	build-deps clean-deps install-tools clean-tools

help:
	@echo "Available commands:"
	@echo "  make build			- Build the binary"
	@echo "  make clean			- Clean the build directory"
	@echo "  make test			- Run tests with coverage"
	@echo "  make lint			- Run linting"
	@echo "  make run			- Build and run the application"
	@echo "  make update		- Update go modules"
	@echo "  make build-deps	- Build project dependencies"
	@echo "  make clean-deps	- Clean project dependencies"


default: build

build:
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o $(BIN_DIR)/$(APP_NAME) $(SRC)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)/*

test:
	@echo "Running tests..."
	go test -v ./... -cover -coverprofile=$(COVERAGE_PROFILE)

lint: install-tools
	@echo "Running linting..."
	@go vet ./...
	@$(BIN_TOOLS_DIR)/golangci-lint run ./...
	@$(BIN_TOOLS_DIR)/modernize ./...

run: build
	@echo "Running $(APP_NAME)..."
	@$(BIN_DIR)/$(APP_NAME)

update:
	@echo "Updating modules..."
	@go get -u ./...
	@go mod tidy


# Tools targets
install-tools: $(BIN_TOOLS_DIR)/golangci-lint $(BIN_TOOLS_DIR)/modernize
	@echo "Tools installed successfully."

clean-tools:
	@echo "Removing installed tools..."
	@rm -rf $(BIN_TOOLS_DIR)/*


$(BIN_TOOLS_DIR):
	@mkdir -p $(BIN_TOOLS_DIR)

$(BIN_TOOLS_DIR)/golangci-lint: | $(BIN_TOOLS_DIR)
	@echo "Installing golangci-lint..."
	@env GOBIN=$(BIN_TOOLS_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(BIN_TOOLS_DIR)/modernize: | $(BIN_TOOLS_DIR)
	@echo "Installing modernize..."
	@env GOBIN=$(BIN_TOOLS_DIR) go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest
