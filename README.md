
# 🐦 Google Search Agent

> Author: ruslanmv.com  
> A minimal Go-based MCP agent that exposes a `google_search` tool, powered by the Google Custom Search API.

[![Go Version](https://img.shields.io/badge/go-1.23-blue)]()
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue)]()

---

## Features

- Implements the **`google_search`** MCP tool  
- Queries the Google Custom Search API (requires API key + CSE ID)  
- HTTP (JSON-RPC 2.0) transport at `/`  
- **Health** (`/health`) and **Version** (`/version`) endpoints  
- Single static binary (~2 MiB)  
- Dockerfile for a lightweight container  
- Unit tests with mocked HTTP responses  
- Makefile for build, test, lint, and Docker targets  

---

## Quick Start

```bash
git clone https://github.com/ruslannv/google-search-agent.git
cd google-search-agent

# Build the binary
make build

# Set your credentials
export GOOGLE_API_KEY=your_api_key
export GOOGLE_CSE_ID=your_cse_id

# Run locally on port 8080
make run
````

---

## Installation

Requires Go 1.23+:

```bash
go install github.com/ruslannv/google-search-agent@latest
```

---

## Configuration

All settings are via environment variables:

| Variable         | Description                                   |
| ---------------- | --------------------------------------------- |
| `GOOGLE_API_KEY` | Your Google API key for Custom Search API     |
| `GOOGLE_CSE_ID`  | Your Custom Search Engine ID (the `cx` value) |

You can also override defaults with flags:

```bash
google-search-agent -listen=0.0.0.0 -port=8080
```

* `-listen` (default `0.0.0.0`) — interface to bind
* `-port`   (default `8080`)    — port to listen on
* `-help`                        — show usage

---

## API Reference

### Health Check

**GET** `/health`

```bash
curl http://localhost:8080/health
# → {"status":"ok"}
```

### Version Info

**GET** `/version`

```bash
curl http://localhost:8080/version
# → {"name":"google-search-agent","version":"0.1.0"}
```

### Perform a Search

**POST** `/`
Content-Type: `application/json`

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "google_search",
    "arguments": {
      "query": "IBM WatsonX"
    }
  },
  "id": 1
}
```

**Example response**:

```json
{
  "jsonrpc": "2.0",
  "result": [
    {
      "title": "IBM WatsonX",
      "link": "https://www.ibm.com/cloud/watsonx"
    },
    {
      "title": "Introducing Watsonx",
      "link": "https://www.ibm.com/cloud/watsonx-introducing"
    }
  ],
  "id": 1
}
```

---

## Docker

Build and run via the Makefile (it loads `.env` for your credentials):

```bash
make docker-build
make docker-run
```

Under the hood this runs:

```bash
docker build -t google-search-agent:latest .
docker run --rm --env-file .env -p 8080:8080 google-search-agent:latest
```

---

## Makefile Targets

| Target              | Description                           |
| ------------------- | ------------------------------------- |
| `make help`         | Show this help summary                |
| `make build`        | Compile the binary into `dist/`       |
| `make run`          | Build & run locally on port 8080      |
| `make test`         | Run unit tests                        |
| `make lint`         | Run `golangci-lint`                   |
| `make tidy`         | Run `go mod tidy` + `go mod verify`   |
| `make docker-build` | Build the Docker image                |
| `make docker-run`   | Run the Docker container on port 8080 |
| `make clean`        | Remove build & coverage artifacts     |

---

## Testing & Coverage

```bash
make test
make coverage     # HTML report in dist/coverage.html
```

```bash
go mod tidy
```

