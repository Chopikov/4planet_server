package payments

import (
	"time"

	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DonationService handles donation creation and management
type DonationService struct {
	db *gorm.DB
}

// NewDonationService creates a new donation service
func NewDonationService() *DonationService {
	return &DonationService{
		db: database.GetDB(),
	}
}

// CreateDonationFromPayment creates a donation record from a completed payment
// This combines payment completion and donation creation in a single transaction
func (s *DonationService) CreateDonationFromPayment(payment *models.Payment) error {
	// Get tree price for the payment currency
	var treePrice models.TreePrice
	if err := s.db.Where("currency = ?", payment.Currency).First(&treePrice).Error; err != nil {
		return err
	}

	// Calculate trees count
	treesCount := int(payment.AmountMinor / treePrice.PriceMinor)

	// Get project ID and referral user ID from payment meta if available
	var projectID *uuid.UUID
	var referralUserID *string
	if meta, ok := payment.Meta.(map[string]interface{}); ok {
		if projectIDStr, exists := meta["project_id"]; exists && projectIDStr != nil {
			if id, ok := projectIDStr.(string); ok {
				if parsedID, err := uuid.Parse(id); err == nil {
					projectID = &parsedID
				}
			}
		}
		// Get referral user ID from payment meta if available
		if refUserID, exists := meta["referral_user_id"]; exists && refUserID != nil {
			if id, ok := refUserID.(string); ok {
				referralUserID = &id
			}
		}
	}

	// Create donation and update payment in a single transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create donation
		donation := &models.Donation{
			ID:             uuid.New(),
			AuthUserID:     *payment.AuthUserID,
			PaymentID:      payment.ID,
			ProjectID:      projectID,
			ReferralUserID: referralUserID,
			TreesCount:     treesCount,
		}

		if err := tx.Create(donation).Error; err != nil {
			return err
		}

		// Update user counters
		updates := map[string]interface{}{
			"total_trees":      gorm.Expr("total_trees + ?", treesCount),
			"donations_count":  gorm.Expr("donations_count + 1"),
			"last_donation_at": time.Now(),
		}

		if err := tx.Model(&models.User{}).Where("auth_user_id = ?", *payment.AuthUserID).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}
