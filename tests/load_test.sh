TOKEN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"loadtest","password":"test123"}')

TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

curl -s -X POST http://localhost:8080/crypto \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"btc"}' > /dev/null

curl -s -X POST http://localhost:8080/crypto \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"eth"}' > /dev/null

# Perform load test - 50 concurrent requests
echo "Starting load test with 50 requests..."

for i in {1..50}; do
  (
    case $((i % 4)) in
      0) curl -s http://localhost:8080/crypto -H "Authorization: Bearer $TOKEN" ;;
      1) curl -s http://localhost:8080/crypto/btc -H "Authorization: Bearer $TOKEN" ;;
      2) curl -s http://localhost:8080/crypto/btc/history -H "Authorization: Bearer $TOKEN" ;;
      3) curl -s http://localhost:8080/crypto/btc/stats -H "Authorization: Bearer $TOKEN" ;;
    esac
  ) &
done

wait
echo "Load test completed!"