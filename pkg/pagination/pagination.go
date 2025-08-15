package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// DefaultLimit is the default number of items per page
const DefaultLimit = 20

// MaxLimit is the maximum number of items per page
const MaxLimit = 100

// ExtractPagination extracts pagination parameters from Gin context
func ExtractPagination(c *gin.Context) PaginationParams {
	limitStr := c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}

// PaginatedResponse represents a paginated response with generic type
type PaginatedResponse[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](items []T, total int, params PaginationParams) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Items:  items,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
}
