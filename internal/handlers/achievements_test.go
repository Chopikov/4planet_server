package handlers

import (
	"testing"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/pkg/achievements"
	"github.com/stretchr/testify/assert"
)

func TestNewAchievementsHandler(t *testing.T) {
	cfg := &config.Config{}
	achievementsService := &achievements.Service{}

	handler := NewAchievementsHandler(achievementsService, cfg)
	assert.NotNil(t, handler)
	assert.Equal(t, achievementsService, handler.achievementsService)
	assert.Equal(t, cfg, handler.config)
}

func TestAchievementsHandler_GetAchievements(t *testing.T) {
	cfg := &config.Config{}
	achievementsService := &achievements.Service{}

	handler := NewAchievementsHandler(achievementsService, cfg)
	assert.NotNil(t, handler)
}
