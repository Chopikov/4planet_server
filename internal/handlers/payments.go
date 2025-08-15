package handlers

import (
	"net/http"

	"github.com/4planet/backend/pkg/payments"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentsHandler handles payment-related requests
type PaymentsHandler struct {
	paymentService *payments.CloudPaymentsService
}

// NewPaymentsHandler creates a new payments handler
func NewPaymentsHandler(paymentService *payments.CloudPaymentsService) *PaymentsHandler {
	return &PaymentsHandler{
		paymentService: paymentService,
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

	response, err := h.paymentService.CreatePaymentIntent(paymentReq, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, response)
}
