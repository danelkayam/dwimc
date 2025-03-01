APP_NAME := dwimc
BIN_DIR := bin
BUILD_DIR := build
SRC := cmd/dwimc/main.go

COVERAGE_PROFILE := coverage.out

DEFAULT_DB := dwimc.db
MIGRATIONS_DIR := migrations
DEP_GOOSE_REPO := https://github.com/pressly/goose.git
DEP_GOOSE_VERSION := v3.24.1

default: build

.PHONY: build build-deps clean test lint run update migrate-up migrate-down help

build: build-deps
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
		go build -o $(BIN_DIR)/$(APP_NAME) $(SRC)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)/*

test: build-deps
	@echo "Running tests..."
	@MIGRATIONS_DIR=`readlink -f $(MIGRATIONS_DIR)` \
		GOOSE_PATH=`readlink -f $(BIN_DIR)/goose` \
		go test -v ./... -cover -coverprofile=$(COVERAGE_PROFILE)

lint:
	@echo "Running linting and vetting..."
	@go vet ./...

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

help:
	@echo "Available commands:"
	@echo "  make build	       - Build the binary"
	@echo "  make clean        - Clean the build directory"
	@echo "  make test         - Run tests with coverage"
	@echo "  make lint         - Run linting and vetting"
	@echo "  make run          - Build and run the application"
	@echo "  make update       - Update Go modules"
	@echo "  make migrate-up   - Apply local database migrations"
	@echo "  make migrate-down - Roll back local database migrations"


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
			-o ../../$(BIN_DIR)/goose ./cmd/goose