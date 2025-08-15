package projects

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

// GetProjects retrieves projects with pagination and returns projects and total count
func (s *Service) GetProjects(limit int, offset int) ([]models.Project, int, error) {
	var projects []models.Project
	var total int64

	// Get total count
	err := s.db.Model(&models.Project{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = s.db.Limit(limit).Offset(offset).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, int(total), nil
}

// GetProjectByID retrieves a project by its ID with its media files
func (s *Service) GetProjectByID(id string) (*models.Project, error) {
	var project models.Project
	err := s.db.Where("id = ?", id).
		Preload("MediaFiles").
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}
