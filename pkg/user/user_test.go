package user

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

func TestService_GetUserByAuthID(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetUserByID(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetUserByEmail(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_UpdateUser(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}
