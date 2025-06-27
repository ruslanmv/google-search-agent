# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
#   🦫 GOOGLE-SEARCH-AGENT – Makefile
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# Author : Your Name
# Usage  : make <target>   or just `make help`
#
# help: 🦫 GOOGLE-SEARCH-AGENT (Go build & automation helpers)
# ─────────────────────────────────────────────────────────────────────────

MODULE       := github.com/ruslannv/google-search-agent
BIN_NAME     := google-search-agent
VERSION      ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo "v0.0.0-dev")

DIST_DIR     := dist
COVERPROFILE := $(DIST_DIR)/coverage.out
COVERHTML    := $(DIST_DIR)/coverage.html

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
# help: tidy    - go mod tidy + verify
# help: fmt     - gofmt & goimports
tidy:
	@$(GO) mod tidy
	@$(GO) mod verify

fmt:
	@$(GO) fmt ./...
	@go run golang.org/x/tools/cmd/goimports@latest -w .

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🔍 Linting & Static Analysis
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: vet lint

# help: 🔍 LINTING
# help: vet     - go vet
# help: lint    - golangci-lint run
vet:
	@$(GO) vet ./...

lint:
	@golangci-lint run

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🧪 Tests & Coverage
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: test coverage

# help: 🧪 TESTS & COVERAGE
# help: test      - Run unit tests
# help: coverage  - HTML coverage report
test:
	@$(GO) test -timeout=60s ./...

coverage:
	@mkdir -p $(DIST_DIR)
	@$(GO) test -covermode=count -coverprofile=$(COVERPROFILE) ./...
	@$(GO) tool cover -html=$(COVERPROFILE) -o $(COVERHTML)
	@echo "HTML coverage → $(COVERHTML)"

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🛠 Build & Run
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: build install release run

# help: 🛠 BUILD & RUN
# help: build    - Build binary into ./dist
# help: install  - go install into GOPATH/bin
# help: release  - Cross-compile for GOOS/GOARCH
# help: run      - Build then run agent on :8080
build: tidy
	@mkdir -p $(DIST_DIR)
	@$(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(DIST_DIR)/$(BIN_NAME) .

install:
	@$(GO) install -trimpath -ldflags '$(LDFLAGS)' .

release:
	@mkdir -p $(DIST_DIR)/$(GOOS)-$(GOARCH)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	  $(GO) build -trimpath -ldflags '$(LDFLAGS)' \
	  -o $(DIST_DIR)/$(GOOS)-$(GOARCH)/$(BIN_NAME) .

run: build
	@echo "Starting $(BIN_NAME) on :8080"
	@GOOGLE_API_KEY=$${GOOGLE_API_KEY} GOOGLE_CSE_ID=$${GOOGLE_CSE_ID} \
	  $(DIST_DIR)/$(BIN_NAME)

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🐳 Docker
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: docker-build docker-run

IMAGE ?= $(BIN_NAME):$(VERSION)

# help: 🐳 DOCKER
# help: docker-build  - Build Docker image
# help: docker-run    - Run container on :8080 with .env
docker-build:
	@docker build --build-arg VERSION=$(VERSION) -t $(IMAGE) .

docker-run: docker-build
	@docker run --rm \
	  --env-file .env \
	  -p 8080:8080 \
	  $(IMAGE)

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🧹 Clean
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
.PHONY: clean

# help: 🧹 CLEANUP
# help: clean  - Remove build & coverage artifacts
clean:
	@rm -rf $(DIST_DIR) $(COVERPROFILE) $(COVERHTML)
	@echo "Cleaned up."
