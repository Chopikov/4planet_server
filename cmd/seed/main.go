package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database connection string
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/planet?sslmode=disable"
	}

	// Connect to database without auto-migration
	if err := database.ConnectWithoutMigration(dsn); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	db := database.GetDB()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("🌱 Starting database seeding...")

	// Seed tree prices
	if err := seedTreePrices(ctx, db); err != nil {
		log.Printf("Failed to seed tree prices: %v", err)
	} else {
		log.Println("✅ Tree prices seeded")
	}

	// Seed achievements
	if err := seedAchievements(ctx, db); err != nil {
		log.Printf("Failed to seed achievements: %v", err)
	} else {
		log.Println("✅ Achievements seeded")
	}

	// Seed projects
	if err := seedProjects(ctx, db); err != nil {
		log.Printf("Failed to seed projects: %v", err)
	} else {
		log.Println("✅ Projects seeded")
	}

	// Seed news items
	if err := seedNews(ctx, db); err != nil {
		log.Printf("Failed to seed news: %v", err)
	} else {
		log.Println("✅ News items seeded")
	}

	// Seed media files
	if err := seedMediaFiles(ctx, db); err != nil {
		log.Printf("Failed to seed media files: %v", err)
	} else {
		log.Println("✅ Media files seeded")
	}

	// Seed test users
	if err := seedTestUsers(ctx, db); err != nil {
		log.Printf("Failed to seed test users: %v", err)
	} else {
		log.Println("✅ Test users seeded")
	}

	// Seed payments
	if err := seedPayments(ctx, db); err != nil {
		log.Printf("Failed to seed payments: %v", err)
	} else {
		log.Println("✅ Payments seeded")
	}

	// Seed donations
	if err := seedDonations(ctx, db); err != nil {
		log.Printf("Failed to seed donations: %v", err)
	} else {
		log.Println("✅ Donations seeded")
	}

	// Seed user achievements
	if err := seedUserAchievements(ctx, db); err != nil {
		log.Printf("Failed to seed user achievements: %v", err)
	} else {
		log.Println("✅ User achievements seeded")
	}

	// Seed share tokens
	if err := seedShareTokens(ctx, db); err != nil {
		log.Printf("Failed to seed share tokens: %v", err)
	} else {
		log.Println("✅ Share tokens seeded")
	}

	// Update user stats
	if err := updateUserStats(ctx, db); err != nil {
		log.Printf("Failed to update user stats: %v", err)
	} else {
		log.Println("✅ User stats updated")
	}

	log.Println("🎉 Database seeding completed successfully!")
}

// Helper functions to convert values to pointers
func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }

