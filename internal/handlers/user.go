package handlers

import (
	"net/http"
	"time"

	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/achievements"
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
	achievementsService *achievements.Service
}

// LeaderboardUser represents a limited user profile for leaderboard display
type LeaderboardUser struct {
	Username       *string    `json:"username"`
	AvatarURL      *string    `json:"avatar_url"`
	TotalTrees     int        `json:"total_trees"`
	DonationsCount int        `json:"donations_count"`
	LastDonationAt *time.Time `json:"last_donation_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *user.Service, donationService *donations.Service, subscriptionService *subscriptions.Service, achievementsService *achievements.Service) *UserHandler {
	return &UserHandler{
		userService:         userService,
		donationService:     donationService,
		subscriptionService: subscriptionService,
		achievementsService: achievementsService,
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

// GetMyAchievements returns the current user's achievements
func (h *UserHandler) GetMyAchievements(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	u := user.(*models.User)

	// Extract pagination parameters
	params := pagination.ExtractPagination(c)

	// Get achievements for the current user using the achievements service
	achievements, total, err := h.achievementsService.GetUserAchievements(u.AuthUserID, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	response := pagination.NewPaginatedResponse(achievements, int(total), params)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetLeaderboard(c *gin.Context) {
	params := pagination.ExtractPagination(c)
	users, total, err := h.userService.GetLeaderboard(params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	// Convert to limited leaderboard response
	leaderboardUsers := make([]LeaderboardUser, len(users))
	for i, user := range users {
		leaderboardUsers[i] = LeaderboardUser{
			Username:       user.Username,
			AvatarURL:      user.AvatarURL,
			TotalTrees:     user.TotalTrees,
			DonationsCount: user.DonationsCount,
			LastDonationAt: user.LastDonationAt,
			CreatedAt:      user.CreatedAt,
		}
	}

	response := pagination.NewPaginatedResponse(leaderboardUsers, total, params)
	c.JSON(http.StatusOK, response)
}
