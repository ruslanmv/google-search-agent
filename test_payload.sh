curl http://localhost:8080/health
curl -X POST http://localhost:8080/ \
     -H "Content-Type: application/json" \
     -d @payload.json
