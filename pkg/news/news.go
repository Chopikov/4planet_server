package news

import (
	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"gorm.io/gorm"
)

// NewsFilter represents optional filters for news queries
type NewsFilter struct {
	Type      *models.NewsType
	ProjectID *string
}

type Service struct {
	db *gorm.DB
}

func NewService() *Service {
	return &Service{
		db: database.GetDB(),
	}
}

// GetNews retrieves news items with pagination and optional filters
func (s *Service) GetNews(limit int, offset int, filter *NewsFilter) ([]models.News, int, error) {
	var newsItems []models.News
	var total int64

	// Build the query
	query := s.db.Model(&models.News{})

	// Apply filters if provided
	if filter != nil {
		if filter.Type != nil {
			query = query.Where("type = ?", *filter.Type)
		}
		if filter.ProjectID != nil {
			query = query.Where("project_id = ?", *filter.ProjectID)
		}
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results, ordered by published date (newest first)
	err = query.Limit(limit).Offset(offset).
		Order("published_at DESC NULLS LAST, created_at DESC").
		Find(&newsItems).Error
	if err != nil {
		return nil, 0, err
	}

	return newsItems, int(total), nil
}

// GetNewsByID retrieves a news item by its ID with its project
func (s *Service) GetNewsByID(id string) (*models.News, error) {
	var news models.News
	err := s.db.Where("id = ?", id).
		Preload("Project").
		First(&news).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

// Convenience methods for backward compatibility and common use cases

// GetAllNews retrieves all news items with pagination (no filters)
func (s *Service) GetAllNews(limit int, offset int) ([]models.News, int, error) {
	return s.GetNews(limit, offset, nil)
}

// GetNewsByType retrieves news items by type with pagination
func (s *Service) GetNewsByType(newsType models.NewsType, limit int, offset int) ([]models.News, int, error) {
	filter := &NewsFilter{Type: &newsType}
	return s.GetNews(limit, offset, filter)
}

// GetNewsByProject retrieves news items for a specific project with pagination
func (s *Service) GetNewsByProject(projectID string, limit int, offset int) ([]models.News, int, error) {
	filter := &NewsFilter{ProjectID: &projectID}
	return s.GetNews(limit, offset, filter)
}
