#!/bin/bash

set -e

# The host and port where your server is running
HOST="http://localhost:8080"
echo "INFO: Testing server at $HOST"

# 1) INITIALIZE
# Send the 'initialize' request and use `grep` to extract the Mcp-Session-Id header.
echo "INFO: Sending 'initialize' request..."
RESPONSE_HEADERS=$(curl -s -i -X POST "$HOST/" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"initialize",
    "params":{"clientInfo":{"name":"curl","version":"1.0"}},
    "id":1
  }'
)

SESSION_ID=$(echo "$RESPONSE_HEADERS" | grep -i 'Mcp-Session-Id:' | awk '{print $2}' | tr -d '\r\n')

# Check if we successfully captured the Session ID
if [ -z "$SESSION_ID" ]; then
    echo "ERROR: Could not capture Mcp-Session-Id header."
    echo "Is the agent running in another terminal? (e.g., via 'make run')"
    exit 1
fi

echo "SUCCESS: Captured Session ID: $SESSION_ID"
echo "----------------------------------------------------"


# 2) CALL THE TOOL
# Send the 'tools/call' request, including the captured session ID in a new header.
echo "INFO: Sending 'tools/call' request for 'google_search'..."
curl -s -X POST "$HOST/" \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: $SESSION_ID" \
  -d '{
    "jsonrpc":"2.0",
    "method":"tools/call",
    "params":{
      "name":"google_search",
      "arguments":{"query":"IBM WatsonX"}
    },
    "id":2
  }' | json_pp

echo ""
echo "SUCCESS: Test completed."