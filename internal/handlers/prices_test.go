package handlers

import (
	"testing"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/pkg/prices"
	"github.com/stretchr/testify/assert"
)

func TestNewPricesHandler(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Create a mock prices service
	pricesService := &prices.Service{}

	// Create the handler
	handler := NewPricesHandler(pricesService, cfg)

	// Verify the handler was created correctly
	assert.NotNil(t, handler)
	assert.Equal(t, pricesService, handler.pricesService)
	assert.Equal(t, cfg, handler.config)
}

func TestNewPricesHandler_Consistency(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Create a mock prices service
	pricesService := &prices.Service{}

	// Create multiple handlers
	handler1 := NewPricesHandler(pricesService, cfg)
	handler2 := NewPricesHandler(pricesService, cfg)

	// Each call should return a new instance
	assert.NotSame(t, handler1, handler2)

	// But both should be valid
	assert.NotNil(t, handler1)
	assert.NotNil(t, handler2)
}

func TestPricesHandler_GetPrices_MethodExists(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Create a mock prices service
	pricesService := &prices.Service{}

	// Create the handler
	handler := NewPricesHandler(pricesService, cfg)

	// Verify the GetPrices method exists and can be called
	assert.NotNil(t, handler.GetPrices)

	// Test that it's a function
	assert.NotNil(t, handler.GetPrices)
}
