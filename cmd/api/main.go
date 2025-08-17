package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/handlers"
	"github.com/4planet/backend/internal/middleware"
	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/achievements"
	"github.com/4planet/backend/pkg/auth"
	"github.com/4planet/backend/pkg/donations"
	"github.com/4planet/backend/pkg/mailer"
	"github.com/4planet/backend/pkg/news"
	"github.com/4planet/backend/pkg/payments"
	"github.com/4planet/backend/pkg/prices"
	"github.com/4planet/backend/pkg/projects"
	"github.com/4planet/backend/pkg/shares"
	"github.com/4planet/backend/pkg/subscriptions"
	"github.com/4planet/backend/pkg/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set log level
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Connect to database
	if err := database.Connect(cfg.Database.DSN); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Auth Entity
	authService := auth.NewService()
	var mailerService mailer.Mailer
	if cfg.SMTP.Host != "" {
		mailerService = mailer.NewSMTPMailer(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.User,
			cfg.SMTP.Password,
			cfg.SMTP.From,
			mailer.EmailTexts{
				Subjects: mailer.EmailSubjects{
					Verification:  cfg.Email.Subjects.Verification,
					PasswordReset: cfg.Email.Subjects.PasswordReset,
				},
				URLs: mailer.EmailURLs{
					BaseURL:       cfg.Email.URLs.BaseURL,
					VerifyEmail:   cfg.Email.URLs.VerifyEmail,
					ResetPassword: cfg.Email.URLs.ResetPassword,
				},
				TeamName: cfg.Email.TeamName,
			},
		)
	} else {
		mailerService = mailer.NewNoOpMailer()
	}
	authHandler := handlers.NewAuthHandler(authService, mailerService, cfg)

	// User Entity
	userService := user.NewService()
	donationService := donations.NewService()
	subscriptionService := subscriptions.NewService()
	achievementsService := achievements.NewService()
	userHandler := handlers.NewUserHandler(userService, donationService, subscriptionService, achievementsService)

	// Projects Entity
	projectsService := projects.NewService()
	projectsHandler := handlers.NewProjectsHandler(projectsService, cfg)

	// News Entity
	newsService := news.NewService()
	newsHandler := handlers.NewNewsHandler(newsService, cfg)

	// Prices Entity
	pricesService := prices.NewService()
	pricesHandler := handlers.NewPricesHandler(pricesService, cfg)

	// Achievements Entity
	achievementsHandler := handlers.NewAchievementsHandler(achievementsService, cfg)

	// Shares Entity
	sharesService := shares.NewService()
	sharesHandler := handlers.NewSharesHandler(sharesService, cfg.App.BaseURL)

	// Payments Entity
	paymentProviderConfigs := map[string]payments.PaymentProviderConfig{
		"cloudpayments": {
			ProviderName: "cloudpayments",
			PublicID:     cfg.CloudPayments.PublicID,
			Secret:       cfg.CloudPayments.Secret,
			BaseURL:      cfg.App.BaseURL,
			Enabled:      true,
			PaymentTexts: payments.PaymentTexts{
				DefaultDonationDescription:  cfg.Payments.DefaultDonationDescription,
				BaseSubscriptionDescription: cfg.Payments.BaseSubscriptionDescription,
				SubscriptionDescriptions: payments.SubscriptionDescriptions{
					Monthly: cfg.Payments.SubscriptionDescriptions.Monthly,
					Yearly:  cfg.Payments.SubscriptionDescriptions.Yearly,
					Custom:  cfg.Payments.SubscriptionDescriptions.Custom,
				},
			},
		},
		"mock": {
			ProviderName: "mock",
			PublicID:     "mock-public-id",
			Secret:       "mock-secret",
			BaseURL:      cfg.App.BaseURL,
			Enabled:      true,
			PaymentTexts: payments.PaymentTexts{
				DefaultDonationDescription:  cfg.Payments.DefaultDonationDescription,
				BaseSubscriptionDescription: cfg.Payments.BaseSubscriptionDescription,
				SubscriptionDescriptions: payments.SubscriptionDescriptions{
					Monthly: cfg.Payments.SubscriptionDescriptions.Monthly,
					Yearly:  cfg.Payments.SubscriptionDescriptions.Yearly,
					Custom:  cfg.Payments.SubscriptionDescriptions.Custom,
				},
			},
		},
		// Add more providers here as they become available
		// "stripe": {
		//     ProviderName: "stripe",
		//     PublicID:     cfg.Stripe.PublicKey,
		//     Secret:       cfg.Stripe.SecretKey,
		//     BaseURL:      cfg.App.BaseURL,
		//     Enabled:      false, // Disabled until implemented
		// },
	}

	paymentFactory := payments.NewPaymentProviderFactory(paymentProviderConfigs)
	paymentsHandler := handlers.NewPaymentsHandler(paymentFactory)
	subscriptionsHandler := handlers.NewSubscriptionsHandler(paymentFactory)

	// Set Gin mode
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// ========================================
	// API ROUTES - Grouped by Entity
	// ========================================
	v1 := router.Group("/v1")
	{
		// ========= AUTH ENTITY =========
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.RequireAuth(authService, cfg), authHandler.Logout)
			auth.POST("/verify-email/request", middleware.RequireAuth(authService, cfg), authHandler.RequestVerificationEmail)
			auth.POST("/verify-email/confirm", authHandler.ConfirmEmail)
			auth.POST("/password/forgot", authHandler.ForgotPassword)
			auth.POST("/password/reset", authHandler.ResetPassword)
		}

		// ========= USER ENTITY =========
		// User profile and data (requires authentication)
		me := v1.Group("/me")
		me.Use(middleware.RequireAuth(authService, cfg))
		{
			me.GET("", userHandler.Me)
			me.GET("/donations", userHandler.GetMyDonations)
			me.GET("/subscriptions", userHandler.GetMySubscriptions)
			me.GET("/achievements", userHandler.GetMyAchievements)
		}

		// User leaderboard (requires authentication)
		users := v1.Group("/users")
		users.Use(middleware.RequireAuth(authService, cfg))
		{
			users.GET("/leaderboard", userHandler.GetLeaderboard)
		}

		// ========= PROJECTS ENTITY =========
		projects := v1.Group("/projects")
		{
			projects.GET("", projectsHandler.GetProjects)
			projects.GET("/:id", projectsHandler.GetProject)
		}

		// ========= NEWS ENTITY =========
		news := v1.Group("/news")
		{
			news.GET("", newsHandler.GetNews)
			news.GET("/:id", newsHandler.GetNewsItem)
		}

		// ========= PRICES ENTITY =========
		prices := v1.Group("/prices")
		{
			prices.GET("", pricesHandler.GetPrices)
			prices.GET("/:currency", pricesHandler.GetPriceByCurrency)
		}

		// ========= ACHIEVEMENTS ENTITY =========
		// Public catalog of all achievements
		v1.GET("/badges", achievementsHandler.GetAchievements)

		// ========= PAYMENTS ENTITY =========
		payments := v1.Group("/payments")
		{
			// Public endpoint to get supported providers
			payments.GET("/providers", paymentsHandler.GetSupportedProviders)

			// Protected endpoints (auth required)
			paymentsProtected := payments.Group("")
			paymentsProtected.Use(middleware.RequireAuth(authService, cfg))
			{
				paymentsProtected.POST("/intents", paymentsHandler.CreatePaymentIntent)
			}
		}

		// ========= SUBSCRIPTIONS ENTITY =========
		subscriptions := v1.Group("/subscriptions")
		subscriptions.Use(middleware.RequireAuth(authService, cfg))
		{
			subscriptions.POST("/intents", subscriptionsHandler.CreateSubscriptionIntent)
		}

		// ========= SHARES ENTITY =========
		shares := v1.Group("/shares")
		{
			// Public endpoint (no auth required)
			shares.GET("/resolve/:slug", sharesHandler.ResolveShare)

			// Protected endpoints (auth required)
			sharesProtected := shares.Group("")
			sharesProtected.Use(middleware.RequireAuth(authService, cfg))
			{
				sharesProtected.POST("/profile", sharesHandler.CreateProfileShare)
				sharesProtected.POST("/donation", sharesHandler.CreateDonationShare)
				sharesProtected.GET("", sharesHandler.GetMyShares)
				sharesProtected.DELETE("/:id", sharesHandler.DeleteShare)
				sharesProtected.GET("/stats", sharesHandler.GetReferralStats)
			}
		}
	}

	// ========================================
	// WEBHOOKS - Payment Provider Callbacks
	// ========================================
	router.POST("/webhooks/:provider", func(c *gin.Context) {
		provider := c.Param("provider")

		switch provider {
		case "cloudpayments":
			// TODO: Implement CloudPayments webhook handler
			c.JSON(http.StatusOK, gin.H{"message": "CloudPayments webhook received"})
		case "mock":
			// Mock provider webhooks are processed internally
			c.JSON(http.StatusOK, gin.H{"message": "Mock provider webhook received"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider"})
		}
	})

	// ========================================
	// ADMIN INTERFACE - Administrative Tools
	// ========================================
	adminRouter := router.Group("/admin")
	adminRouter.Use(middleware.AdminAuth(cfg))
	{
		// TODO: Implement QOR Admin integration
		// Note: QOR Admin requires GORM v1, but we're using GORM v2
		// For now, provide a simple admin interface
		adminRouter.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin.html", gin.H{
				"title": "4Planet Admin",
			})
		})

		adminRouter.GET("/users", func(c *gin.Context) {
			var users []struct {
				models.User
				Status string `json:"status"`
			}
			database.GetDB().Table("users").
				Select("users.*, user_auth.status").
				Joins("JOIN user_auth ON users.auth_user_id = user_auth.auth_user_id").
				Find(&users)
			c.JSON(http.StatusOK, users)
		})

		adminRouter.GET("/projects", func(c *gin.Context) {
			var projects []models.Project
			database.GetDB().Find(&projects)
			c.JSON(http.StatusOK, projects)
		})

		adminRouter.GET("/donations", func(c *gin.Context) {
			var donations []models.Donation
			database.GetDB().Find(&donations)
			c.JSON(http.StatusOK, donations)
		})
	}

	// ========================================
	// STATIC FILES & UTILITIES
	// ========================================

	// Load HTML templates
	router.LoadHTMLGlob("web/**/*.html")

	// Serve OpenAPI spec
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("openapi.yaml")
	})

	// Serve Swagger UI
	router.GET("/docs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "swagger.html", gin.H{
			"title": "4Planet API Documentation",
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logrus.Infof("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}

	logrus.Info("Server exited")
}
