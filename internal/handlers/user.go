package handlers

import (
	"net/http"

	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/donations"
	"github.com/4planet/backend/pkg/pagination"
	"github.com/4planet/backend/pkg/subscriptions"
	"github.com/4planet/backend/pkg/user"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-specific requests
type UserHandler struct {
	userService         *user.Service
	donationService     *donations.Service
	subscriptionService *subscriptions.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *user.Service, donationService *donations.Service, subscriptionService *subscriptions.Service) *UserHandler {
	return &UserHandler{
		userService:         userService,
		donationService:     donationService,
		subscriptionService: subscriptionService,
	}
}

// Me returns the current authenticated user
func (h *UserHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetMyDonations returns the current user's donations
func (h *UserHandler) GetMyDonations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	u := user.(*models.User)

	// Extract pagination parameters
	params := pagination.ExtractPagination(c)

	// Get donations for the current user using the donation service
	donations, total, err := h.donationService.GetUserDonations(u.AuthUserID, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donations"})
		return
	}

	response := pagination.NewPaginatedResponse(donations, total, params)
	c.JSON(http.StatusOK, response)
}

// GetMySubscriptions returns the current user's subscriptions
func (h *UserHandler) GetMySubscriptions(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	u := user.(*models.User)

	// Extract pagination parameters
	params := pagination.ExtractPagination(c)

	// Get subscriptions for the current user using the subscription service
	subscriptions, total, err := h.subscriptionService.GetUserSubscriptions(u.AuthUserID, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	response := pagination.NewPaginatedResponse(subscriptions, total, params)
	c.JSON(http.StatusOK, response)
}
