#!/bin/bash
# API test script
# Start server first: go run main.go

HOST="http://localhost:1323"

echo "========== /demo/search =========="
curl -s "${HOST}/demo/search?tag=go&tag=web&tag=api" | python3 -m json.tool

echo ""
echo "========== /demo/err/debug/hello =========="
curl -s "${HOST}/demo/err/debug/hello" | python3 -m json.tool

echo ""
echo "========== /demo/user/phone =========="
curl -s "${HOST}/demo/user/phone" | python3 -m json.tool
