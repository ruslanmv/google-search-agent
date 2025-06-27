
#  Running the Agent

Once your **google-search-agent** is running (for example on `localhost:8080`), you can interact with it over plain HTTP using the MCP “streamable” HTTP endpoint.

---

## 1. Health & Version

Before making any search calls, confirm your agent is up and healthy:

```bash
# Health check
curl http://localhost:8080/health
# → {"status":"ok"}

# Version info
curl http://localhost:8080/version
# → {"name":"google-search-agent","version":"0.1.0"}
````

---

## 2. JSON-RPC Search Call

All search requests are sent as JSON-RPC 2.0 `tools/call` calls to the root path (`/`). Here’s a `curl` example:

```bash
curl -X POST http://localhost:8080/ \
     -H "Content-Type: application/json" \
     -d '{
       "jsonrpc": "2.0",
       "method": "tools/call",
       "params": {
         "name": "google_search",
         "arguments": {
           "query": "IBM WatsonX"
         }
       },
       "id": 1
     }'
```

A successful response looks like:

```json
{
  "jsonrpc": "2.0",
  "result": [
    {
      "title": "IBM WatsonX",
      "link": "https://www.ibm.com/cloud/watsonx"
    },
    {
      "title": "Watsonx.ai Overview",
      "link": "https://www.ibm.com/cloud/watsonx-ai"
    }
  ],
  "id": 1
}
```

---

## 3. Using a Saved Payload

Instead of embedding JSON in the command line, store your request in a file (e.g. `payload.json`):

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "google_search",
    "arguments": {
      "query": "WatsonX pricing"
    }
  },
  "id": 2
}
```

Then invoke:

```bash
curl -X POST http://localhost:8080/ \
     -H "Content-Type: application/json" \
     -d @payload.json
```

---

## 4. Programmatic Calls

If you’re writing your own client, you can issue the same JSON-RPC over HTTP. For example, in Go:

```go
package main

import (
  "bytes"
  "encoding/json"
  "net/http"
)

func main() {
  reqBody := map[string]interface{}{
    "jsonrpc": "2.0",
    "method":  "tools/call",
    "params": map[string]interface{}{
      "name": "google_search",
      "arguments": map[string]string{
        "query": "WatsonX tutorials",
      },
    },
    "id": 3,
  }
  buf, _ := json.Marshal(reqBody)

  resp, err := http.Post("http://localhost:8080/", "application/json", bytes.NewReader(buf))
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  // decode resp.Body ...
}
```

---

## Summary

* **GET** `/health` and **GET** `/version` are simple status checks.
* **POST** `/` (JSON-RPC) with `"method":"tools/call"` and `"name":"google_search"` is how you perform searches.
* You can use `curl`, Postman, or any HTTP-capable client to drive the agent.
