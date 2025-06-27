# Simple Setup Guide

This document walks you through the steps to get the **google-search-agent** up and running, either as a standalone binary or in Docker.

---

## 1. Prerequisites

- **Go 1.23+** (for building from source)
- **Make** & **Git**
- **Docker** (optional, for containerized run)
- A Google Custom Search API key and CSE ID:
  - Create or retrieve these from the [Google Cloud Console](https://console.cloud.google.com/).

---

## 2. Clone the Repository

```bash
git clone https://github.com/ruslannv/google-search-agent.git
cd google-search-agent
````

---

## 3. Configure Credentials

Create a file named `.env` in the project root (this will be used by both Makefile and Docker):

```dotenv
GOOGLE_API_KEY=your_actual_api_key_here
GOOGLE_CSE_ID=your_custom_search_engine_id_here
```

Alternatively, you can export them directly in your shell:

```bash
export GOOGLE_API_KEY=your_actual_api_key_here
export GOOGLE_CSE_ID=your_custom_search_engine_id_here
```

---

## 4. Build & Run as a Binary

1. **Build** the agent:

   ```bash
   make build
   ```

   This produces `dist/google-search-agent`.

2. **Run** the agent on port 8080:

   ```bash
   make run
   ```

   If you prefer to run the binary directly:

   ```bash
   GOOGLE_API_KEY=$GOOGLE_API_KEY \
   GOOGLE_CSE_ID=$GOOGLE_CSE_ID \
   ./dist/google-search-agent \
     -listen=0.0.0.0 -port=8080
   ```

3. **Verify** it’s healthy:

   ```bash
   curl http://localhost:8080/health
   # → {"status":"ok"}
   ```

---

## 5. Invoke a Search

Send a JSON-RPC request to `/`:

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d @payload.json
```

Where `payload.json` contains:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "google_search",
    "arguments": { "query": "IBM WatsonX" }
  },
  "id": 1
}
```

---

## 6. Run in Docker

1. **Build** the image:

   ```bash
   make docker-build
   ```

2. **Run** the container (loads `.env` automatically):

   ```bash
   make docker-run
   ```

3. **Test** health:

   ```bash
   curl http://localhost:8080/health
   # → {"status":"ok"}
   ```

---

## 7. Run the Tests

```bash
make test
```

All unit tests will execute (including the mocked HTTP tests for `googleSearch`).

---

## 8. Clean Up

To remove build artifacts:

```bash
make clean
```

