package achievements

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

// GetUserAchievements retrieves achievements for a specific user with pagination
func (s *Service) GetUserAchievements(authUserID string, limit, offset int) ([]models.UserAchievement, int64, error) {
	var userAchievements []models.UserAchievement
	var total int64

	// Count total user achievements
	if err := s.db.Model(&models.UserAchievement{}).Where("auth_user_id = ?", authUserID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get user achievements with pagination and preload achievement details
	err := s.db.Where("auth_user_id = ?", authUserID).
		Preload("Achievement").
		Order("awarded_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&userAchievements).Error

	if err != nil {
		return nil, 0, err
	}

	return userAchievements, total, nil
}

// GetAllAchievements retrieves all available achievements (catalog)
func (s *Service) GetAllAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement

	err := s.db.Order("threshold_trees ASC, title ASC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}

	return achievements, nil
}

// GetAchievementByCode retrieves an achievement by its code
func (s *Service) GetAchievementByCode(code string) (*models.Achievement, error) {
	var achievement models.Achievement

	err := s.db.Where("code = ?", code).First(&achievement).Error
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

// AwardAchievement awards an achievement to a user
func (s *Service) AwardAchievement(authUserID string, achievementCode string, reason *string) error {
	// Check if user already has this achievement
	var existing models.UserAchievement
	err := s.db.Where("auth_user_id = ? AND achievement_id = (SELECT id FROM achievements WHERE code = ?)",
		authUserID, achievementCode).First(&existing).Error

	if err == nil {
		// Achievement already awarded
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Get the achievement
	var achievement models.Achievement
	if err := s.db.Where("code = ?", achievementCode).First(&achievement).Error; err != nil {
		return err
	}

	// Create user achievement record
	userAchievement := models.UserAchievement{
		AuthUserID:    authUserID,
		AchievementID: achievement.ID,
		Reason:        reason,
	}

	return s.db.Create(&userAchievement).Error
}

// CheckAndAwardTreeBasedAchievements checks if a user qualifies for tree-based achievements
func (s *Service) CheckAndAwardTreeBasedAchievements(authUserID string, totalTrees int) error {
	var achievements []models.Achievement

	// Get all achievements with threshold_trees that the user might qualify for
	err := s.db.Where("threshold_trees IS NOT NULL AND threshold_trees <= ?", totalTrees).Find(&achievements).Error
	if err != nil {
		return err
	}

	// Award each qualifying achievement
	for _, achievement := range achievements {
		if err := s.AwardAchievement(authUserID, achievement.Code, nil); err != nil {
			return err
		}
	}

	return nil
}
