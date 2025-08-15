package subscriptions

import (
	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService() *Service {
	return &Service{
		db: database.GetDB(),
	}
}

// GetUserSubscriptions retrieves subscriptions for a specific user with pagination
func (s *Service) GetUserSubscriptions(authUserID string, limit int, offset int) ([]models.Subscription, int, error) {
	var subscriptions []models.Subscription
	var total int64

	// Get total count
	err := s.db.Model(&models.Subscription{}).Where("auth_user_id = ?", authUserID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = s.db.Where("auth_user_id = ?", authUserID).
		Order("started_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&subscriptions).Error
	return subscriptions, int(total), err
}

// GetSubscriptionByID retrieves a subscription by its ID
func (s *Service) GetSubscriptionByID(id string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := s.db.Where("id = ?", id).
		Preload("User").
		Preload("Payments").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetActiveSubscriptions retrieves all active subscriptions for a user
func (s *Service) GetActiveSubscriptions(authUserID string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := s.db.Where("auth_user_id = ? AND status = ?", authUserID, "active").
		Order("started_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}
