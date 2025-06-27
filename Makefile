# -----------------------------------------------------------------------------
#  🐦 GOOGLE-SEARCH-AGENT – Makefile
# -----------------------------------------------------------------------------

# --- Variables
# The binary name and the Go module path.
BIN_NAME := google-search-agent
MODULE   := github.com/ruslannv/google-search-agent

# Get version from git tags. Fallback for CI environments.
VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo "v0.1.0-dev")

# Tools and directories
GO       := go
DIST_DIR := dist
IMAGE    := $(BIN_NAME):$(VERSION)

# Go build flags.
# -s -w: strip debug information to reduce binary size.
# -X: inject the version into the main.appVersion variable in our Go code.
LDFLAGS := -s -w -X 'main.appVersion=$(VERSION)'

# --- Setup
# Explicitly set the default goal to 'help'.
.DEFAULT_GOAL := help

# Declare all non-file targets as .PHONY for correctness and performance.
.PHONY: help tidy fmt lint test coverage build run docker-build docker-run clean

# -----------------------------------------------------------------------------
# 📖 Help
# Self-documenting Makefile. Targets with '# help:' comments will be shown.
# -----------------------------------------------------------------------------
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "──────────────────────────────────────────────────────────────────"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo "──────────────────────────────────────────────────────────────────"
	@echo ""


# -----------------------------------------------------------------------------
# 🛠️ Development & Formatting
# -----------------------------------------------------------------------------
tidy: ## Tidy and verify Go modules.
	@echo "🧹 Tidying Go modules..."
	@$(GO) mod tidy
	@$(GO) mod verify

fmt: ## Format Go source code.
	@echo "🎨 Formatting code with gofmt and goimports..."
	@$(GO) fmt ./...
	@$(GO) run golang.org/x/tools/cmd/goimports@latest -w .

lint: ## Lint Go source code with golangci-lint.
	@echo "🔎 Linting source code..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Error: golangci-lint is not installed."; \
		echo "Please install it: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi
	@golangci-lint run


# -----------------------------------------------------------------------------
# ✅ Testing
# -----------------------------------------------------------------------------
test: tidy ## Run unit tests.
	@echo "🔬 Running tests..."
	@$(GO) test -v -race -timeout 30s ./...

coverage: ## Run tests and generate an HTML coverage report.
	@echo "📊 Generating test coverage report..."
	@mkdir -p $(DIST_DIR)
	@$(GO) test -coverprofile=$(DIST_DIR)/coverage.out ./...
	@$(GO) tool cover -html=$(DIST_DIR)/coverage.out -o $(DIST_DIR)/coverage.html
	@echo "Coverage report available at: $(DIST_DIR)/coverage.html"


# -----------------------------------------------------------------------------
# 📦 Build & Run
# -----------------------------------------------------------------------------
build: tidy ## Build the Go binary.
	@echo "🏗️  Building binary for $(VERSION)..."
	@mkdir -p $(DIST_DIR)
	@$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BIN_NAME) .
	@echo "Binary available at: $(DIST_DIR)/$(BIN_NAME)"

run: build ## Run the server locally on port 8080 (uses .env).
	@echo "🚀 Starting server on http://localhost:8080"
	@if [ ! -f .env ]; then \
		echo "Warning: .env file not found. Export GOOGLE_API_KEY and GOOGLE_CSE_ID manually."; \
		$(DIST_DIR)/$(BIN_NAME); \
	else \
		export $(shell cat .env | sed 's/#.*//g' | xargs) && \
		$(DIST_DIR)/$(BIN_NAME); \
	fi


# -----------------------------------------------------------------------------
# 🐳 Docker
# -----------------------------------------------------------------------------
docker-build: ## Build the Docker image.
	@echo "🐳 Building Docker image: $(IMAGE)"
	@docker build --build-arg VERSION=$(VERSION) -t $(IMAGE) .

docker-run: ## Run the server in a Docker container (uses .env).
	@echo "🐳 Running Docker container on http://localhost:8080"
	@docker run --rm --env-file .env -p 8080:8080 $(IMAGE)


# -----------------------------------------------------------------------------
# 🧹 Clean
# -----------------------------------------------------------------------------
clean: ## Remove build artifacts.
	@echo "🗑️  Cleaning up..."
	@rm -rf $(DIST_DIR)