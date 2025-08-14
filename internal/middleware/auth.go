package middleware

import (
	"fmt"
	"net/http"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequireAuth middleware that requires authentication
func RequireAuth(authService *auth.Service, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session cookie
		cookie, err := c.Cookie(config.App.CookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Parse session ID
		sessionID, err := uuid.Parse(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Get user by session
		user, err := authService.GetUserBySession(sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		// Check if user is active
		if user.Status != "active" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account not active"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.AuthUserID)
		c.Next()
	}
}

// OptionalAuth middleware that optionally loads user if authenticated
func OptionalAuth(authService *auth.Service, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session cookie
		cookie, err := c.Cookie(config.App.CookieName)
		if err != nil {
			c.Next()
			return
		}

		// Parse session ID
		sessionID, err := uuid.Parse(cookie)
		if err != nil {
			c.Next()
			return
		}

		// Get user by session
		user, err := authService.GetUserBySession(sessionID)
		if err != nil {
			c.Next()
			return
		}

		// Check if user is active
		if user.Status != "active" {
			c.Next()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.AuthUserID)
		c.Next()
	}
}

// AdminAuth middleware for admin routes
func AdminAuth(config *config.Config) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.Admin.Username: config.Admin.Password,
	})
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// LoggingMiddleware logs request details
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		requestID := param.Keys["request_id"]
		if requestID == nil {
			requestID = "unknown"
		}

		return fmt.Sprintf("[%s] %s | %d | %v | %s | %s | %s | %s\n",
			requestID,
			param.Method,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	})
}
