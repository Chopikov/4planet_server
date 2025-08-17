package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App struct {
		BaseURL        string
		CookieName     string
		CookieDomain   string
		CookieSecure   bool
		CookieHTTPOnly bool
		SessionTTL     time.Duration
	}

	Database struct {
		DSN string
	}

	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
		From     string
	}

	Email struct {
		Subjects struct {
			Verification  string
			PasswordReset string
		}
		URLs struct {
			BaseURL       string
			VerifyEmail   string
			ResetPassword string
		}
		TeamName string
	}

	CloudPayments struct {
		PublicID string
		Secret   string
	}

	Payments struct {
		DefaultDonationDescription  string
		BaseSubscriptionDescription string
		SubscriptionDescriptions    struct {
			Monthly string
			Yearly  string
			Custom  string
		}
	}

	Log struct {
		Level string
	}

	Admin struct {
		Username string
		Password string
	}
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Debug("No .env file found, using environment variables")
	}

	config := &Config{}

	// App config
	config.App.BaseURL = getEnv("APP_BASE_URL", "http://localhost:8080")
	config.App.CookieName = getEnv("APP_COOKIE_NAME", "session_id")
	config.App.CookieDomain = getEnv("APP_COOKIE_DOMAIN", "")
	config.App.CookieSecure = getEnvBool("APP_COOKIE_SECURE", false)
	config.App.CookieHTTPOnly = true
	config.App.SessionTTL = getEnvDuration("APP_SESSION_TTL", 30*24*time.Hour) // 30 days

	// Database config
	config.Database.DSN = getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/planet?sslmode=disable")

	// SMTP config
	config.SMTP.Host = getEnv("SMTP_HOST", "")
	config.SMTP.Port = getEnvInt("SMTP_PORT", 587)
	config.SMTP.User = getEnv("SMTP_USER", "")
	config.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	config.SMTP.From = getEnv("SMTP_FROM", "noreply@4planet.local")

	// Email config
	config.Email.Subjects.Verification = getEnv("EMAIL_VERIFICATION_SUBJECT", "Verify your email address")
	config.Email.Subjects.PasswordReset = getEnv("EMAIL_PASSWORD_RESET_SUBJECT", "Password reset request")
	config.Email.URLs.BaseURL = getEnv("APP_BASE_URL", "http://localhost:8080")
	config.Email.URLs.VerifyEmail = getEnv("EMAIL_VERIFY_EMAIL_URL", "/verify-email")
	config.Email.URLs.ResetPassword = getEnv("EMAIL_RESET_PASSWORD_URL", "/reset-password")
	config.Email.TeamName = getEnv("EMAIL_TEAM_NAME", "Planet")

	// CloudPayments config
	config.CloudPayments.PublicID = getEnv("CLOUDPAYMENTS_PUBLIC_ID", "")
	config.CloudPayments.Secret = getEnv("CLOUDPAYMENTS_SECRET", "")

	// Payments config
	config.Payments.DefaultDonationDescription = getEnv("PAYMENTS_DEFAULT_DONATION_DESCRIPTION", "Donation to Planet")
	config.Payments.BaseSubscriptionDescription = getEnv("PAYMENTS_BASE_SUBSCRIPTION_DESCRIPTION", "Subscription to Planet")
	config.Payments.SubscriptionDescriptions.Monthly = getEnv("PAYMENTS_SUBSCRIPTION_MONTHLY_DESCRIPTION", "Monthly subscription")
	config.Payments.SubscriptionDescriptions.Yearly = getEnv("PAYMENTS_SUBSCRIPTION_YEARLY_DESCRIPTION", "Yearly subscription")
	config.Payments.SubscriptionDescriptions.Custom = getEnv("PAYMENTS_SUBSCRIPTION_CUSTOM_DESCRIPTION", "Custom subscription")

	// Log config
	config.Log.Level = getEnv("LOG_LEVEL", "info")

	// Admin config
	config.Admin.Username = getEnv("ADMIN_USERNAME", "admin")
	config.Admin.Password = getEnv("ADMIN_PASSWORD", "admin")

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
