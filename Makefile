APP_NAME := dwimc
BIN_DIR := $(CURDIR)/bin
BUILD_DIR := $(CURDIR)/build
BIN_TOOLS_DIR := $(CURDIR)/bin-tools
SRC := $(CURDIR)/cmd/dwimc/main.go

COVERAGE_PROFILE := $(CURDIR)/coverage.out

DEFAULT_DB := $(CURDIR)/dwimc.db
MIGRATIONS_DIR := $(CURDIR)/migrations
DEP_GOOSE_REPO := https://github.com/pressly/goose.git
DEP_GOOSE_VERSION := v3.24.1


.PHONY: \
	help build clean test lint run \
	update migrate-up migrate-down \
	build-deps clean-deps install-tools clean-tools

help:
	@echo "Available commands:"
	@echo "  make build			- Build the binary"
	@echo "  make clean			- Clean the build directory"
	@echo "  make test			- Run tests with coverage"
	@echo "  make lint			- Run linting"
	@echo "  make run			- Build and run the application"
	@echo "  make update		- Update go modules"
	@echo "  make migrate-up	- Apply local database migrations"
	@echo "  make migrate-down	- Roll back local database migrations"
	@echo "  make build-deps	- Build project dependencies"
	@echo "  make clean-deps	- Clean project dependencies"


default: build

build: build-deps
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
		go build -o $(BIN_DIR)/$(APP_NAME) $(SRC)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)/*

test: build-deps
	@echo "Running tests..."
	MIGRATIONS_DIR=$(MIGRATIONS_DIR) \
		GOOSE_PATH=$(BIN_DIR)/goose \
		go test -v ./... -cover -coverprofile=$(COVERAGE_PROFILE)

lint: install-tools
	@echo "Running linting..."
	@go vet ./...
	@$(BIN_TOOLS_DIR)/golangci-lint run ./...
	@$(BIN_TOOLS_DIR)/modernize ./...

run: build
	@echo "Running $(APP_NAME)..."
	@./$(BIN_DIR)/$(APP_NAME)

update:
	@echo "Updating modules..."
	@go get -u ./...
	@go mod tidy

migrate-up: $(BIN_DIR)/goose
	@echo "Applying database migrations..."
	@./$(BIN_DIR)/goose -dir $(MIGRATIONS_DIR) sqlite3 $(DEFAULT_DB) up

migrate-down: $(BIN_DIR)/goose
	@echo "Rolling back database migrations..."
	@./$(BIN_DIR)/goose -dir $(MIGRATIONS_DIR) sqlite3 $(DEFAULT_DB) down



# Dependencies targets
build-deps: $(BIN_DIR)/goose

clean-deps:
	@echo "Cleaning dependencies..."
	@rm -rf $(BUILD_DIR)/*


$(BUILD_DIR)/goose:
	@echo "Cloning goose..."
	@mkdir -p build
	git -c advice.detachedHead=false clone --depth=1 \
		--branch $(DEP_GOOSE_VERSION) $(DEP_GOOSE_REPO) $(BUILD_DIR)/goose

$(BIN_DIR)/goose: $(BUILD_DIR)/goose
	@echo "Building goose..."
	cd $(BUILD_DIR)/goose && \
		go mod tidy && \
		go build \
			-ldflags="-s -w" \
			-tags='no_postgres no_clickhouse no_mssql no_mysql' \
			-o $(BIN_DIR)/goose ./cmd/goose


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
