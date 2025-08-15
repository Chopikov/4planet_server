# Pagination System

This document describes the pagination system implemented in the 4Planet API.

## Overview

The pagination system provides a consistent way to handle large datasets by allowing clients to request data in smaller chunks. It uses offset-based pagination with `limit` and `offset` query parameters.

## Query Parameters

- **`limit`** (optional): Number of items to return per page
  - Default: 20
  - Minimum: 1
  - Maximum: 100
  - Example: `?limit=50`

- **`offset`** (optional): Number of items to skip
  - Default: 0
  - Minimum: 0
  - Example: `?offset=100`

## Response Format

All paginated endpoints return responses in the following format:

```json
{
  "items": [...],
  "total": 150,
  "limit": 20,
  "offset": 0
}
```

- **`items`**: Array of the requested items
- **`total`**: Total number of items available
- **`limit`**: Number of items returned in this response
- **`offset`**: Number of items skipped

## Usage Examples

### Basic Pagination
```
GET /v1/projects?limit=10&offset=0
```

### Next Page
```
GET /v1/projects?limit=10&offset=10
```

### Custom Page Size
```
GET /v1/projects?limit=50&offset=0
```

### Default Values
```
GET /v1/projects
# Equivalent to: GET /v1/projects?limit=20&offset=0
```

## Endpoints with Pagination

The following endpoints support pagination:

- `GET /v1/projects` - List all projects
- `GET /v1/projects/{id}/media` - List media files for a project
- `GET /v1/me/donations` - List user's donations
- `GET /v1/me/subscriptions` - List user's subscriptions
- `GET /v1/news` - News feed
- `GET /v1/leaderboard` - User leaderboard

## Implementation Details

### Backend

The pagination system is implemented using:

1. **`pkg/pagination/pagination.go`** - Core pagination utilities with generic types
2. **Service layer** - Updated to support pagination parameters
3. **Handler layer** - Uses pagination utilities with type-safe generic responses
4. **OpenAPI specification** - Documents pagination parameters and responses

### Frontend

When implementing pagination on the frontend:

1. Start with `offset=0` and your desired `limit`
2. For the next page, increment `offset` by `limit`
3. Continue until `offset >= total`
4. Show pagination controls based on `total`, `limit`, and current `offset`

## Example Frontend Implementation

```javascript
async function fetchProjects(page = 0, pageSize = 20) {
  const offset = page * pageSize;
  const response = await fetch(`/v1/projects?limit=${pageSize}&offset=${offset}`);
  const data = await response.json();
  
  return {
    items: data.items,
    total: data.total,
    currentPage: page,
    totalPages: Math.ceil(data.total / pageSize),
    hasNextPage: (page + 1) * pageSize < data.total,
    hasPrevPage: page > 0
  };
}
```

## Backend Usage Example

```go
func (h *Handler) GetItems(c *gin.Context) {
    params := pagination.ExtractPagination(c)
    items, total, err := h.service.GetItems(params.Limit, params.Offset)
    if err != nil {
        // handle error
        return
    }
    
    // Type-safe generic response - no conversion needed!
    response := pagination.NewPaginatedResponse(items, total, params)
    c.JSON(http.StatusOK, response)
}
```

## Best Practices

1. **Always use pagination** for endpoints that return lists
2. **Set reasonable limits** - don't request more than 100 items at once
3. **Cache results** when possible to improve performance
4. **Handle edge cases** like empty results and last page
5. **Show loading states** during pagination requests
6. **Provide navigation controls** (previous/next, page numbers)
7. **Use the generic pagination system** - no need for manual type conversion
