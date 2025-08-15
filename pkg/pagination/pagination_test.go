package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExtractPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "default values",
			queryParams:    "",
			expectedLimit:  DefaultLimit,
			expectedOffset: 0,
		},
		{
			name:           "custom limit and offset",
			queryParams:    "?limit=50&offset=100",
			expectedLimit:  50,
			expectedOffset: 100,
		},
		{
			name:           "limit exceeds max",
			queryParams:    "?limit=200&offset=0",
			expectedLimit:  MaxLimit,
			expectedOffset: 0,
		},
		{
			name:           "negative offset",
			queryParams:    "?limit=10&offset=-10",
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "invalid limit",
			queryParams:    "?limit=abc&offset=0",
			expectedLimit:  DefaultLimit,
			expectedOffset: 0,
		},
		{
			name:           "invalid offset",
			queryParams:    "?limit=10&offset=abc",
			expectedLimit:  10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin context
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test"+tt.queryParams, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Extract pagination
			params := ExtractPagination(c)

			// Assert results
			assert.Equal(t, tt.expectedLimit, params.Limit)
			assert.Equal(t, tt.expectedOffset, params.Offset)
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	total := 10
	params := PaginationParams{Limit: 5, Offset: 0}

	response := NewPaginatedResponse(items, total, params)

	assert.Equal(t, items, response.Items)
	assert.Equal(t, total, response.Total)
	assert.Equal(t, params.Limit, response.Limit)
	assert.Equal(t, params.Offset, response.Offset)
}
