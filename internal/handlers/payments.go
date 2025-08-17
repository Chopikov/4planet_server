package handlers

import (
	"net/http"

	"github.com/4planet/backend/pkg/payments"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentsHandler handles payment-related requests
type PaymentsHandler struct {
	paymentFactory *payments.PaymentProviderFactory
}

// NewPaymentsHandler creates a new payments handler
func NewPaymentsHandler(paymentFactory *payments.PaymentProviderFactory) *PaymentsHandler {
	return &PaymentsHandler{
		paymentFactory: paymentFactory,
	}
}

// CreatePaymentIntent creates a new payment intent
func (h *PaymentsHandler) CreatePaymentIntent(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")

	var req struct {
		Provider         string  `json:"provider" binding:"required"`
		AmountMinor      int64   `json:"amount_minor" binding:"required"`
		Currency         string  `json:"currency" binding:"required"`
		SuccessReturnURL string  `json:"success_return_url" binding:"required"`
		FailReturnURL    string  `json:"fail_return_url" binding:"required"`
		Description      *string `json:"description"`
		ProjectID        *string `json:"project_id"`
		ReferralUserID   *string `json:"referral_user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse project ID if provided
	var projectID *uuid.UUID
	if req.ProjectID != nil {
		if parsedID, err := uuid.Parse(*req.ProjectID); err == nil {
			projectID = &parsedID
		}
	}

	// Get the payment service for the specified provider
	paymentService, err := h.paymentFactory.CreateService(req.Provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported payment provider"})
		return
	}

	paymentReq := &payments.PaymentIntentRequest{
		Provider:         req.Provider,
		AmountMinor:      req.AmountMinor,
		Currency:         req.Currency,
		SuccessReturnURL: req.SuccessReturnURL,
		FailReturnURL:    req.FailReturnURL,
		Description:      req.Description,
		ProjectID:        projectID,
		ReferralUserID:   req.ReferralUserID,
	}

	response, err := paymentService.CreatePaymentIntent(paymentReq, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSupportedProviders returns a list of supported payment providers
func (h *PaymentsHandler) GetSupportedProviders(c *gin.Context) {
	providers := h.paymentFactory.GetSupportedProviders()

	// Get provider details for each supported provider
	var providerDetails []gin.H
	for _, providerName := range providers {
		providerDetails = append(providerDetails, gin.H{
			"name":    providerName,
			"enabled": true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"providers": providerDetails,
	})
}
