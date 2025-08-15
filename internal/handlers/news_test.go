package handlers

import (
	"testing"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/pkg/news"
	"github.com/stretchr/testify/assert"
)

func TestNewNewsHandler(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Create a mock news service
	newsService := &news.Service{}

	// Create the handler
	handler := NewNewsHandler(newsService, cfg)

	// Verify the handler was created correctly
	assert.NotNil(t, handler)
	assert.Equal(t, newsService, handler.newsService)
	assert.Equal(t, cfg, handler.config)
}

func TestNewNewsHandler_Consistency(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Create a mock news service
	newsService := &news.Service{}

	// Create multiple handlers
	handler1 := NewNewsHandler(newsService, cfg)
	handler2 := NewNewsHandler(newsService, cfg)

	// Each call should return a new instance
	assert.NotSame(t, handler1, handler2)

	// But both should be valid
	assert.NotNil(t, handler1)
	assert.NotNil(t, handler2)
}
