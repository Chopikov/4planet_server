#!/bin/bash

# Test script for payment provider endpoints
BASE_URL="http://localhost:8080"

echo "🌳 Testing Payment Provider Structure"
echo "===================================="

# Test 1: Health check
echo -e "\n1️⃣ Testing health endpoint..."
curl -s "$BASE_URL/health" | jq '.'

# Test 2: Get supported payment providers (public endpoint)
echo -e "\n2️⃣ Testing payment providers endpoint..."
curl -s "$BASE_URL/v1/payments/providers" | jq '.'

# Test 3: Test payment intent with unsupported provider (should fail)
echo -e "\n3️⃣ Testing payment intent with unsupported provider..."
curl -s -X POST "$BASE_URL/v1/payments/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "stripe",
    "amount_minor": 1000,
    "currency": "RUB",
    "success_return_url": "http://example.com/success",
    "fail_return_url": "http://example.com/fail"
  }' | jq '.'

# Test 4: Test subscription intent with unsupported provider (should fail)
echo -e "\n4️⃣ Testing subscription intent with unsupported provider..."
curl -s -X POST "$BASE_URL/v1/subscriptions/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "paypal",
    "amount_minor": 1000,
    "currency": "RUB",
    "success_return_url": "http://example.com/success",
    "fail_return_url": "http://example.com/fail",
    "interval": "monthly",
    "interval_count": 1
  }' | jq '.'

# Test 5: Test payment intent with supported provider but no auth (should fail)
echo -e "\n5️⃣ Testing payment intent with supported provider but no auth..."
curl -s -X POST "$BASE_URL/v1/payments/intents" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cloudpayments",
    "amount_minor": 1000,
    "currency": "RUB",
    "success_return_url": "http://example.com/success",
    "fail_return_url": "http://example.com/fail"
  }' | jq '.'

echo -e "\n✅ Payment provider structure tests completed!"
echo -e "\n📝 Test Results Summary:"
echo "   - Test 1: Health check should succeed"
echo "   - Test 2: Should return list of supported providers (cloudpayments)"
echo "   - Test 3: Should return 'Unsupported payment provider' for stripe"
echo "   - Test 4: Should return 'Unsupported payment provider' for paypal"
echo "   - Test 5: Should return 'Authentication required' for cloudpayments"
echo -e "\n🏗️ Architecture Benefits:"
echo "   ✅ Provider-agnostic handlers"
echo "   ✅ Easy to add new providers"
echo "   ✅ Centralized provider configuration"
echo "   ✅ Clean separation of concerns"
echo "   ✅ Future-ready for Stripe, PayPal, etc."
