package handlers

import (
	"net/http"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/pkg/achievements"
	"github.com/gin-gonic/gin"
)

// AchievementsHandler handles achievement-related requests
type AchievementsHandler struct {
	achievementsService *achievements.Service
	config              *config.Config
}

// NewAchievementsHandler creates a new achievements handler
func NewAchievementsHandler(achievementsService *achievements.Service, config *config.Config) *AchievementsHandler {
	return &AchievementsHandler{
		achievementsService: achievementsService,
		config:              config,
	}
}

// GetAchievements returns all available achievements (catalog)
func (h *AchievementsHandler) GetAchievements(c *gin.Context) {
	achievements, err := h.achievementsService.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}
