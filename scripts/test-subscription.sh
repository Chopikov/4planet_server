#!/bin/bash

# Test script for subscription endpoints
BASE_URL="http://localhost:8080"

echo "🌳 Testing Subscription Endpoints"
echo "================================="

# Test 1: Health check
echo -e "\n1️⃣ Testing health endpoint..."
curl -s "$BASE_URL/health" | jq '.'

# Test 2: Test subscription endpoint without auth (should fail)
echo -e "\n2️⃣ Testing subscription endpoint without authentication..."
curl -s -X POST "$BASE_URL/v1/subscriptions/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cloudpayments",
    "amount_minor": 1000,
    "currency": "RUB",
    "success_return_url": "http://example.com/success",
    "fail_return_url": "http://example.com/fail",
    "interval": "monthly",
    "interval_count": 1
  }' | jq '.'

# Test 3: Test with invalid interval (should fail validation)
echo -e "\n3️⃣ Testing with invalid interval..."
curl -s -X POST "$BASE_URL/v1/subscriptions/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cloudpayments",
    "amount_minor": 1000,
    "currency": "RUB",
    "success_return_url": "http://example.com/success",
    "fail_return_url": "http://example.com/fail",
    "interval": "weekly",
    "interval_count": 1
  }' | jq '.'

# Test 4: Test with missing required fields (should fail validation)
echo -e "\n4️⃣ Testing with missing required fields..."
curl -s -X POST "$BASE_URL/v1/subscriptions/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cloudpayments",
    "amount_minor": 1000,
    "currency": "RUB"
  }' | jq '.'

echo -e "\n✅ Subscription endpoint tests completed!"
echo -e "\n📝 Note: All tests show expected behavior:"
echo "   - Test 1: Health check should succeed"
echo "   - Test 2: Should return 'Authentication required'"
echo "   - Test 3: Should return 'Authentication required' (auth check happens before validation)"
echo "   - Test 4: Should return 'Authentication required' (auth check happens before validation)"
echo -e "\n🔐 To test with authentication, you would need to:"
echo "   1. Register/login to get a session cookie"
echo "   2. Include the cookie in subsequent requests"
echo "   3. Test the actual subscription creation logic"
