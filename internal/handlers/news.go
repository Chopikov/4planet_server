package handlers

import (
	"net/http"
	"strconv"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/news"
	"github.com/4planet/backend/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsService *news.Service
	config      *config.Config
}

func NewNewsHandler(newsService *news.Service, config *config.Config) *NewsHandler {
	return &NewsHandler{
		newsService: newsService,
		config:      config,
	}
}

// GetNews retrieves a paginated list of news items
func (h *NewsHandler) GetNews(c *gin.Context) {
	// Extract pagination parameters
	params := pagination.ExtractPagination(c)

	// Build filter based on query parameters
	var filter *news.NewsFilter

	// Check for type filter
	if newsType := c.Query("type"); newsType != "" {
		// Validate news type using the enum's validation method
		validType := models.NewsType(newsType)
		if !validType.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news type"})
			return
		}

		filter = &news.NewsFilter{Type: &validType}
	}

	// Check for project filter
	if projectID := c.Query("project_id"); projectID != "" {
		// Validate project ID format
		if _, err := strconv.ParseUint(projectID, 10, 64); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		if filter == nil {
			filter = &news.NewsFilter{ProjectID: &projectID}
		} else {
			filter.ProjectID = &projectID
		}
	}

	// Get news with filters
	newsItems, total, err := h.newsService.GetNews(params.Limit, params.Offset, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
		return
	}

	response := pagination.NewPaginatedResponse(newsItems, total, params)
	c.JSON(http.StatusOK, response)
}

// GetNewsItem retrieves a specific news item by ID
func (h *NewsHandler) GetNewsItem(c *gin.Context) {
	id := c.Param("id")
	newsItem, err := h.newsService.GetNewsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News item not found"})
		return
	}
	c.JSON(http.StatusOK, newsItem)
}
