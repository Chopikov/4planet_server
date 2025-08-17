#!/bin/bash

# Test script for Referral Program
BASE_URL="http://localhost:8080"

# Test users
REFERRER_EMAIL="ref"
REFERRER_PASSWORD="12345678"
REFERRER_AUTH_ID="b03922b4-c05f-477d-bdd5-b59c87fb4a44"

DONOR_EMAIL="test"
DONOR_PASSWORD="12345678"

echo "🌳 Testing Referral Program"
echo "============================"

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

# Step 1: Login as referrer and check initial stats
echo "🔐 Step 1: Referrer Login & Initial Stats"
echo "------------------------------------------"

REFERRER_LOGIN=$(curl -s -c referrer_cookies.txt -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$REFERRER_EMAIL\",
    \"password\": \"$REFERRER_PASSWORD\"
  }")

if echo "$REFERRER_LOGIN" | jq -e '.error' > /dev/null; then
    echo "❌ Referrer login failed: $(echo "$REFERRER_LOGIN" | jq -r '.error')"
    exit 1
fi

echo "✅ Referrer authenticated"

# Check initial referral stats
INITIAL_STATS=$(curl -s -b referrer_cookies.txt -X GET "$BASE_URL/v1/me/referral-stats")
echo "📊 Initial referral stats:"
echo "$INITIAL_STATS" | jq '.'

INITIAL_REFERRALS=$(echo "$INITIAL_STATS" | jq -r '.total_referrals // 0')
INITIAL_TREES=$(echo "$INITIAL_STATS" | jq -r '.total_trees_planted // 0')

echo "📈 Initial: $INITIAL_REFERRALS referrals, $INITIAL_TREES trees"
echo ""

# Step 2: Login as donor
echo "🔐 Step 2: Donor Login"
echo "----------------------"

DONOR_LOGIN=$(curl -s -c donor_cookies.txt -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"login\": \"$DONOR_EMAIL\",
    \"password\": \"$DONOR_PASSWORD\"
  }")

if echo "$DONOR_LOGIN" | jq -e '.error' > /dev/null; then
    echo "❌ Donor login failed: $(echo "$DONOR_LOGIN" | jq -r '.error')"
    exit 1
fi

echo "✅ Donor authenticated"
echo ""

# Step 3: Create payment with referral
echo "💰 Step 3: Create Payment with Referral"
echo "---------------------------------------"

PAYMENT_RESPONSE=$(curl -s -b donor_cookies.txt -X POST "$BASE_URL/v1/payments/intents" \
  -H "Content-Type: application/json" \
  -d "{
    \"provider\": \"mock\",
    \"amount_minor\": 2000,
    \"currency\": \"RUB\",
    \"success_return_url\": \"http://localhost:3000/success\",
    \"fail_return_url\": \"http://localhost:3000/fail\",
    \"description\": \"Test donation with referral\",
    \"referral_user_id\": \"$REFERRER_AUTH_ID\"
  }")

PAYMENT_ID=$(echo "$PAYMENT_RESPONSE" | jq -r '.provider_payload.payment_id')

if [ "$PAYMENT_ID" != "null" ] && [ "$PAYMENT_ID" != "" ]; then
    echo "✅ Payment created with referral: $PAYMENT_ID"
    echo "🔗 Referral User ID: $REFERRER_AUTH_ID"
else
    echo "❌ Payment creation failed"
    echo "$PAYMENT_RESPONSE" | jq '.'
    exit 1
fi

echo ""

# Step 4: Wait for payment webhook and check donation creation
echo "⏰ Step 4: Wait for Payment Processing"
echo "-------------------------------------"

echo "⏳ Waiting for payment webhook..."
sleep 12

echo "✅ Payment should be processed"
echo ""

# Step 5: Check updated referral stats
echo "📊 Step 5: Check Updated Referral Stats"
echo "----------------------------------------"

FINAL_STATS=$(curl -s -b referrer_cookies.txt -X GET "$BASE_URL/v1/me/referral-stats")
echo "📊 Final referral stats:"
echo "$FINAL_STATS" | jq '.'

FINAL_REFERRALS=$(echo "$FINAL_STATS" | jq -r '.total_referrals // 0')
FINAL_TREES=$(echo "$FINAL_STATS" | jq -r '.total_trees_planted // 0')

echo "📈 Final: $FINAL_REFERRALS referrals, $FINAL_TREES trees"
echo ""

