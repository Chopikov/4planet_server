package user

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

// GetUserByAuthID retrieves a user by their auth user ID
func (s *Service) GetUserByAuthID(authUserID string) (*models.User, error) {
	var user models.User
	err := s.db.Where("auth_user_id = ?", authUserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their UUID
func (s *Service) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

func (s *Service) GetLeaderboard(limit int, offset int) ([]models.User, int, error) {
	var users []models.User
	var total int64
	err := s.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = s.db.Order("total_trees DESC").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}
