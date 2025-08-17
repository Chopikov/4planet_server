package payments

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockPaymentService handles mock payment integration for testing
type MockPaymentService struct {
	db      *gorm.DB
	baseURL string
	texts   PaymentTexts
}

// NewMockPaymentService creates a new mock payment service
func NewMockPaymentService(baseURL string, paymentTexts PaymentTexts) *MockPaymentService {
	return &MockPaymentService{
		db:      database.GetDB(),
		baseURL: baseURL,
		texts:   paymentTexts,
	}
}

// GetProviderName returns the name of the payment provider
func (s *MockPaymentService) GetProviderName() string {
	return "mock"
}

// CreatePaymentIntent creates a payment intent for one-time payment
func (s *MockPaymentService) CreatePaymentIntent(req *PaymentIntentRequest, authUserID string) (*PaymentIntentResponse, error) {
	// Generate a realistic-looking provider payment ID
	providerPaymentID := fmt.Sprintf("mock_pay_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Create payment record
	payment := &models.Payment{
		ID:                uuid.New(),
		Provider:          models.PaymentProviderMock,
		ProviderPaymentID: &providerPaymentID,
		AuthUserID:        &authUserID,
		AmountMinor:       req.AmountMinor,
		Currency:          models.Currency(req.Currency),
		Status:            models.PaymentStatusPending,
		Meta: map[string]interface{}{
			"description": req.Description,
			"project_id":  req.ProjectID,
		},
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// Create mock payment intent response
	response := &PaymentIntentResponse{
		Provider:    "mock",
		RedirectURL: fmt.Sprintf("%s/mock/payment/%s", s.baseURL, payment.ID),
		ProviderPayload: map[string]interface{}{
			"payment_id": payment.ID.String(),
			"amount":     req.AmountMinor,
			"currency":   req.Currency,
			"status":     "pending",
		},
	}

	// Schedule webhook after 10 seconds
	go s.schedulePaymentWebhook(payment, req)

	return response, nil
}

// CreateSubscriptionIntent creates a subscription intent for recurring payments
func (s *MockPaymentService) CreateSubscriptionIntent(req *SubscriptionIntentRequest, authUserID string) (*SubscriptionIntentResponse, error) {
	// Generate a realistic-looking provider subscription ID
	providerSubscriptionID := fmt.Sprintf("mock_sub_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Create subscription record
	subscription := &models.Subscription{
		ID:                     uuid.New(),
		Provider:               models.PaymentProviderMock,
		ProviderSubscriptionID: &providerSubscriptionID,
		AuthUserID:             authUserID,
		AmountMinor:            req.AmountMinor,
		Currency:               models.Currency(req.Currency),
		Status:                 models.SubscriptionStatusPending,
		IntervalMonths:         req.IntervalMonths,
		Meta: map[string]interface{}{
			"description": req.Description,
			"project_id":  req.ProjectID,
		},
		StartedAt: time.Now(),
	}

	if err := s.db.Create(subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscription record: %w", err)
	}

	// Create mock subscription intent response
	response := &SubscriptionIntentResponse{
		Provider:    "mock",
		RedirectURL: fmt.Sprintf("%s/mock/subscription/%s", s.baseURL, subscription.ID),
		ProviderPayload: map[string]interface{}{
			"subscription_id": subscription.ID.String(),
			"amount":          req.AmountMinor,
			"currency":        req.Currency,
			"status":          "pending",
			"interval_months": req.IntervalMonths,
		},
	}

	// Schedule webhook after 10 seconds
	go s.scheduleSubscriptionWebhook(subscription, req)

	return response, nil
}

// ProcessWebhook processes webhook events from the mock payment provider
func (s *MockPaymentService) ProcessWebhook(payload []byte, signature string) error {
	var webhookData map[string]interface{}
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		return fmt.Errorf("failed to unmarshal webhook payload: %w", err)
	}

	// Check if this is a subscription charge webhook
	if _, isSubscriptionCharge := webhookData["subscription_charge"]; isSubscriptionCharge {
		// This is a subscription charge webhook - process as payment
		paymentIDStr, ok := webhookData["payment_id"].(string)
		if !ok {
			return fmt.Errorf("invalid subscription charge webhook: missing payment_id")
		}

		paymentID, err := uuid.Parse(paymentIDStr)
		if err != nil {
			return fmt.Errorf("invalid payment ID in subscription charge webhook: %w", err)
		}

		status := webhookData["status"].(string)
		return s.processPaymentWebhook(paymentID, status)
	}

	// Extract payment/subscription ID and status
	paymentIDStr, ok := webhookData["payment_id"].(string)
	if ok {
		// This is a payment webhook
		paymentID, err := uuid.Parse(paymentIDStr)
		if err != nil {
			return fmt.Errorf("invalid payment ID in webhook: %w", err)
		}

		status := webhookData["status"].(string)
		return s.processPaymentWebhook(paymentID, status)
	}

	subscriptionIDStr, ok := webhookData["subscription_id"].(string)
	if ok {
		// This is a subscription webhook
		subscriptionID, err := uuid.Parse(subscriptionIDStr)
		if err != nil {
			return fmt.Errorf("invalid subscription ID in webhook: %w", err)
		}

		status := webhookData["status"].(string)
		return s.processSubscriptionWebhook(subscriptionID, status)
	}

	return fmt.Errorf("invalid webhook payload: missing payment_id or subscription_id")
}

// schedulePaymentWebhook schedules a webhook to be sent after 10 seconds
func (s *MockPaymentService) schedulePaymentWebhook(payment *models.Payment, req *PaymentIntentRequest) {
	time.Sleep(10 * time.Second)

	// Simulate successful payment
	webhookPayload := map[string]interface{}{
		"payment_id":          payment.ID.String(),
		"provider_payment_id": payment.ProviderPaymentID,
		"amount":              req.AmountMinor,
		"currency":            req.Currency,
		"status":              "completed",
		"timestamp":           time.Now().Unix(),
	}

	// Send webhook to our own endpoint
	if err := s.sendWebhook(webhookPayload); err != nil {
		log.Printf("Failed to send mock payment webhook: %v", err)
	}
}

// scheduleSubscriptionWebhook schedules a webhook to be sent after 10 seconds
func (s *MockPaymentService) scheduleSubscriptionWebhook(subscription *models.Subscription, req *SubscriptionIntentRequest) {
	time.Sleep(10 * time.Second)

	// Simulate successful subscription activation
	webhookPayload := map[string]interface{}{
		"subscription_id":          subscription.ID.String(),
		"provider_subscription_id": subscription.ProviderSubscriptionID,
		"amount":                   req.AmountMinor,
		"currency":                 req.Currency,
		"status":                   "active",
		"interval_months":          req.IntervalMonths,
		"timestamp":                time.Now().Unix(),
	}

	// Send webhook to our own endpoint
	if err := s.sendWebhook(webhookPayload); err != nil {
		log.Printf("Failed to send mock subscription webhook: %v", err)
	}
}

// sendWebhook sends a webhook to our own webhook endpoint
func (s *MockPaymentService) sendWebhook(payload map[string]interface{}) error {
	// In a real implementation, this would send an HTTP request to the webhook endpoint
	// For now, we'll just log it and process it directly
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Process the webhook directly since we're in the same process
	return s.ProcessWebhook(payloadBytes, "mock-signature")
}

// processPaymentWebhook processes a payment webhook
func (s *MockPaymentService) processPaymentWebhook(paymentID uuid.UUID, status string) error {
	var payment models.Payment
	if err := s.db.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment status
	switch status {
	case "completed":
		payment.Status = models.PaymentStatusSucceeded
		payment.OccurredAt = &time.Time{}
		*payment.OccurredAt = time.Now()
	case "failed":
		payment.Status = models.PaymentStatusFailed
	case "cancelled":
		payment.Status = models.PaymentStatusCanceled
	default:
		return fmt.Errorf("unknown payment status: %s", status)
	}

	if err := s.db.Save(&payment).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	log.Printf("Mock payment %s status updated to: %s", paymentID, status)
	return nil
}

// processSubscriptionWebhook processes a subscription webhook
func (s *MockPaymentService) processSubscriptionWebhook(subscriptionID uuid.UUID, status string) error {
	var subscription models.Subscription
	if err := s.db.Where("id = ?", subscriptionID).First(&subscription).Error; err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Update subscription status
	switch status {
	case "active":
		subscription.Status = models.SubscriptionStatusActive
		subscription.StartedAt = time.Now()

		// Simulate subscription charge by creating a payment record
		if err := s.createSubscriptionCharge(&subscription); err != nil {
			log.Printf("Failed to create subscription charge: %v", err)
		}

	case "cancelled":
		subscription.Status = models.SubscriptionStatusCanceled
		subscription.CanceledAt = &time.Time{}
		*subscription.CanceledAt = time.Now()
	case "failed":
		subscription.Status = models.SubscriptionStatusFailed
	default:
		return fmt.Errorf("unknown subscription status: %s", status)
	}

	if err := s.db.Save(&subscription).Error; err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	log.Printf("Mock subscription %s status updated to: %s", subscriptionID, status)
	return nil
}

// createSubscriptionCharge creates a payment record for a subscription charge
func (s *MockPaymentService) createSubscriptionCharge(subscription *models.Subscription) error {
	// Generate a realistic-looking provider payment ID for the charge
	providerPaymentID := fmt.Sprintf("mock_charge_%s_%d", uuid.New().String()[:8], time.Now().Unix())

	// Create payment record for this subscription charge
	payment := &models.Payment{
		ID:                uuid.New(),
		Provider:          models.PaymentProviderMock,
		ProviderPaymentID: &providerPaymentID,
		AuthUserID:        &subscription.AuthUserID,
		SubscriptionID:    &subscription.ID,
		AmountMinor:       subscription.AmountMinor,
		Currency:          subscription.Currency,
		Status:            models.PaymentStatusSucceeded,
		OccurredAt:        &time.Time{},
		Meta: map[string]interface{}{
			"subscription_charge": true,
			"subscription_id":     subscription.ID.String(),
			"webhook_processed":   true,
		},
		CreatedAt: time.Now(),
	}
	*payment.OccurredAt = time.Now()

	if err := s.db.Create(payment).Error; err != nil {
		return fmt.Errorf("failed to create subscription charge payment: %w", err)
	}

	// Send payment webhook for the subscription charge
	go s.scheduleSubscriptionChargeWebhook(payment, subscription)

	return nil
}

// scheduleSubscriptionChargeWebhook schedules a webhook for subscription charge
func (s *MockPaymentService) scheduleSubscriptionChargeWebhook(payment *models.Payment, subscription *models.Subscription) {
	time.Sleep(5 * time.Second) // Shorter delay for subscription charges

	// Simulate payment webhook for subscription charge
	webhookPayload := map[string]interface{}{
		"payment_id":          payment.ID.String(),
		"provider_payment_id": payment.ProviderPaymentID,
		"subscription_id":     subscription.ID.String(),
		"amount":              payment.AmountMinor,
		"currency":            payment.Currency,
		"status":              "completed",
		"subscription_charge": true,
		"timestamp":           time.Now().Unix(),
	}

	// Send webhook to our own endpoint
	if err := s.sendWebhook(webhookPayload); err != nil {
		log.Printf("Failed to send mock subscription charge webhook: %v", err)
	}
}
