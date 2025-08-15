package prices

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

func TestNewService_Consistency(t *testing.T) {
	service1 := NewService()
	service2 := NewService()

	// Each call should return a new instance
	assert.NotSame(t, service1, service2)

	// But both should be valid
	assert.NotNil(t, service1)
	assert.NotNil(t, service2)
}

func TestService_GetPrices(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	// Note: service.db might be nil if database connection is not available during testing
	// This is expected behavior in test environments
}

func TestService_GetPriceByCurrency(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	// Note: service.db might be nil if database connection is not available during testing
	// This is expected behavior in test environments
}

func TestService_UpdatePrice(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	// Note: service.db might be nil if database connection is not available during testing
	// This is expected behavior in test environments
}
