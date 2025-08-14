package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service provides authentication functionality
type Service struct {
	db *gorm.DB
}

// NewService creates a new auth service
func NewService() *Service {
	return &Service{
		db: database.GetDB(),
	}
}

// HashPassword hashes a password using bcrypt
func (s *Service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword checks if a password matches a hash
func (s *Service) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a random token for email verification and password reset
func (s *Service) GenerateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CreateSession creates a new session for a user
func (s *Service) CreateSession(authUserID string, userAgent, ipAddr string, expiresAt time.Time) (*models.Session, error) {
	session := &models.Session{
		ID:         uuid.New(),
		AuthUserID: authUserID,
		ExpiresAt:  expiresAt,
		UserAgent:  &userAgent,
		IPAddr:     &ipAddr,
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (s *Service) GetSession(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.db.Where("id = ? AND expires_at > ? AND revoked_at IS NULL", sessionID, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// RevokeSession revokes a session
func (s *Service) RevokeSession(sessionID uuid.UUID) error {
	now := time.Now()
	return s.db.Model(&models.Session{}).Where("id = ?", sessionID).Update("revoked_at", now).Error
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *Service) RevokeAllUserSessions(authUserID string) error {
	now := time.Now()
	return s.db.Model(&models.Session{}).Where("auth_user_id = ?", authUserID).Update("revoked_at", now).Error
}

// GetUserBySession retrieves a user by session ID
func (s *Service) GetUserBySession(sessionID uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.Joins("JOIN sessions ON sessions.auth_user_id = users.auth_user_id").
		Where("sessions.id = ? AND sessions.expires_at > ? AND sessions.revoked_at IS NULL", sessionID, time.Now()).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateEmailVerificationToken creates a new email verification token
func (s *Service) CreateEmailVerificationToken(authUserID string, expiresAt time.Time) (*models.EmailVerificationToken, error) {
	token := &models.EmailVerificationToken{
		ID:         uuid.New(),
		AuthUserID: authUserID,
		Token:      s.GenerateToken(),
		ExpiresAt:  expiresAt,
	}

	if err := s.db.Create(token).Error; err != nil {
		return nil, fmt.Errorf("failed to create email verification token: %w", err)
	}

	return token, nil
}

// VerifyEmailToken verifies an email verification token
func (s *Service) VerifyEmailToken(tokenStr string) (*models.EmailVerificationToken, error) {
	var token models.EmailVerificationToken
	err := s.db.Where("token = ? AND expires_at > ? AND used_at IS NULL", tokenStr, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}

	// Mark token as used
	now := time.Now()
	if err := s.db.Model(&token).Update("used_at", now).Error; err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	return &token, nil
}

// CreatePasswordResetToken creates a new password reset token
func (s *Service) CreatePasswordResetToken(authUserID string, expiresAt time.Time) (*models.PasswordResetToken, error) {
	token := &models.PasswordResetToken{
		ID:         uuid.New(),
		AuthUserID: authUserID,
		Token:      s.GenerateToken(),
		ExpiresAt:  expiresAt,
	}

	if err := s.db.Create(token).Error; err != nil {
		return nil, fmt.Errorf("failed to create password reset token: %w", err)
	}

	return token, nil
}

// VerifyPasswordResetToken verifies a password reset token
func (s *Service) VerifyPasswordResetToken(tokenStr string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	err := s.db.Where("token = ? AND expires_at > ? AND used_at IS NULL", tokenStr, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}

	// Mark token as used
	now := time.Now()
	if err := s.db.Model(&token).Update("used_at", now).Error; err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	return &token, nil
}

// GetUserByEmail retrieves a user by email
func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *Service) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByLogin retrieves a user by email or username
func (s *Service) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ? OR username = ?", login, login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(email, username, passwordHash string, displayName *string) (*models.User, error) {
	user := &models.User{
		AuthUserID:   uuid.New().String(),
		Email:        email,
		Username:     &username,
		PasswordHash: &passwordHash,
		DisplayName:  displayName,
		Status:       models.UserStatusPending,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUserPassword updates a user's password
func (s *Service) UpdateUserPassword(authUserID, passwordHash string) error {
	return s.db.Model(&models.User{}).Where("auth_user_id = ?", authUserID).Update("password_hash", passwordHash).Error
}

// VerifyUserEmail verifies a user's email
func (s *Service) VerifyUserEmail(authUserID string) error {
	now := time.Now()
	return s.db.Model(&models.User{}).Where("auth_user_id = ?", authUserID).Updates(map[string]interface{}{
		"email_verified_at": now,
		"status":            models.UserStatusActive,
	}).Error
}

// CleanupExpiredTokens removes expired tokens
func (s *Service) CleanupExpiredTokens() error {
	now := time.Now()

	// Clean up expired email verification tokens
	if err := s.db.Where("expires_at < ?", now).Delete(&models.EmailVerificationToken{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired email verification tokens: %w", err)
	}

	// Clean up expired password reset tokens
	if err := s.db.Where("expires_at < ?", now).Delete(&models.PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired password reset tokens: %w", err)
	}

	// Clean up expired sessions
	if err := s.db.Where("expires_at < ?", now).Delete(&models.Session{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return nil
}
