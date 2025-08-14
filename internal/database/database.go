package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/4planet/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(dsn string) error {
	return connectWithMigration(dsn, true)
}

// ConnectWithoutMigration establishes a connection without auto-migration
func ConnectWithoutMigration(dsn string) error {
	return connectWithMigration(dsn, false)
}

// connectWithMigration establishes a connection to the database
func connectWithMigration(dsn string, runMigration bool) error {
	var err error

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate models (for development) if requested
	if runMigration {
		if err := autoMigrate(); err != nil {
			return fmt.Errorf("failed to auto-migrate: %w", err)
		}
	}

	return nil
}

// autoMigrate automatically migrates the database schema
func autoMigrate() error {
	// Skip AutoMigrate in production
	if os.Getenv("ENV") == "production" {
		log.Println("Production environment detected, skipping AutoMigrate")
		return nil
	}

	// Check if migrations have already run
	var count int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'schema_migrations'").Scan(&count)

	if count > 0 {
		log.Println("Migrations detected, skipping AutoMigrate")
		return nil
	}

	log.Println("No migrations found, running AutoMigrate for development")

	models := []interface{}{
		&models.User{},
		&models.Session{},
		&models.EmailVerificationToken{},
		&models.PasswordResetToken{},
		&models.TreePrice{},
		&models.Project{},
		&models.MediaFile{},
		&models.News{},
		&models.Achievement{},
		&models.UserAchievement{},
		&models.Subscription{},
		&models.Payment{},
		&models.Donation{},
		&models.ShareToken{},
		&models.WebhookEvent{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
