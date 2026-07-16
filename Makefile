# Build variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BIN_DIR ?= bin
BIN_NAME ?= $(shell basename $(CURDIR))

# Deploy variables
DEPLOY_HOST ?=
DEPLOY_PATH ?= /opt/src/main
DEPLOY_SUPERVISOR ?=

.DEFAULT_GOAL := help

# Macro: check_required verifies a variable is not empty
define check_required
	@if [ -z "$($1)" ]; then \
		echo "Error: $(2) is not set."; \
		echo "  Usage: make deploy-qa DEPLOY_HOST=<host> DEPLOY_SUPERVISOR=<name> [DEPLOY_PATH=<path>]"; \
		exit 1; \
	fi
endef

.PHONY: help
help: ## help: Show available commands
	@echo "Usage: make [target]"
	@echo ""
	@echo "Development:"
	@grep -E '^[a-zA-Z_-]+:.*?## dev:' $(firstword $(MAKEFILE_LIST)) | sort | \
		awk -F':.*?## dev: ' '{printf "  make %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Build:"
	@grep -E '^[a-zA-Z_-]+:.*?## build:' $(firstword $(MAKEFILE_LIST)) | sort | \
		awk -F':.*?## build: ' '{printf "  make %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Test:"
	@grep -E '^[a-zA-Z_-]+:.*?## test:' $(firstword $(MAKEFILE_LIST)) | sort | \
		awk -F':.*?## test: ' '{printf "  make %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Deploy:"
	@grep -E '^[a-zA-Z_-]+:.*?## deploy:' $(firstword $(MAKEFILE_LIST)) | sort | \
		awk -F':.*?## deploy: ' '{printf "  make %-15s %s\n", $$1, $$2}'

.PHONY: rundev
rundev: ## dev: Start local development server
	@echo "Starting local dev server..."
	@go run main.go

.PHONY: fmt
fmt: ## dev: Format Go source code
	@if command -v gofumpt >/dev/null 2>&1; then \
		echo "Formatting with gofumpt..."; \
		gofumpt -l -w .; \
	else \
		gofmt -l -w .; \
		echo "gofumpt not found, using go fmt. Install: go install mvdan.cc/gofumpt@latest"; \
	fi

.PHONY: lint
lint: ## dev: Run static analysis (go vet)
	@go vet ./...

.PHONY: build
build: ## build: Compile for current platform
	@mkdir -p $(BIN_DIR)
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@go build -o $(BIN_DIR)/$(BIN_NAME) main.go
	@echo "Build complete: $(BIN_DIR)/$(BIN_NAME)"

.PHONY: build-linux
build-linux: ## build: Cross-compile for Linux amd64
	@mkdir -p $(BIN_DIR)
	@echo "Building for linux/amd64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BIN_NAME) main.go
	@echo "Build complete: $(BIN_DIR)/$(BIN_NAME)"

.PHONY: test
test: ## test: Run all unit tests
	@echo "Running tests..."
	@go test -v ./... -count=1 && echo "All tests passed."

.PHONY: clean
clean: ## build: Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "Done."

.PHONY: deploy-qa
deploy-qa: build-linux ## deploy: Build, upload, restart QA server
	$(call check_required,DEPLOY_HOST,DEPLOY_HOST)
	$(call check_required,DEPLOY_SUPERVISOR,DEPLOY_SUPERVISOR)
	@echo "Uploading to $(DEPLOY_HOST):$(DEPLOY_PATH)/ ..."
	@scp -O $(BIN_DIR)/$(BIN_NAME) root@$(DEPLOY_HOST):$(DEPLOY_PATH)/
	@echo "Restarting service $(DEPLOY_SUPERVISOR) on $(DEPLOY_HOST)..."
	@ssh root@$(DEPLOY_HOST) "supervisorctl restart $(DEPLOY_SUPERVISOR)"
	@echo "Deploy complete!"
