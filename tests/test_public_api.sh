#!/bin/bash

# Public API Test Script

GATEWAY_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================"
echo "MockingCode Public API Test"
echo "========================================"
echo ""

# 1. Register and Login
echo "=== 1. Setup: Register and Login ==="
EMAIL="publictest@test.com"
PASSWORD="test123"

curl -s -X POST "$GATEWAY_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" > /dev/null

TOKEN=$(curl -s -X POST "$GATEWAY_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" | jq -r '.access_token')

echo -e "${GREEN}✓ Logged in${NC}"
echo ""

# 2. Create Project
echo "=== 2. Create Project ==="
PROJECT=$(curl -s -X POST "$GATEWAY_URL/projects" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Public API Test","description":"Testing public API"}')

API_KEY=$(echo "$PROJECT" | jq -r '.api_key')

if [ -n "$API_KEY" ] && [ "$API_KEY" != "null" ]; then
    echo -e "${GREEN}✓ Project created${NC}"
    echo "API Key: $API_KEY"
else
    echo -e "${RED}✗ Failed to create project${NC}"
    echo "$PROJECT"
    exit 1
fi
echo ""

# 3. Test Public API GET
echo "=== 3. Test Public API GET ==="
RESPONSE=$(curl -s "$GATEWAY_URL/$API_KEY/products")
echo "$RESPONSE" | jq .
echo ""

# 4. Test Public API POST
echo "=== 4. Test Public API POST ==="
RESPONSE=$(curl -s -X POST "$GATEWAY_URL/$API_KEY/products" \
  -H "Content-Type: application/json" \
  -d '{"name":"Product 1","price":99.99}')
echo "$RESPONSE" | jq .
echo ""

echo "========================================"
echo "Public API Test Completed"
echo "========================================"

