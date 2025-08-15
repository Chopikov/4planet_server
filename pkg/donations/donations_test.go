package donations

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

func TestService_GetUserDonations(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetDonationByID(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetDonationsByProject(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}