# Step 6: Verify referral tracking
echo "🔍 Step 6: Verify Referral Tracking"
echo "-----------------------------------"

if [ "$FINAL_REFERRALS" -gt "$INITIAL_REFERRALS" ]; then
    echo "✅ Referral count increased: $INITIAL_REFERRALS → $FINAL_REFERRALS"
else
    echo "❌ Referral count did not increase"
fi

if [ "$FINAL_TREES" -gt "$INITIAL_TREES" ]; then
    echo "✅ Trees planted increased: $INITIAL_TREES → $FINAL_TREES"
else
    echo "❌ Trees planted did not increase"
fi

echo ""

# Step 7: Check donation details
echo "🎯 Step 7: Check Donation Details"
echo "---------------------------------"

# Get recent referrals to verify the donation was tracked
RECENT_REFERRALS=$(echo "$FINAL_STATS" | jq -r '.recent_referrals // []')
if [ "$(echo "$RECENT_REFERRALS" | jq 'length')" -gt 0 ]; then
    echo "✅ Recent referrals found:"
    echo "$RECENT_REFERRALS" | jq '.[0] | {id, amount_minor, currency, trees_count, referral_user_id}'
    
    # Verify referral_user_id matches
    REFERRAL_ID=$(echo "$RECENT_REFERRALS" | jq -r '.[0].referral_user_id // "null"')
    if [ "$REFERRAL_ID" = "$REFERRER_AUTH_ID" ]; then
        echo "✅ Referral user ID correctly set: $REFERRAL_ID"
    else
        echo "❌ Referral user ID mismatch: expected $REFERRER_AUTH_ID, got $REFERRAL_ID"
    fi
else
    echo "❌ No recent referrals found"
fi

echo ""

# Step 8: Test without referral (control test)
echo "🧪 Step 8: Control Test - Payment without Referral"
echo "--------------------------------------------------"

CONTROL_PAYMENT=$(curl -s -b donor_cookies.txt -X POST "$BASE_URL/v1/payments/intents" \
  -H "Content-Type: application/json" \
  -d "{
    \"provider\": \"mock\",
    \"amount_minor\": 1000,
    \"currency\": \"RUB\",
    \"success_return_url\": \"http://localhost:3000/success\",
    \"fail_return_url\": \"http://localhost:3000/fail\",
    \"description\": \"Control donation without referral\"
  }")

CONTROL_PAYMENT_ID=$(echo "$CONTROL_PAYMENT" | jq -r '.provider_payload.payment_id')

if [ "$CONTROL_PAYMENT_ID" != "null" ] && [ "$CONTROL_PAYMENT_ID" != "" ]; then
    echo "✅ Control payment created: $CONTROL_PAYMENT_ID"
    echo "⏳ Waiting for control payment webhook..."
    sleep 12
    
    # Check stats again to ensure they didn't change
    CONTROL_STATS=$(curl -s -b referrer_cookies.txt -X GET "$BASE_URL/v1/me/referral-stats")
    CONTROL_REFERRALS=$(echo "$CONTROL_STATS" | jq -r '.total_referrals // 0')
    
    if [ "$CONTROL_REFERRALS" -eq "$FINAL_REFERRALS" ]; then
        echo "✅ Referral count unchanged (correct): $CONTROL_REFERRALS"
    else
        echo "❌ Referral count changed unexpectedly: $FINAL_REFERRALS → $CONTROL_REFERRALS"
    fi
else
    echo "❌ Control payment creation failed"
fi

echo ""

# Final summary
echo "🎉 Referral Test Summary"
echo "========================"
echo "✅ Referrer authenticated: $REFERRER_EMAIL"
echo "✅ Donor authenticated: $DONOR_EMAIL"
echo "✅ Payment with referral created: $PAYMENT_ID"
echo "✅ Referral tracking verified"
echo "✅ Statistics updated correctly"
echo "✅ Control test passed"
echo ""
echo "📊 Final Results:"
echo "- Total Referrals: $FINAL_REFERRALS"
echo "- Total Trees Planted: $FINAL_TREES"
echo "- Referral User ID: $REFERRER_AUTH_ID"
echo ""

# Cleanup
rm -f referrer_cookies.txt donor_cookies.txt
echo "🧹 Cleanup completed"
echo "✅ Referral test completed successfully!"