func seedTreePrices(ctx context.Context, db *gorm.DB) error {
	prices := []models.TreePrice{
		{Currency: "RUB", PriceMinor: 1900},
		{Currency: "KZT", PriceMinor: 950},
		{Currency: "USD", PriceMinor: 2500},
		{Currency: "EUR", PriceMinor: 2300},
	}

	for _, price := range prices {
		price.UpdatedAt = time.Now()
		if err := db.WithContext(ctx).Where("currency = ?", price.Currency).
			Assign(price).FirstOrCreate(&price).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAchievements(ctx context.Context, db *gorm.DB) error {
	achievements := []models.Achievement{
		{
			Code:           "first_tree",
			Title:          "First Tree",
			Description:    stringPtr("Planted your first tree"),
			ThresholdTrees: intPtr(1),
		},
		{
			Code:           "tree_planter",
			Title:          "Tree Planter",
			Description:    stringPtr("Planted 10 trees"),
			ThresholdTrees: intPtr(10),
		},
		{
			Code:           "forest_guardian",
			Title:          "Forest Guardian",
			Description:    stringPtr("Planted 100 trees"),
			ThresholdTrees: intPtr(100),
		},
		{
			Code:           "earth_saver",
			Title:          "Earth Saver",
			Description:    stringPtr("Planted 1000 trees"),
			ThresholdTrees: intPtr(1000),
		},
		{
			Code:        "volunteer",
			Title:       "Volunteer",
			Description: stringPtr("Participated in a planting event"),
		},
	}

	for _, achievement := range achievements {
		achievement.ID = uuid.New()
		if err := db.WithContext(ctx).Where("code = ?", achievement.Code).
			FirstOrCreate(&achievement).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedProjects(ctx context.Context, db *gorm.DB) error {
	projects := []models.Project{
		{
			Title:       "Moscow Forest Restoration",
			Description: stringPtr("Restoring forests around Moscow after recent fires"),
			Status:      models.ProjectStatusInProgress,
			CountryCode: stringPtr("RU"),
			Region:      stringPtr("Moscow Oblast"),
			LocationGeoJSON: map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{37.6173, 55.7558},
			},
			TreesTarget:  intPtr(10000),
			TreesPlanted: intPtr(2500),
			CoverURL:     stringPtr("https://picsum.photos/id/237/536/354"),
		},
		{
			Title:       "Almaty Green Belt",
			Description: stringPtr("Creating green belt around Almaty to improve air quality"),
			Status:      models.ProjectStatusPlanned,
			CountryCode: stringPtr("KZ"),
			Region:      stringPtr("Almaty"),
			LocationGeoJSON: map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{76.9285, 43.2220},
			},
			TreesTarget:  intPtr(5000),
			TreesPlanted: intPtr(0),
			CoverURL:     stringPtr("https://picsum.photos/id/237/536/354"),
		},
		{
			Title:       "Kazakhstan Steppe Restoration",
			Description: stringPtr("Restoring native vegetation in the steppe regions"),
			Status:      models.ProjectStatusCompleted,
			CountryCode: stringPtr("KZ"),
			Region:      stringPtr("Central Kazakhstan"),
			LocationGeoJSON: map[string]interface{}{
				"type": "Polygon",
				"coordinates": [][][]float64{
					{{70.0, 50.0}, {75.0, 50.0}, {75.0, 52.0}, {70.0, 52.0}, {70.0, 50.0}},
				},
			},
			TreesTarget:  intPtr(2000),
			TreesPlanted: intPtr(2000),
			CoverURL:     stringPtr("https://picsum.photos/id/237/536/354"),
		},
	}

	for _, project := range projects {
		project.ID = uuid.New()
		if err := db.WithContext(ctx).Where("title = ?", project.Title).
			FirstOrCreate(&project).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedNews(ctx context.Context, db *gorm.DB) error {
	// Get project IDs for linking
	var moscowProject, almatyProject models.Project
	if err := db.WithContext(ctx).Where("title = ?", "Moscow Forest Restoration").First(&moscowProject).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("title = ?", "Almaty Green Belt").First(&almatyProject).Error; err != nil {
		return err
	}

	newsItems := []models.News{
		{
			Type:        models.NewsTypeUpdate,
			Title:       "Moscow Forest Project Reaches 25% Goal",
			BodyMD:      stringPtr("We are excited to announce that our Moscow Forest Restoration project has reached 25% of its target!"),
			CoverURL:    stringPtr("https://example.com/images/moscow-forest.jpg"),
			ProjectID:   &moscowProject.ID,
			PublishedAt: &time.Time{},
		},
		{
			Type:        models.NewsTypeAchievement,
			Title:       "First 1000 Trees Planted",
			BodyMD:      stringPtr("Congratulations to our community for planting the first 1000 trees!"),
			CoverURL:    stringPtr("https://example.com/images/1000-trees.jpg"),
			PublishedAt: &time.Time{},
		},
		{
			Type:        models.NewsTypeInvite,
			Title:       "Join Our Next Planting Event",
			BodyMD:      stringPtr("Join us this weekend for a community tree planting event in Almaty!"),
			CoverURL:    stringPtr("https://example.com/images/planting-event.jpg"),
			ProjectID:   &almatyProject.ID,
			PublishedAt: &time.Time{},
		},
	}

	for _, news := range newsItems {
		news.ID = uuid.New()
		*news.PublishedAt = time.Now()
		if err := db.WithContext(ctx).Where("title = ?", news.Title).
			FirstOrCreate(&news).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedMediaFiles(ctx context.Context, db *gorm.DB) error {
	// Get project IDs
	var moscowProject, almatyProject, steppeProject models.Project
	if err := db.WithContext(ctx).Where("title = ?", "Moscow Forest Restoration").First(&moscowProject).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("title = ?", "Almaty Green Belt").First(&almatyProject).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("title = ?", "Kazakhstan Steppe Restoration").First(&steppeProject).Error; err != nil {
		return err
	}

	mediaFiles := []models.MediaFile{
		{
			ProjectID: moscowProject.ID,
			Kind:      models.MediaKindImage,
			URL:       "https://example.com/images/moscow-forest-plan.jpg",
			MimeType:  stringPtr("image/jpeg"),
			Title:     stringPtr("Moscow Forest Restoration Plan"),
			AltText:   stringPtr("Aerial view of the planned forest restoration area"),
			Meta: map[string]interface{}{
				"width":  1200,
				"height": 800,
			},
		},
		{
			ProjectID: almatyProject.ID,
			Kind:      models.MediaKindImage,
			URL:       "https://example.com/images/almaty-greenbelt.jpg",
			MimeType:  stringPtr("image/jpeg"),
			Title:     stringPtr("Almaty Green Belt Plan"),
			AltText:   stringPtr("Map showing the planned green belt around Almaty"),
			Meta: map[string]interface{}{
				"width":  1200,
				"height": 800,
			},
		},
		{
			ProjectID: steppeProject.ID,
			Kind:      models.MediaKindImage,
			URL:       "https://example.com/images/steppe-restoration.jpg",
			MimeType:  stringPtr("image/jpeg"),
			Title:     stringPtr("Steppe Restoration Progress"),
			AltText:   stringPtr("Before and after photos of steppe restoration"),
			Meta: map[string]interface{}{
				"width":  1600,
				"height": 900,
			},
		},
	}

	for _, media := range mediaFiles {
		media.ID = uuid.New()
		if err := db.WithContext(ctx).Where("url = ?", media.URL).
			FirstOrCreate(&media).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedTestUsers(ctx context.Context, db *gorm.DB) error {
	// Create UserAuth records first
	userAuths := []models.UserAuth{
		{
			ID:           uuid.New(),
			AuthUserID:   "test-user-1",
			Email:        "john@example.com",
			PasswordHash: nil, // No password for demo
			Status:       models.UserStatusActive,
			VerifiedAt:   &time.Time{},
		},
		{
			ID:           uuid.New(),
			AuthUserID:   "test-user-2",
			Email:        "jane@example.com",
			PasswordHash: nil, // No password for demo
			Status:       models.UserStatusActive,
			VerifiedAt:   &time.Time{},
		},
		{
			ID:           uuid.New(),
			AuthUserID:   "test-user-3",
			Email:        "bob@example.com",
			PasswordHash: nil, // No password for demo
			Status:       models.UserStatusPending,
			VerifiedAt:   nil, // Not verified yet
		},
	}

	for i := range userAuths {
		if i == 0 || i == 1 {
			*userAuths[i].VerifiedAt = time.Now().Add(-time.Duration(i+1) * 24 * time.Hour)
		}
		if err := db.WithContext(ctx).Where("auth_user_id = ?", userAuths[i].AuthUserID).
			FirstOrCreate(&userAuths[i]).Error; err != nil {
			return err
		}
	}

	// Create User records (profile data)
	users := []models.User{
		{
			AuthUserID:     "test-user-1",
			Username:       stringPtr("john_doe"),
			DisplayName:    stringPtr("John Doe"),
			Email:          "john@example.com",
			TotalTrees:     15,
			DonationsCount: 3,
		},
		{
			AuthUserID:     "test-user-2",
			Username:       stringPtr("jane_smith"),
			DisplayName:    stringPtr("Jane Smith"),
			Email:          "jane@example.com",
			TotalTrees:     42,
			DonationsCount: 8,
		},
		{
			AuthUserID:     "test-user-3",
			Username:       stringPtr("bob_wilson"),
			DisplayName:    stringPtr("Bob Wilson"),
			Email:          "bob@example.com",
			TotalTrees:     0,
			DonationsCount: 0,
		},
	}

	for _, user := range users {
		if err := db.WithContext(ctx).Where("auth_user_id = ?", user.AuthUserID).
			FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedPayments(ctx context.Context, db *gorm.DB) error {
	payments := []models.Payment{
		{
			ID:                uuid.New(),
			Provider:          models.PaymentProviderCloudPayments,
			ProviderPaymentID: stringPtr("cp_test_1"),
			AuthUserID:        stringPtr("test-user-1"),
			AmountMinor:       9500, // 5 trees * 1900 RUB
			Currency:          "RUB",
			Status:            models.PaymentStatusSucceeded,
			OccurredAt:        &time.Time{},
		},
		{
			ID:                uuid.New(),
			Provider:          models.PaymentProviderCloudPayments,
			ProviderPaymentID: stringPtr("cp_test_2"),
			AuthUserID:        stringPtr("test-user-1"),
			AmountMinor:       9500, // 10 trees * 950 KZT
			Currency:          "KZT",
			Status:            models.PaymentStatusSucceeded,
			OccurredAt:        &time.Time{},
		},
		{
			ID:                uuid.New(),
			Provider:          models.PaymentProviderCloudPayments,
			ProviderPaymentID: stringPtr("cp_test_3"),
			AuthUserID:        stringPtr("test-user-2"),
			AmountMinor:       50000, // 20 trees * 2500 USD
			Currency:          "USD",
			Status:            models.PaymentStatusSucceeded,
			OccurredAt:        &time.Time{},
		},
		{
			ID:                uuid.New(),
			Provider:          models.PaymentProviderCloudPayments,
			ProviderPaymentID: stringPtr("cp_test_4"),
			AuthUserID:        stringPtr("test-user-2"),
			AmountMinor:       50600, // 22 trees * 2300 EUR
			Currency:          "EUR",
			Status:            models.PaymentStatusSucceeded,
			OccurredAt:        &time.Time{},
		},
	}

	for i := range payments {
		*payments[i].OccurredAt = time.Now().Add(-time.Duration(i) * 24 * time.Hour)
		if err := db.WithContext(ctx).Create(&payments[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDonations(ctx context.Context, db *gorm.DB) error {
	// Get project IDs
	var moscowProject, almatyProject, steppeProject models.Project
	if err := db.WithContext(ctx).Where("title = ?", "Moscow Forest Restoration").First(&moscowProject).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("title = ?", "Almaty Green Belt").First(&almatyProject).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("title = ?", "Kazakhstan Steppe Restoration").First(&steppeProject).Error; err != nil {
		return err
	}

	// Get payment IDs
	var payments []models.Payment
	if err := db.WithContext(ctx).Where("auth_user_id IN ?", []string{"test-user-1", "test-user-2"}).Find(&payments).Error; err != nil {
		return err
	}

	if len(payments) < 4 {
		return fmt.Errorf("expected 4 payments, got %d", len(payments))
	}

	donations := []models.Donation{
		{
			AuthUserID: "test-user-1",
			PaymentID:  payments[0].ID,
			ProjectID:  &moscowProject.ID,
			TreesCount: 5,
		},
		{
			AuthUserID: "test-user-1",
			PaymentID:  payments[1].ID,
			ProjectID:  &almatyProject.ID,
			TreesCount: 10,
		},
		{
			AuthUserID: "test-user-2",
			PaymentID:  payments[2].ID,
			ProjectID:  &steppeProject.ID,
			TreesCount: 20,
		},
		{
			AuthUserID: "test-user-2",
			PaymentID:  payments[3].ID,
			ProjectID:  nil, // General donation
			TreesCount: 22,
		},
	}

	for _, donation := range donations {
		donation.ID = uuid.New()
		if err := db.WithContext(ctx).Create(&donation).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedUserAchievements(ctx context.Context, db *gorm.DB) error {
	// Get achievement IDs
	var firstTree, treePlanter, forestGuardian models.Achievement
	if err := db.WithContext(ctx).Where("code = ?", "first_tree").First(&firstTree).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("code = ?", "tree_planter").First(&treePlanter).Error; err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("code = ?", "forest_guardian").First(&forestGuardian).Error; err != nil {
		return err
	}

	userAchievements := []models.UserAchievement{
		{
			AuthUserID:    "test-user-1",
			AchievementID: firstTree.ID,
			Reason:        stringPtr("First donation"),
		},
		{
			AuthUserID:    "test-user-1",
			AchievementID: treePlanter.ID,
			Reason:        stringPtr("Reached 10 trees"),
		},
		{
			AuthUserID:    "test-user-2",
			AchievementID: firstTree.ID,
			Reason:        stringPtr("First donation"),
		},
		{
			AuthUserID:    "test-user-2",
			AchievementID: forestGuardian.ID,
			Reason:        stringPtr("Reached 100 trees"),
		},
	}

	for _, ua := range userAchievements {
		if err := db.WithContext(ctx).Where("auth_user_id = ? AND achievement_id = ?", ua.AuthUserID, ua.AchievementID).
			FirstOrCreate(&ua).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedShareTokens(ctx context.Context, db *gorm.DB) error {
	// Get a donation ID for linking
	var donation models.Donation
	if err := db.WithContext(ctx).Where("auth_user_id = ?", "test-user-1").First(&donation).Error; err != nil {
		return err
	}

	shareTokens := []models.ShareToken{
		{
			AuthUserID: "test-user-1",
			Kind:       "profile",
			RefID:      nil,
			Slug:       "john-doe-profile",
		},
		{
			AuthUserID: "test-user-2",
			Kind:       "profile",
			RefID:      nil,
			Slug:       "jane-smith-profile",
		},
		{
			AuthUserID: "test-user-1",
			Kind:       "donation",
			RefID:      &donation.ID,
			Slug:       "john-first-donation",
		},
	}

	for _, st := range shareTokens {
		st.ID = uuid.New()
		if err := db.WithContext(ctx).Where("slug = ?", st.Slug).
			FirstOrCreate(&st).Error; err != nil {
			return err
		}
	}
	return nil
}

func updateUserStats(ctx context.Context, db *gorm.DB) error {
	// Update user stats based on actual donations
	return db.WithContext(ctx).Exec(`
		UPDATE users SET 
			total_trees = (
				SELECT COALESCE(SUM(trees_count), 0) 
				FROM donations 
				WHERE donations.auth_user_id = users.auth_user_id
			),
			donations_count = (
				SELECT COUNT(*) 
				FROM donations 
				WHERE donations.auth_user_id = users.auth_user_id
			),
			last_donation_at = (
				SELECT MAX(created_at) 
				FROM donations 
				WHERE donations.auth_user_id = users.auth_user_id
			)
	`).Error
}
