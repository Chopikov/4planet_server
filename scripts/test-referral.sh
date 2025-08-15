#!/bin/bash

# Test script for the referral system
# This demonstrates the complete referral flow

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/v1"

echo "🌱 Testing 4Planet Referral System"
echo "=================================="

# Test 1: Health check
echo -e "\n1. Testing API health..."
curl -s "$BASE_URL/health" | jq '.'

# Test 2: Resolve a non-existent share (should return 404)
echo -e "\n2. Testing share resolution (non-existent)..."
curl -s "$API_URL/shares/resolve/non-existent-slug" | jq '.'

# Test 3: Test payment intent endpoint structure
echo -e "\n3. Testing payment intent endpoint structure..."
curl -s -X POST "$API_URL/payments/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cloudpayments",
    "amount_minor": 1900,
    "currency": "RUB",
    "success_return_url": "https://example.com/success",
    "fail_return_url": "https://example.com/fail",
    "referral_user_id": "test-referral-user"
  }' | jq '.'

# Test 4: Test share endpoints (should require auth)
echo -e "\n4. Testing share endpoints (should require auth)..."
echo "Profile share:"
curl -s -X POST "$API_URL/shares/profile" | jq '.'

echo "Share stats:"
curl -s "$API_URL/shares/stats" | jq '.'

# Test 5: Check OpenAPI spec for referral fields
echo -e "\n5. Checking OpenAPI spec for referral fields..."
curl -s "$BASE_URL/openapi.yaml" | grep -A 3 -B 3 "referral_user_id" || echo "referral_user_id not found in OpenAPI spec"

echo -e "\n✅ Referral system test completed!"
echo -e "\nTo test with authentication, you would need to:"
echo "1. Register/login to get a session"
echo "2. Create share tokens"
echo "3. Make payments with referral_user_id"
echo "4. View referral statistics"
