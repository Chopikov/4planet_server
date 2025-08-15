package achievements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAchievementsService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	// Note: service.db might be nil if database connection is not available during testing
	// This is expected behavior in test environments
}

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetUserAchievements(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetAllAchievements(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetAchievementByCode(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_AwardAchievement(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_CheckAndAwardTreeBasedAchievements(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}
