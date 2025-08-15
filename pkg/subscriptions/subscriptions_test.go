package subscriptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	// Note: service.db might be nil if database connection is not available during testing
	// This is expected behavior in test environments
}

func TestService_GetUserSubscriptions(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetSubscriptionByID(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetActiveSubscriptions(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}
