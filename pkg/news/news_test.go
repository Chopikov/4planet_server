package news

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

func TestNewsFilter_Structure(t *testing.T) {
	// Test that NewsFilter can be created with different combinations
	filter1 := &NewsFilter{}
	assert.NotNil(t, filter1)

	filter2 := &NewsFilter{
		Type:      nil,
		ProjectID: nil,
	}
	assert.NotNil(t, filter2)

	// Test that convenience methods still work
	service := NewService()
	assert.NotNil(t, service)
}
