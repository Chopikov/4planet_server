package donations

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

// GetUserDonations retrieves donations for a specific user with pagination
func (s *Service) GetUserDonations(authUserID string, limit int, offset int) ([]models.Donation, int, error) {
	var donations []models.Donation
	var total int64

	// Get total count
	err := s.db.Model(&models.Donation{}).Where("auth_user_id = ?", authUserID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = s.db.Where("auth_user_id = ?", authUserID).
		Preload("Payment").
		Preload("Project").
		Order("created_at DESC").
		Omit("User").
		Limit(limit).
		Offset(offset).
		Find(&donations).Error
	return donations, int(total), err
}

// GetDonationByID retrieves a donation by its ID
func (s *Service) GetDonationByID(id string) (*models.Donation, error) {
	var donation models.Donation
	err := s.db.Where("id = ?", id).
		Preload("Payment").
		Preload("Project").
		First(&donation).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

// GetDonationsByProject retrieves all donations for a specific project
func (s *Service) GetDonationsByProject(projectID string) ([]models.Donation, error) {
	var donations []models.Donation
	err := s.db.Where("project_id = ?", projectID).
		Preload("Payment").
		Preload("User").
		Order("created_at DESC").
		Find(&donations).Error
	return donations, err
}
