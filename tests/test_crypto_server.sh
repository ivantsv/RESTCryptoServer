BASE_URL="http://localhost:8080"
TOKEN=""

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_test() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

extract_token() {
    echo "$1" | grep -o '"token":"[^"]*"' | cut -d'"' -f4
}

echo -e "${BLUE}=== Starting Crypto Server Tests ===${NC}\n"

# Test 1: Register new user
echo "Test 1: User Registration"
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}')

HTTP_CODE=$(echo "$REGISTER_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "201" ]; then
    TOKEN=$(extract_token "$RESPONSE_BODY")
    print_test 0 "User registration successful"
    echo "Response: $RESPONSE_BODY"
else
    print_test 1 "User registration failed (HTTP: $HTTP_CODE)"
    echo "Response: $RESPONSE_BODY"
fi
echo ""

# Test 2: Login with created user
echo "Test 2: User Login"
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}')

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    TOKEN=$(extract_token "$RESPONSE_BODY")
    print_test 0 "User login successful"
    echo "Token: $TOKEN"
else
    print_test 1 "User login failed (HTTP: $HTTP_CODE)"
    echo "Response: $RESPONSE_BODY"
fi
echo ""

# Test 3: Try to access protected endpoint without token
echo "Test 3: Access without authentication"
NO_AUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/crypto")
HTTP_CODE=$(echo "$NO_AUTH_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "401" ]; then
    print_test 0 "Protected endpoint correctly requires authentication"
else
    print_test 1 "Protected endpoint should require authentication (HTTP: $HTTP_CODE)"
fi
echo ""

# Test 4: Get empty crypto list
echo "Test 4: Get initial crypto list (should be empty)"
if [ -n "$TOKEN" ]; then
    CRYPTO_LIST_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$CRYPTO_LIST_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$CRYPTO_LIST_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully retrieved crypto list"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to get crypto list (HTTP: $HTTP_CODE)"
    fi
else
    print_test 1 "No token available for authentication"
fi
echo ""

# Test 5: Add Bitcoin
echo "Test 5: Add Bitcoin (BTC)"
if [ -n "$TOKEN" ]; then
    ADD_BTC_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
      "$BASE_URL/crypto" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"symbol":"btc"}')
    
    HTTP_CODE=$(echo "$ADD_BTC_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$ADD_BTC_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully added Bitcoin"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to add Bitcoin (HTTP: $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
    fi
else
    print_test 1 "No token available for authentication"
fi
echo ""

# Test 6: Add Ethereum
echo "Test 6: Add Ethereum (ETH)"
if [ -n "$TOKEN" ]; then
    ADD_ETH_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
      "$BASE_URL/crypto" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"symbol":"eth"}')
    
    HTTP_CODE=$(echo "$ADD_ETH_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$ADD_ETH_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully added Ethereum"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to add Ethereum (HTTP: $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
    fi
fi
echo ""

# Test 7: Try to add duplicate
echo "Test 7: Try to add duplicate Bitcoin"
if [ -n "$TOKEN" ]; then
    DUP_BTC_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
      "$BASE_URL/crypto" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"symbol":"btc"}')
    
    HTTP_CODE=$(echo "$DUP_BTC_RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" = "409" ]; then
        print_test 0 "Correctly rejected duplicate cryptocurrency"
    else
        print_test 1 "Should reject duplicate cryptocurrency (HTTP: $HTTP_CODE)"
    fi
fi
echo ""

# Test 8: Get updated crypto list
echo "Test 8: Get updated crypto list"
if [ -n "$TOKEN" ]; then
    UPDATED_LIST_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$UPDATED_LIST_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$UPDATED_LIST_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully retrieved updated crypto list"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to get updated crypto list"
    fi
fi
echo ""

# Test 9: Get specific cryptocurrency
echo "Test 9: Get Bitcoin details"
if [ -n "$TOKEN" ]; then
    BTC_DETAILS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto/btc" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$BTC_DETAILS_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$BTC_DETAILS_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully retrieved Bitcoin details"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to get Bitcoin details"
    fi
fi
echo ""

# Test 10: Refresh cryptocurrency price
echo "Test 10: Refresh Bitcoin price"
if [ -n "$TOKEN" ]; then
    sleep 2 # Wait a bit to potentially get different price
    
    REFRESH_RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT \
      "$BASE_URL/crypto/btc/refresh" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")
    
    HTTP_CODE=$(echo "$REFRESH_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$REFRESH_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully refreshed Bitcoin price"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to refresh Bitcoin price"
    fi
fi
echo ""

# Test 11: Get price history
echo "Test 11: Get Bitcoin price history"
if [ -n "$TOKEN" ]; then
    HISTORY_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto/btc/history" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$HISTORY_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$HISTORY_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully retrieved Bitcoin price history"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to get Bitcoin price history"
    fi
fi
echo ""

# Test 12: Get price statistics
echo "Test 12: Get Bitcoin price statistics"
if [ -n "$TOKEN" ]; then
    STATS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto/btc/stats" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$STATS_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$STATS_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully retrieved Bitcoin statistics"
        echo "Response: $RESPONSE_BODY"
    else
        print_test 1 "Failed to get Bitcoin statistics"
    fi
fi
echo ""

# Test 13: Try to get non-existent cryptocurrency
echo "Test 13: Try to get non-existent cryptocurrency"
if [ -n "$TOKEN" ]; then
    NONEXISTENT_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto/nonexistent" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$NONEXISTENT_RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" = "400" ]; then
        print_test 0 "Correctly handled non-existent cryptocurrency"
    else
        print_test 1 "Should return error for non-existent crypto (HTTP: $HTTP_CODE)"
    fi
fi
echo ""

# Test 14: Delete cryptocurrency
echo "Test 14: Delete Ethereum"
if [ -n "$TOKEN" ]; then
    DELETE_RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE \
      "$BASE_URL/crypto/eth" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$DELETE_RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully deleted Ethereum"
    else
        print_test 1 "Failed to delete Ethereum (HTTP: $HTTP_CODE)"
    fi
fi
echo ""

# Test 15: Verify deletion
echo "Test 15: Verify Ethereum deletion"
if [ -n "$TOKEN" ]; then
    VERIFY_DELETE_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
      "$BASE_URL/crypto" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$VERIFY_DELETE_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$VERIFY_DELETE_RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" = "200" ]; then
        print_test 0 "Successfully verified deletion"
        echo "Updated crypto list: $RESPONSE_BODY"
    else
        print_test 1 "Failed to verify deletion"
    fi
fi
echo ""

echo -e "${BLUE}=== Database Inspection Commands ===${NC}\n"

echo "To inspect PostgreSQL database:"
echo "docker exec -it \$(docker-compose ps -q db) psql -U user -d myapp"
echo ""
echo "Once connected to PostgreSQL, run:"
echo "\\dt                    -- List all tables"
echo "SELECT * FROM users;   -- Show all users"
echo "SELECT * FROM crypto;  -- Show all cryptocurrencies"
echo "\\q                     -- Quit"
echo ""

echo "To inspect Redis cache:"
echo "docker exec -it \$(docker-compose ps -q redis) redis-cli"
echo ""
echo "Once connected to Redis, run:"
echo "KEYS *                           -- Show all keys"
echo "LLEN price_history:btc           -- Get Bitcoin history length"
echo "LRANGE price_history:btc 0 5     -- Get first 6 Bitcoin price entries"
echo "LRANGE price_history:eth 0 -1    -- Get all Ethereum price entries"
echo "TYPE price_history:btc           -- Check data type"
echo "EXIT                             -- Quit"

echo -e "\n${BLUE}=== Advanced Tests ===${NC}"