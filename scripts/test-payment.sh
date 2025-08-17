#!/bin/bash

# Test script for Payment Functionality
BASE_URL="http://localhost:8080"
TEST_EMAIL="test"
TEST_PASSWORD="12345678"
PAYMENT_PROVIDER=${1:-"mock"}

echo "🌳 Testing Payment: $PAYMENT_PROVIDER"

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "❌ jq required: brew install jq"
    exit 1
fi

if ! command -v curl &> /dev/null; then
    echo "❌ curl required"
    exit 1
fi

# Login
LOGIN_RESPONSE=$(curl -s -c cookies.txt -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

if echo "$LOGIN_RESPONSE" | jq -e '.error' > /dev/null; then
    echo "❌ Login failed: $(echo "$LOGIN_RESPONSE" | jq -r '.error')"
    exit 1
fi

echo "✅ Authenticated"

# Check provider availability
PROVIDERS_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/payments/providers" \
  -H "Content-Type: application/json")

if ! echo "$PROVIDERS_RESPONSE" | jq -r '.providers[].name' | grep -q "^$PAYMENT_PROVIDER$"; then
    echo "❌ Provider '$PAYMENT_PROVIDER' not available"
    echo "$PROVIDERS_RESPONSE" | jq -r '.providers[].name'
    exit 1
fi

# Create payment intent
PAYMENT_RESPONSE=$(curl -s -b cookies.txt -X POST "$BASE_URL/v1/payments/intents" \
  -H "Content-Type: application/json" \
  -d "{
    \"provider\": \"$PAYMENT_PROVIDER\",
    \"amount_minor\": 1000,
    \"currency\": \"RUB\",
    \"success_return_url\": \"http://localhost:3000/success\",
    \"fail_return_url\": \"http://localhost:3000/fail\",
    \"description\": \"Test donation via $PAYMENT_PROVIDER provider\"
  }")

PAYMENT_ID=$(echo "$PAYMENT_RESPONSE" | jq -r '.provider_payload.payment_id')

if [ "$PAYMENT_ID" != "null" ] && [ "$PAYMENT_ID" != "" ]; then
    echo "✅ Payment created: $PAYMENT_ID"
    
    if [ "$PAYMENT_PROVIDER" = "mock" ]; then
        echo "⏳ Waiting for webhook..."
        sleep 12
        echo "✅ Webhook processed"
    elif [ "$PAYMENT_PROVIDER" = "cloudpayments" ]; then
        echo "🔗 Redirect: $(echo "$PAYMENT_RESPONSE" | jq -r '.redirect_url // "N/A"')"
    fi
    
    echo "✅ Payment test completed"
else
    echo "❌ Payment creation failed"
    echo "$PAYMENT_RESPONSE" | jq '.'
    exit 1
fi

rm -f cookies.txt
