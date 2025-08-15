package handlers

import (
	"net/http"
	"time"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/auth"
	"github.com/4planet/backend/pkg/mailer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *auth.Service
	mailer      mailer.Mailer
	config      *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.Service, mailer mailer.Mailer, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		mailer:      mailer,
		config:      config,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Username    string  `json:"username" binding:"required"`
	Password    string  `json:"password" binding:"required,min=8"`
	DisplayName *string `json:"display_name,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// VerifyEmailRequest represents an email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// NewPasswordRequest represents a new password request
type NewPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if _, err := h.authService.GetUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	if _, err := h.authService.GetUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	passwordHash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user, err := h.authService.CreateUser(req.Email, req.Username, passwordHash, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create email verification token
	expiresAt := time.Now().Add(24 * time.Hour)
	token, err := h.authService.CreateEmailVerificationToken(user.AuthUserID, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification token"})
		return
	}

	// Send verification email
	if err := h.mailer.SendVerificationEmail(user.Email, token.Token); err != nil {
		// Log error but don't fail registration
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.Status(http.StatusCreated)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by login (email or username)
	user, err := h.authService.GetUserByLogin(req.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Get user auth data to check status and password
	userAuth, err := h.authService.GetUserAuthByAuthUserID(user.AuthUserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is active
	if userAuth.Status != models.UserStatusActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not active"})
		return
	}

	// Check password
	if userAuth.PasswordHash == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account has no password set"})
		return
	}
	if !h.authService.CheckPassword(req.Password, *userAuth.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Revoke all existing sessions
	if err := h.authService.RevokeAllUserSessions(user.AuthUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		return
	}

	// Create new session
	expiresAt := time.Now().Add(h.config.App.SessionTTL)
	session, err := h.authService.CreateSession(user.AuthUserID, c.GetHeader("User-Agent"), c.ClientIP(), expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set session cookie
	c.SetCookie(
		h.config.App.CookieName,
		session.ID.String(),
		int(h.config.App.SessionTTL.Seconds()),
		"/",
		h.config.App.CookieDomain,
		h.config.App.CookieSecure,
		h.config.App.CookieHTTPOnly,
	)

	c.JSON(http.StatusOK, user)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.Status(http.StatusNoContent)
		return
	}

	// Get session cookie
	cookie, err := c.Cookie(h.config.App.CookieName)
	if err == nil {
		// Parse and revoke session
		if sessionID, err := uuid.Parse(cookie); err == nil {
			h.authService.RevokeSession(sessionID)
		}
	}

	// Clear cookie
	c.SetCookie(
		h.config.App.CookieName,
		"",
		-1,
		"/",
		h.config.App.CookieDomain,
		h.config.App.CookieSecure,
		h.config.App.CookieHTTPOnly,
	)

	c.Status(http.StatusNoContent)
}

// RequestVerificationEmail handles requesting a new verification email
func (h *AuthHandler) RequestVerificationEmail(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	u := user.(*models.User)

	// Create new verification token
	expiresAt := time.Now().Add(24 * time.Hour)
	token, err := h.authService.CreateEmailVerificationToken(u.AuthUserID, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification token"})
		return
	}

	// Send verification email
	if err := h.mailer.SendVerificationEmail(u.Email, token.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ConfirmEmail handles email verification
func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify token
	token, err := h.authService.VerifyEmailToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Verify user email
	if err := h.authService.VerifyUserEmail(token.AuthUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ForgotPassword handles password reset request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.authService.GetUserByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not
		c.Status(http.StatusNoContent)
		return
	}

	// Create password reset token
	expiresAt := time.Now().Add(1 * time.Hour)
	token, err := h.authService.CreatePasswordResetToken(user.AuthUserID, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	// Send password reset email
	if err := h.mailer.SendPasswordResetEmail(user.Email, token.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req NewPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify token
	token, err := h.authService.VerifyPasswordResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Hash new password
	passwordHash, err := h.authService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Update user password
	if err := h.authService.UpdateUserPassword(token.AuthUserID, passwordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Revoke all user sessions
	if err := h.authService.RevokeAllUserSessions(token.AuthUserID); err != nil {
		// Log error but don't fail password reset
	}

	c.Status(http.StatusNoContent)
}
