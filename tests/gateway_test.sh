#!/bin/bash

# Gateway Integration Test Script
# Tests all endpoints of the API Gateway

GATEWAY_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "MockingCode API Gateway Integration Test"
echo "========================================"
echo ""

# Test variables
ACCESS_TOKEN=""
PROJECT_ID=""

# Helper functions
test_endpoint() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local auth=$5
    
    echo -n "Testing: $name... "
    
    if [ -n "$auth" ]; then
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $auth" \
                -d "$data")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
                -H "Authorization: Bearer $auth")
        fi
    else
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
                -H "Content-Type: application/json" \
                -d "$data")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$url")
        fi
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$status_code" -ge 200 ] && [ "$status_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASSED${NC} (HTTP $status_code)"
        echo "$body" | jq . 2>/dev/null || echo "$body"
        echo "$body"
    else
        echo -e "${RED}✗ FAILED${NC} (HTTP $status_code)"
        echo "$body"
        return 1
    fi
    echo ""
}

# 1. Test Health Check
echo "=== 1. Health Check ==="
test_endpoint "Gateway Health" "GET" "$GATEWAY_URL/health"

# 2. Test Registration
echo "=== 2. User Registration ==="
RANDOM_EMAIL="test$(date +%s)@example.com"
REGISTER_DATA="{\"email\":\"$RANDOM_EMAIL\",\"password\":\"password123\"}"
REGISTER_RESPONSE=$(test_endpoint "Register User" "POST" "$GATEWAY_URL/auth/register" "$REGISTER_DATA")

# 3. Test Login
echo "=== 3. User Login ==="
LOGIN_DATA="{\"email\":\"$RANDOM_EMAIL\",\"password\":\"password123\"}"
LOGIN_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA")

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token' 2>/dev/null)

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ Login successful, token received${NC}"
    echo "Token: ${ACCESS_TOKEN:0:20}..."
else
    echo -e "${RED}✗ Login failed or no token received${NC}"
    echo "$LOGIN_RESPONSE"
    exit 1
fi
echo ""

# 4. Test Project Creation
echo "=== 4. Project Management ==="
PROJECT_DATA='{"name":"Test Project","description":"Test project via gateway"}'
PROJECT_RESPONSE=$(curl -s -X POST "$GATEWAY_URL/api/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$PROJECT_DATA")

PROJECT_ID=$(echo "$PROJECT_RESPONSE" | jq -r '.project.id' 2>/dev/null)

if [ -n "$PROJECT_ID" ] && [ "$PROJECT_ID" != "null" ]; then
    echo -e "${GREEN}✓ Project created${NC}"
    echo "Project ID: $PROJECT_ID"
else
    echo -e "${YELLOW}⚠ Project creation response:${NC}"
    echo "$PROJECT_RESPONSE"
fi
echo ""

# 5. Test Get Projects
echo "=== 5. Get Projects List ==="
test_endpoint "Get Projects" "GET" "$GATEWAY_URL/api/projects" "" "$ACCESS_TOKEN"

# 6. Test Get Specific Project (if we have ID)
if [ -n "$PROJECT_ID" ] && [ "$PROJECT_ID" != "null" ]; then
    echo "=== 6. Get Specific Project ==="
    test_endpoint "Get Project $PROJECT_ID" "GET" "$GATEWAY_URL/api/projects/$PROJECT_ID" "" "$ACCESS_TOKEN"
fi

# 7. Test Unauthorized Access
echo "=== 7. Test Authentication ==="
echo -n "Testing unauthorized access... "
UNAUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$GATEWAY_URL/api/projects")
UNAUTH_CODE=$(echo "$UNAUTH_RESPONSE" | tail -n1)

if [ "$UNAUTH_CODE" -eq 401 ]; then
    echo -e "${GREEN}✓ PASSED${NC} (Correctly blocked with HTTP 401)"
else
    echo -e "${RED}✗ FAILED${NC} (Expected 401, got $UNAUTH_CODE)"
fi
echo ""

# Summary
echo "========================================"
echo "Test Suite Completed"
echo "========================================"
echo ""
echo -e "${YELLOW}Note:${NC} Some tests may fail if services are not fully initialized."
echo "Run the tests again after all services are up and running."

