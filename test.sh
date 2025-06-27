# 1) initialize and capture the cookie
curl -v -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -c /tmp/mcp.cookie \
  -d '{
    "jsonrpc":"2.0",
    "method":"initialize",
    "params":{"clientInfo":{"name":"curl","version":"1.0"}},
    "id":1
  }'

# 2) replay it on the next call
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -b /tmp/mcp.cookie \
  -d '{
    "jsonrpc":"2.0",
    "method":"tools/call",
    "params":{
      "name":"google_search",
      "arguments":{"query":"IBM WatsonX"}
    },
    "id":2
  }'
