package handlers

import (
	"net/http"

	"github.com/4planet/backend/pkg/payments"
	"github.com/gin-gonic/gin"
)

// SubscriptionsHandler handles subscription-related requests
type SubscriptionsHandler struct {
	paymentService *payments.CloudPaymentsService
}

// NewSubscriptionsHandler creates a new subscriptions handler
func NewSubscriptionsHandler(paymentService *payments.CloudPaymentsService) *SubscriptionsHandler {
	return &SubscriptionsHandler{
		paymentService: paymentService,
	}
}

// CreateSubscriptionIntent creates a new subscription intent
func (h *SubscriptionsHandler) CreateSubscriptionIntent(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")

	var req struct {
		Provider         string  `json:"provider" binding:"required"`
		AmountMinor      int64   `json:"amount_minor" binding:"required"`
		Currency         string  `json:"currency" binding:"required"`
		SuccessReturnURL string  `json:"success_return_url" binding:"required"`
		FailReturnURL    string  `json:"fail_return_url" binding:"required"`
		Description      *string `json:"description"`
		Interval         string  `json:"interval" binding:"required"`
		IntervalCount    int     `json:"interval_count" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate interval
	if req.Interval != "monthly" && req.Interval != "yearly" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interval. Must be 'monthly' or 'yearly'"})
		return
	}

	// Convert interval to months
	var intervalMonths int
	switch req.Interval {
	case "monthly":
		intervalMonths = req.IntervalCount
	case "yearly":
		intervalMonths = req.IntervalCount * 12
	}

	// Create subscription intent request
	subscriptionReq := &payments.SubscriptionIntentRequest{
		Provider:         req.Provider,
		AmountMinor:      req.AmountMinor,
		Currency:         req.Currency,
		SuccessReturnURL: req.SuccessReturnURL,
		FailReturnURL:    req.FailReturnURL,
		IntervalMonths:   intervalMonths,
		Description:      req.Description,
	}

	// Create subscription intent using CloudPayments service
	response, err := h.paymentService.CreateSubscriptionIntent(subscriptionReq, authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription intent"})
		return
	}

	c.JSON(http.StatusOK, response)
}
