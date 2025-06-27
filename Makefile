# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
#   🐦 GOOGLE-SEARCH-AGENT – Makefile
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

MODULE       := github.com/ruslannv/google-search-agent
BIN_NAME     := google-search-agent
VERSION      ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo "v0.0.0-dev")

DIST_DIR     := dist
IMAGE        := $(BIN_NAME):$(VERSION)
GO           ?= go
GOOS         ?= $(shell $(GO) env GOOS)
GOARCH       ?= $(shell $(GO) env GOARCH)

LDFLAGS      := -s -w -X 'main.appVersion=$(VERSION)'

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 📖 Dynamic help
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: help
help:
	@grep '^# help:' $(firstword $(MAKEFILE_LIST)) | sed 's/^# help: //'

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 📂 Module & Formatting
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: tidy fmt

# help: 📂 MODULE & FORMAT
# help: tidy    - go mod tidy + verify (if go present)
tidy:
	@if command -v $(GO) >/dev/null 2>&1; then \
	  $(GO) mod tidy && $(GO) mod verify; \
	else \
	  echo "warning: 'go' not found, skipping tidy"; \
	fi

# help: fmt     - gofmt & goimports (if go present)
fmt:
	@if command -v $(GO) >/dev/null 2>&1; then \
	  $(GO) fmt ./... && \
	  $(GO) run golang.org/x/tools/cmd/goimports@latest -w .; \
	else \
	  echo "warning: 'go' not found, skipping fmt"; \
	fi

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🛠 Build & Run
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: build run docker-build docker-run clean

# help: 🛠 build  - Build binary (local go) or Docker image fallback
build: tidy
	@if command -v $(GO) >/dev/null 2>&1; then \
	  mkdir -p $(DIST_DIR) && \
	  $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(DIST_DIR)/$(BIN_NAME) .; \
	else \
	  echo "go not found → using Docker build"; \
	  $(MAKE) docker-build; \
	fi

# help: run    - Build & run locally on :8080
run: build
	@if [ -f $(DIST_DIR)/$(BIN_NAME) ]; then \
	  echo "Starting $(BIN_NAME) on :8080"; \
	  GOOGLE_API_KEY=$${GOOGLE_API_KEY} GOOGLE_CSE_ID=$${GOOGLE_CSE_ID} \
	    $(DIST_DIR)/$(BIN_NAME); \
	else \
	  echo "Running via Docker"; \
	  $(MAKE) docker-run; \
	fi

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🐳 Docker
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

# help: docker-build  - Build Docker image
docker-build:
	@docker build --build-arg VERSION=$(VERSION) -t $(IMAGE) .

# help: docker-run    - Run container on :8080 with .env
docker-run:
	@docker run --rm \
	  --env-file .env \
	  -p 8080:8080 \
	  $(IMAGE)

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🧹 Clean
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

# help: clean  - Remove build & coverage artifacts
clean:
	@rm -rf $(DIST_DIR)
	@echo "Cleaned up."
