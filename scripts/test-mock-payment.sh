#!/bin/bash

# Comprehensive Test Script for Mock Payment Provider
echo "🌳 Testing Mock Payment Provider - Complete Suite"

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "❌ jq required: brew install jq"
    exit 1
fi

if ! command -v curl &> /dev/null; then
    echo "❌ curl required"
    exit 1
fi

echo "✅ Dependencies OK"
echo ""

# Test 1: Payment functionality
echo "💰 Test 1: Payment"
if ./scripts/test-payment.sh mock; then
    echo "✅ Payment test passed"
else
    echo "❌ Payment test failed"
    exit 1
fi

echo ""
echo "⏳ Waiting 2s..."
sleep 2

# Test 2: Subscription functionality
echo ""
echo "💳 Test 2: Subscription"
if ./scripts/test-subscription.sh mock; then
    echo "✅ Subscription test passed"
else
    echo "❌ Subscription test failed"
    exit 1
fi

echo ""
echo "🎉 All tests passed!"
echo "✅ Payment + Subscription + Charges working"
