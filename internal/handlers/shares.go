package handlers

import (
	"net/http"

	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/shares"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SharesHandler handles share token operations
type SharesHandler struct {
	shareService *shares.Service
	baseURL      string
}

// NewSharesHandler creates a new shares handler
func NewSharesHandler(shareService *shares.Service, baseURL string) *SharesHandler {
	return &SharesHandler{
		shareService: shareService,
		baseURL:      baseURL,
	}
}

// CreateProfileShareRequest represents a request to create a profile share
type CreateProfileShareRequest struct {
	// No additional fields needed for profile sharing
}

// CreateDonationShareRequest represents a request to create a donation share
type CreateDonationShareRequest struct {
	DonationID uuid.UUID `json:"donation_id" binding:"required"`
}

// ShareTokenResponse represents a share token response
type ShareTokenResponse struct {
	ID        uuid.UUID        `json:"id"`
	Slug      string           `json:"slug"`
	Kind      models.ShareKind `json:"kind"`
	RefID     *uuid.UUID       `json:"ref_id,omitempty"`
	URL       string           `json:"url"`
	CreatedAt string           `json:"created_at"`
}

// CreateProfileShare creates a share token for a user's profile
func (h *SharesHandler) CreateProfileShare(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")
	if authUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	shareToken, err := h.shareService.CreateShareToken(authUserID, models.ShareKindProfile, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share token"})
		return
	}

	response := ShareTokenResponse{
		ID:        shareToken.ID,
		Slug:      shareToken.Slug,
		Kind:      shareToken.Kind,
		RefID:     shareToken.RefID,
		URL:       h.baseURL + "/share/" + shareToken.Slug,
		CreatedAt: shareToken.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

// CreateDonationShare creates a share token for a specific donation
func (h *SharesHandler) CreateDonationShare(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")
	if authUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateDonationShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Verify that the donation belongs to the authenticated user
	// This would require a donation service to check ownership

	shareToken, err := h.shareService.CreateShareToken(authUserID, models.ShareKindDonation, &req.DonationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share token"})
		return
	}

	response := ShareTokenResponse{
		ID:        shareToken.ID,
		Slug:      shareToken.Slug,
		Kind:      shareToken.Kind,
		RefID:     shareToken.RefID,
		URL:       h.baseURL + "/share/" + shareToken.Slug,
		CreatedAt: shareToken.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

// ResolveShare resolves a share token and returns the referral information
func (h *SharesHandler) ResolveShare(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	shareToken, err := h.shareService.ResolveShareToken(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share link not found"})
		return
	}

	// Build enriched response based on share type
	response := gin.H{
		"slug":             shareToken.Slug,
		"kind":             shareToken.Kind,
		"referral_user_id": shareToken.AuthUserID,
		"ref_id":           shareToken.RefID,
		"created_at":       shareToken.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Enrich response based on share type
	switch shareToken.Kind {
	case models.ShareKindDonation:
		if shareToken.RefID != nil {
			// Fetch donation details
			donation, err := h.shareService.GetDonationDetails(*shareToken.RefID)
			if err == nil {
				response["donation"] = donation
			}
		}
	case models.ShareKindProfile:
		// Fetch user profile details
		userProfile, err := h.shareService.GetUserProfile(shareToken.AuthUserID)
		if err == nil {
			response["user_profile"] = userProfile
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetMyShares gets all share tokens for the authenticated user
func (h *SharesHandler) GetMyShares(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")
	if authUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	shareTokens, err := h.shareService.GetUserShareTokens(authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get share tokens"})
		return
	}

	var responses []ShareTokenResponse
	for _, token := range shareTokens {
		response := ShareTokenResponse{
			ID:        token.ID,
			Slug:      token.Slug,
			Kind:      token.Kind,
			RefID:     token.RefID,
			URL:       h.baseURL + "/share/" + token.Slug,
			CreatedAt: token.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, responses)
}

// DeleteShare deletes a share token
func (h *SharesHandler) DeleteShare(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")
	if authUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid share token ID"})
		return
	}

	if err := h.shareService.DeleteShareToken(id, authUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete share token"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetReferralStats gets referral statistics for the authenticated user
func (h *SharesHandler) GetReferralStats(c *gin.Context) {
	authUserID := c.GetString("auth_user_id")
	if authUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := h.shareService.GetReferralStats(authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
