package payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CloudPaymentsService handles CloudPayments integration
type CloudPaymentsService struct {
	db              *gorm.DB
	publicID        string
	secret          string
	baseURL         string
	texts           PaymentTexts
	donationService *DonationService
}

// NewCloudPaymentsService creates a new CloudPayments service
func NewCloudPaymentsService(publicID, secret, baseURL string, paymentTexts PaymentTexts) *CloudPaymentsService {
	return &CloudPaymentsService{
		db:              database.GetDB(),
		publicID:        publicID,
		secret:          secret,
		baseURL:         baseURL,
		texts:           paymentTexts,
		donationService: NewDonationService(),
	}
}

// GetProviderName returns the name of the payment provider
func (s *CloudPaymentsService) GetProviderName() string {
	return "cloudpayments"
}

// PaymentIntentRequest represents a payment intent request
type PaymentIntentRequest struct {
	Provider         string     `json:"provider"`
	AmountMinor      int64      `json:"amount_minor"`
	Currency         string     `json:"currency"`
	SuccessReturnURL string     `json:"success_return_url"`
	FailReturnURL    string     `json:"fail_return_url"`
	Description      *string    `json:"description,omitempty"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty"`
	ReferralUserID   *string    `json:"referral_user_id,omitempty"`
}

// PaymentIntentResponse represents a payment intent response
type PaymentIntentResponse struct {
	Provider        string                 `json:"provider"`
	RedirectURL     string                 `json:"redirect_url"`
	ProviderPayload map[string]interface{} `json:"provider_payload"`
}

// SubscriptionIntentRequest represents a subscription intent request
type SubscriptionIntentRequest struct {
	Provider         string     `json:"provider"`
	AmountMinor      int64      `json:"amount_minor"`
	Currency         string     `json:"currency"`
	SuccessReturnURL string     `json:"success_return_url"`
	FailReturnURL    string     `json:"fail_return_url"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty"`
	IntervalMonths   int        `json:"interval_months"`
	Description      *string    `json:"description,omitempty"`
}

// SubscriptionIntentResponse represents a subscription intent response
type SubscriptionIntentResponse struct {
	Provider        string                 `json:"provider"`
	RedirectURL     string                 `json:"redirect_url"`
	ProviderPayload map[string]interface{} `json:"provider_payload"`
}

// WebhookPayload represents a CloudPayments webhook payload
type WebhookPayload struct {
	Type           string  `json:"Type"`
	TransactionID  string  `json:"TransactionId"`
	Amount         float64 `json:"Amount"`
	Currency       string  `json:"Currency"`
	Status         string  `json:"Status"`
	AccountID      string  `json:"AccountId"`
	OccurredAt     string  `json:"OccurredAt"`
	SubscriptionID *string `json:"SubscriptionId,omitempty"`
	Reason         *string `json:"Reason,omitempty"`
}

