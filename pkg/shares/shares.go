package shares

import (
	"fmt"
	"time"

	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides share token functionality
type Service struct {
	db *gorm.DB
}

// NewService creates a new share service
func NewService() *Service {
	return &Service{
		db: database.GetDB(),
	}
}

// CreateShareToken creates a new share token for a user
func (s *Service) CreateShareToken(authUserID string, kind models.ShareKind, refID *uuid.UUID) (*models.ShareToken, error) {
	var shareToken *models.ShareToken

	// Use a transaction to ensure atomic slug generation and token creation
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Generate a unique slug within the transaction
		slug, err := s.generateUniqueSlugInTx(tx, authUserID, kind)
		if err != nil {
			return fmt.Errorf("failed to generate slug: %w", err)
		}

		token := &models.ShareToken{
			ID:         uuid.New(),
			AuthUserID: authUserID,
			Kind:       kind,
			RefID:      refID,
			Slug:       slug,
			CreatedAt:  time.Now(),
		}

		if err := tx.Create(token).Error; err != nil {
			return fmt.Errorf("failed to create share token: %w", err)
		}

		shareToken = token
		return nil
	})

	if err != nil {
		return nil, err
	}

	return shareToken, nil
}

// ResolveShareToken resolves a share token by slug and returns the referral information
func (s *Service) ResolveShareToken(slug string) (*models.ShareToken, error) {
	var shareToken models.ShareToken
	if err := s.db.Where("slug = ?", slug).First(&shareToken).Error; err != nil {
		return nil, fmt.Errorf("share token not found: %w", err)
	}

	return &shareToken, nil
}

// GetUserShareTokens gets all share tokens for a user
func (s *Service) GetUserShareTokens(authUserID string) ([]models.ShareToken, error) {
	var shareTokens []models.ShareToken
	if err := s.db.Where("auth_user_id = ?", authUserID).Find(&shareTokens).Error; err != nil {
		return nil, fmt.Errorf("failed to get share tokens: %w", err)
	}

	return shareTokens, nil
}

// DeleteShareToken deletes a share token
func (s *Service) DeleteShareToken(id uuid.UUID, authUserID string) error {
	result := s.db.Where("id = ? AND auth_user_id = ?", id, authUserID).Delete(&models.ShareToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete share token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("share token not found or access denied")
	}

	return nil
}

// generateUniqueSlugInTx generates a unique slug within a transaction
func (s *Service) generateUniqueSlugInTx(tx *gorm.DB, authUserID string, kind models.ShareKind) (string, error) {
	// Get user info for slug generation
	var user models.User
	if err := tx.Where("auth_user_id = ?", authUserID).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Generate base slug from username or display name
	baseSlug := "user"
	if user.Username != nil {
		baseSlug = *user.Username
	}

	// Add kind to base slug
	baseSlug = fmt.Sprintf("%s-%s", baseSlug, kind)

	// Generate a guaranteed unique slug using short UUID
	// This eliminates the need for retry logic while maintaining readability
	shortUUID := uuid.New().String()[:8] // First 8 characters of UUID
	slug := fmt.Sprintf("%s-%s", baseSlug, shortUUID)

	// Verify uniqueness (should always be true, but good for safety)
	var existingToken models.ShareToken
	if err := tx.Where("slug = ?", slug).First(&existingToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Slug is unique, return it
			return slug, nil
		}
		// Database error, return it
		return "", fmt.Errorf("database error checking slug uniqueness: %w", err)
	}

	// This should never happen with UUID-based slugs, but handle just in case
	// Generate a longer UUID suffix as fallback
	fullUUID := uuid.New().String()
	slug = fmt.Sprintf("%s-%s", baseSlug, fullUUID)

	return slug, nil
}

// GetReferralStats gets referral statistics for a user
func (s *Service) GetReferralStats(authUserID string) (*ReferralStats, error) {
	var stats ReferralStats

	// Count total referrals
	if err := s.db.Model(&models.Donation{}).
		Where("referral_user_id = ?", authUserID).
		Count(&stats.TotalReferrals).Error; err != nil {
		return nil, fmt.Errorf("failed to count referrals: %w", err)
	}

	// Count total trees planted through referrals
	if err := s.db.Model(&models.Donation{}).
		Where("referral_user_id = ?", authUserID).
		Select("COALESCE(SUM(trees_count), 0)").
		Scan(&stats.TotalTreesPlanted).Error; err != nil {
		return nil, fmt.Errorf("failed to count trees planted: %w", err)
	}

	return &stats, nil
}

// GetDonationDetails retrieves donation details for share resolution
func (s *Service) GetDonationDetails(donationID uuid.UUID) (*models.Donation, error) {
	var donation models.Donation
	if err := s.db.Preload("User").Preload("Project").Where("id = ?", donationID).First(&donation).Error; err != nil {
		return nil, fmt.Errorf("donation not found: %w", err)
	}
	return &donation, nil
}

// GetUserProfile retrieves user profile summary for share resolution
func (s *Service) GetUserProfile(authUserID string) (*UserProfileSummary, error) {
	var user models.User
	if err := s.db.Where("auth_user_id = ?", authUserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get user achievements
	var achievements []models.UserAchievement
	s.db.Preload("Achievement").Where("auth_user_id = ?", authUserID).Find(&achievements)

	profile := &UserProfileSummary{
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		AvatarURL:      user.AvatarURL,
		TotalTrees:     user.TotalTrees,
		DonationsCount: user.DonationsCount,
		LastDonationAt: user.LastDonationAt,
		Achievements:   achievements,
	}

	return profile, nil
}

// UserProfileSummary represents a user profile summary for sharing
type UserProfileSummary struct {
	Username       *string                  `json:"username,omitempty"`
	DisplayName    *string                  `json:"display_name,omitempty"`
	AvatarURL      *string                  `json:"avatar_url,omitempty"`
	TotalTrees     int                      `json:"total_trees"`
	DonationsCount int                      `json:"donations_count"`
	LastDonationAt *time.Time               `json:"last_donation_at,omitempty"`
	Achievements   []models.UserAchievement `json:"achievements"`
}

// ReferralStats represents referral statistics for a user
type ReferralStats struct {
	TotalReferrals    int64 `json:"total_referrals"`
	TotalTreesPlanted int64 `json:"total_trees_planted"`
}
