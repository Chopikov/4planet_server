# Mock Payment Provider

The Mock Payment Provider is a testing and development tool that simulates the full payment flow without requiring real payment processing. It's perfect for development, testing, and demonstration purposes.

## Features

- **Full Payment Flow Simulation**: Creates payment intents and processes webhooks
- **Automatic Webhook Delivery**: Sends webhooks after a 10-second delay to simulate real-world payment processing
- **No External Dependencies**: Works entirely within your application
- **Easy Testing**: Perfect for testing payment flows without real money
- **Production Ready**: Can be easily disabled for production use

## How It Works

1. **Payment Intent Creation**: When you create a payment intent using the mock provider, it:
   - Creates a payment record in the database with `pending` status
   - Returns a mock redirect URL
   - Schedules a webhook to be sent after 10 seconds

2. **Webhook Processing**: After 10 seconds, the mock provider:
   - Automatically sends a webhook with `completed` status
   - Updates the payment record in the database
   - Logs the status change

3. **Subscription Support**: Works the same way for subscriptions, with webhooks sent after 10 seconds

## Configuration

The mock provider is configured in `cmd/api/main.go`:

```go
"mock": {
    ProviderName: "mock",
    PublicID:     "mock-public-id",
    Secret:       "mock-secret",
    BaseURL:      cfg.App.BaseURL,
    Enabled:      true,
    PaymentTexts: payments.PaymentTexts{
        DefaultDonationDescription:  cfg.Payments.DefaultDonationDescription,
        BaseSubscriptionDescription: cfg.Payments.BaseSubscriptionDescription,
        SubscriptionDescriptions: payments.SubscriptionDescriptions{
            Monthly: cfg.Payments.SubscriptionDescriptions.Monthly,
            Yearly:  cfg.Payments.SubscriptionDescriptions.Yearly,
            Custom:  cfg.Payments.SubscriptionDescriptions.Custom,
        },
    },
},
```

## Usage

### Creating a Payment Intent

```bash
curl -X POST "http://localhost:8080/v1/payments/intent" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -d '{
    "provider": "mock",
    "amount_minor": 1000,
    "currency": "RUB",
    "success_return_url": "http://localhost:3000/success",
    "fail_return_url": "http://localhost:3000/fail",
    "description": "Test donation via mock provider"
  }'
```

### Creating a Subscription Intent

```bash
curl -X POST "http://localhost:8080/v1/subscriptions/intent" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -d '{
    "provider": "mock",
    "amount_minor": 500,
    "currency": "RUB",
    "success_return_url": "http://localhost:3000/success",
    "fail_return_url": "http://localhost:3000/fail",
    "interval_months": 1,
    "description": "Test monthly subscription via mock provider"
  }'
```

## Testing

Use the provided test script to verify the mock provider works correctly:

```bash
./scripts/test-mock-payment.sh
```

This script will:
1. Create a test payment intent
2. Wait for the webhook (10 seconds)
3. Create a test subscription intent
4. Wait for the subscription webhook (10 seconds)
5. Verify that both are processed correctly

## Webhook Flow

### Payment Webhook
```json
{
  "payment_id": "uuid-here",
  "amount": 1000,
  "currency": "RUB",
  "status": "completed",
  "timestamp": 1234567890
}
```

### Subscription Webhook
```json
{
  "subscription_id": "uuid-here",
  "amount": 500,
  "currency": "RUB",
  "status": "active",
  "interval_months": 1,
  "timestamp": 1234567890
}
```

## Database Changes

The mock provider automatically updates the database:

- **Payments**: Status changes from `pending` to `succeeded` after webhook
- **Subscriptions**: Status changes from `pending` to `active` after webhook
- **Timestamps**: `occurred_at` and `started_at` fields are updated accordingly

## Production Use

To disable the mock provider in production:

1. Set `Enabled: false` in the configuration
2. Or remove the mock provider configuration entirely

```go
"mock": {
    ProviderName: "mock",
    PublicID:     "mock-public-id",
    Secret:       "mock-secret",
    BaseURL:      cfg.App.BaseURL,
    Enabled:      false, // Disabled for production
    // ... rest of config
},
```

## Advantages

- **No External Dependencies**: Works offline and doesn't require internet
- **Predictable Timing**: Webhooks are always sent after exactly 10 seconds
- **Full Flow Testing**: Tests the complete payment processing pipeline
- **Easy Debugging**: All operations are logged and visible in the database
- **Cost Effective**: No real money is involved in testing

## Limitations

- **Fixed Timing**: Webhooks are always sent after 10 seconds (not configurable)
- **Simple Status Flow**: Only supports basic status transitions
- **No Real Payment Processing**: Cannot test actual payment provider integrations
- **Internal Webhooks**: Webhooks are processed internally, not sent to external endpoints

## Troubleshooting

### Payment Not Updating
- Check that the mock provider is enabled in configuration
- Verify the webhook is being processed (check logs)
- Ensure the database connection is working

### Webhook Not Received
- The mock provider sends webhooks after 10 seconds
- Check application logs for webhook processing
- Verify the payment/subscription IDs exist in the database

### Configuration Issues
- Ensure the mock provider is properly configured in `main.go`
- Check that all required fields are present
- Verify the provider name matches exactly: `"mock"`

## Future Enhancements

Potential improvements for the mock provider:

- **Configurable Delays**: Allow custom webhook timing
- **Multiple Status Flows**: Support for failed payments, cancellations, etc.
- **External Webhook Endpoints**: Send webhooks to actual HTTP endpoints
- **Webhook History**: Track and replay webhook events
- **Status Simulation**: Simulate various payment provider scenarios