// CreatePaymentIntent creates a payment intent for one-time payment
func (s *CloudPaymentsService) CreatePaymentIntent(req *PaymentIntentRequest, authUserID string) (*PaymentIntentResponse, error) {
	// Create payment record
	payment := &models.Payment{
		ID:          uuid.New(),
		Provider:    models.PaymentProviderCloudPayments,
		AuthUserID:  &authUserID,
		AmountMinor: req.AmountMinor,
		Currency:    models.Currency(req.Currency),
		Status:      models.PaymentStatusPending,
		Meta: map[string]interface{}{
			"success_return_url": req.SuccessReturnURL,
			"fail_return_url":    req.FailReturnURL,
			"description":        req.Description,
			"project_id":         req.ProjectID,
			"referral_user_id":   req.ReferralUserID,
		},
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Generate redirect URL (in production, this would be the actual CloudPayments URL)
	redirectURL := fmt.Sprintf("%s/pay/%s", s.baseURL, payment.ID.String())

	// Create provider payload
	providerPayload := map[string]interface{}{
		"publicId": s.publicID,
		"amount":   req.AmountMinor,
		"currency": req.Currency,
		"description": func() string {
			if req.Description != nil {
				return *req.Description
			}
			return s.texts.DefaultDonationDescription
		}(),
		"accountId": authUserID,
		"paymentId": payment.ID.String(),
	}

	return &PaymentIntentResponse{
		Provider:        "cloudpayments",
		RedirectURL:     redirectURL,
		ProviderPayload: providerPayload,
	}, nil
}

// CreateSubscriptionIntent creates a subscription intent for recurring payments
func (s *CloudPaymentsService) CreateSubscriptionIntent(req *SubscriptionIntentRequest, authUserID string) (*SubscriptionIntentResponse, error) {
	// Create subscription record
	subscription := &models.Subscription{
		ID:             uuid.New(),
		AuthUserID:     authUserID,
		Provider:       models.PaymentProviderCloudPayments,
		AmountMinor:    req.AmountMinor,
		Currency:       models.Currency(req.Currency),
		IntervalMonths: req.IntervalMonths,
		Status:         models.SubscriptionStatusPending,
		Meta: map[string]interface{}{
			"success_return_url": req.SuccessReturnURL,
			"fail_return_url":    req.FailReturnURL,
			"project_id":         req.ProjectID,
			"description":        req.Description,
		},
	}

	if err := s.db.Create(subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Generate redirect URL (in production, this would be the actual CloudPayments subscription URL)
	redirectURL := fmt.Sprintf("%s/subscribe/%s", s.baseURL, subscription.ID.String())

	// Determine interval description
	intervalDesc := s.texts.SubscriptionDescriptions.Monthly
	if req.IntervalMonths > 1 {
		if req.IntervalMonths == 12 {
			intervalDesc = s.texts.SubscriptionDescriptions.Yearly
		} else {
			intervalDesc = fmt.Sprintf(s.texts.SubscriptionDescriptions.Custom, req.IntervalMonths)
		}
	}

	// Create provider payload
	providerPayload := map[string]interface{}{
		"publicId": s.publicID,
		"amount":   req.AmountMinor,
		"currency": req.Currency,
		"description": func() string {
			if req.Description != nil {
				return *req.Description
			}
			return fmt.Sprintf("%s %s", intervalDesc, s.texts.BaseSubscriptionDescription)
		}(),
		"accountId":      authUserID,
		"subscriptionId": subscription.ID.String(),
		"interval":       intervalDesc,
		"intervalMonths": req.IntervalMonths,
	}

	return &SubscriptionIntentResponse{
		Provider:        "cloudpayments",
		RedirectURL:     redirectURL,
		ProviderPayload: providerPayload,
	}, nil
}

// ProcessWebhook processes a CloudPayments webhook
func (s *CloudPaymentsService) ProcessWebhook(payload []byte, signature string) error {
	// Verify signature if secret is provided
	if s.secret != "" {
		if !s.verifySignature(payload, signature) {
			return fmt.Errorf("invalid signature")
		}
	}

	// Parse webhook payload
	var webhookPayload WebhookPayload
	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		return fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Create webhook event record
	webhookEvent := &models.WebhookEvent{
		ID:               uuid.New(),
		Provider:         models.PaymentProviderCloudPayments,
		EventType:        webhookPayload.Type,
		EventIdempotency: &webhookPayload.TransactionID,
		RawPayload:       webhookPayload,
		SignatureOK:      s.secret == "" || s.verifySignature(payload, signature),
	}

	// Check for duplicate events
	var existingEvent models.WebhookEvent
	if err := s.db.Where("event_idempotency = ?", webhookPayload.TransactionID).First(&existingEvent).Error; err == nil {
		// Event already processed
		webhookEvent.ProcessedOK = true
		if err := s.db.Create(webhookEvent).Error; err != nil {
			return fmt.Errorf("failed to create duplicate webhook event: %w", err)
		}
		return nil
	}

	// Process the webhook based on type
	if err := s.processWebhookEvent(&webhookPayload); err != nil {
		errStr := err.Error()
		webhookEvent.ProcessingError = &errStr
		if err := s.db.Create(webhookEvent).Error; err != nil {
			return fmt.Errorf("failed to create webhook event: %w", err)
		}
		return fmt.Errorf("failed to process webhook: %w", err)
	}

	webhookEvent.ProcessedOK = true
	if err := s.db.Create(webhookEvent).Error; err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}

	return nil
}

// processWebhookEvent processes different types of webhook events
func (s *CloudPaymentsService) processWebhookEvent(payload *WebhookPayload) error {
	switch payload.Type {
	case "Payment":
		return s.processPaymentEvent(payload)
	case "SubscriptionCharge":
		return s.processSubscriptionChargeEvent(payload)
	case "Refund":
		return s.processRefundEvent(payload)
	default:
		return fmt.Errorf("unknown webhook type: %s", payload.Type)
	}
}

// processPaymentEvent processes a payment event
func (s *CloudPaymentsService) processPaymentEvent(payload *WebhookPayload) error {
	// Find payment by transaction ID
	var payment models.Payment
	if err := s.db.Where("provider_payment_id = ?", payload.TransactionID).First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	var updates map[string]interface{}
	if payload.Status != "Succeeded" {
		log.Printf("Payment failed: %s (%s)", payload.Status, payload.TransactionID)
		updates = map[string]interface{}{
			"status": models.PaymentStatusFailed,
			"meta":   map[string]interface{}{"webhook_processed": true, "reason": payload.Reason},
		}
	} else {
		occurredAt, _ := time.Parse(time.RFC3339, payload.OccurredAt)
		updates = map[string]interface{}{
			"status":      models.PaymentStatusSucceeded,
			"occurred_at": occurredAt,
			"meta":        map[string]interface{}{"webhook_processed": true},
		}
	}

	// Update payment status
	if err := s.db.Model(&payment).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Create donation
	return s.donationService.CreateDonationFromPayment(&payment)
}

// processSubscriptionChargeEvent processes a subscription charge event
func (s *CloudPaymentsService) processSubscriptionChargeEvent(payload *WebhookPayload) error {
	// Find subscription
	var subscription models.Subscription
	if err := s.db.Where("provider_subscription_id = ?", payload.SubscriptionID).First(&subscription).Error; err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	var status models.PaymentStatus
	if payload.Status == "Succeeded" {
		status = models.PaymentStatusSucceeded
	} else {
		status = models.PaymentStatusFailed
	}

	// Create payment record for this charge
	occurredAt, _ := time.Parse(time.RFC3339, payload.OccurredAt)
	amountMinor := int64(payload.Amount * 100) // Assuming amount is in major units
	payment := &models.Payment{
		ID:                uuid.New(),
		Provider:          models.PaymentProviderCloudPayments,
		ProviderPaymentID: &payload.TransactionID,
		AuthUserID:        &subscription.AuthUserID,
		SubscriptionID:    &subscription.ID,
		AmountMinor:       amountMinor,
		Currency:          subscription.Currency,
		Status:            status,
		OccurredAt:        &occurredAt,
		Meta: map[string]interface{}{
			"subscription_charge": true,
			"webhook_processed":   true,
			"reason":              payload.Reason,
		},
	}

	if err := s.db.Create(payment).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// Create donation
	return s.donationService.CreateDonationFromPayment(payment)
}

// processRefundEvent processes a refund event
func (s *CloudPaymentsService) processRefundEvent(payload *WebhookPayload) error {
	// Find payment by transaction ID
	var payment models.Payment
	if err := s.db.Where("provider_payment_id = ?", payload.TransactionID).First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment status to refunded
	occurredAt, _ := time.Parse(time.RFC3339, payload.OccurredAt)
	updates := map[string]interface{}{
		"status":      models.PaymentStatusRefunded,
		"occurred_at": occurredAt,
		"meta":        map[string]interface{}{"refund_reason": payload.Reason, "webhook_processed": true},
	}

	return s.db.Model(&payment).Updates(updates).Error
}

// verifySignature verifies the webhook signature
func (s *CloudPaymentsService) verifySignature(payload []byte, signature string) bool {
	// Create HMAC-SHA256 hash
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}
