# PHPier CLI Makefile
# Provides convenient targets for building, testing, and installing

# Build variables
BINARY_NAME := phpier
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "dev-$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo 'unknown')
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RED := \033[31m
NC := \033[0m

.PHONY: help build install uninstall clean test fmt vet deps

# Default target
help: ## Show this help message
	@echo "PHPier CLI Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Build information:"
	@echo "  Version: $(YELLOW)$(VERSION)$(NC)"
	@echo "  Commit:  $(YELLOW)$(COMMIT)$(NC)"
	@echo "  Date:    $(YELLOW)$(DATE)$(NC)"

build: ## Build the phpier binary for current platform
	@echo "$(BLUE)Building phpier v$(VERSION)...$(NC)"
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)
	@echo "$(GREEN)✅ Build complete: $(BINARY_NAME)$(NC)"

install: ## Build and install phpier locally using the install script
	@echo "$(BLUE)Installing phpier locally...$(NC)"
	@./scripts/local-install.sh

uninstall: ## Uninstall phpier from the system
	@echo "$(BLUE)Uninstalling phpier...$(NC)"
	@./scripts/local-uninstall.sh

clean: ## Clean build artifacts and temporary files
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f .phpier* index.php docker-compose.yml Dockerfile.php
	@echo "$(GREEN)✅ Clean complete$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✅ Tests complete$(NC)"

fmt: ## Format Go code
	@echo "$(BLUE)Formatting Go code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Formatting complete$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Vet complete$(NC)"

deps: ## Download and tidy dependencies
	@echo "$(BLUE)Managing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(NC)"

# Development shortcuts
quick: fmt vet build ## Quick build with formatting and vetting
	@echo "$(GREEN)✅ Quick build complete$(NC)"