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

	// Initialize services
	authService := auth.NewService()
	userService := user.NewService()
	donationService := donations.NewService()
	subscriptionService := subscriptions.NewService()
	projectsService := projects.NewService()
	newsService := news.NewService()
	pricesService := prices.NewService()
	achievementsService := achievements.NewService()
	sharesService := shares.NewService()

	var mailerService mailer.Mailer
	if cfg.SMTP.Host != "" {
		mailerService = mailer.NewSMTPMailer(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.User,
			cfg.SMTP.Password,
			cfg.SMTP.From,
		)
	} else {
		mailerService = mailer.NewNoOpMailer()
	}

	_ = payments.NewCloudPaymentsService(
		cfg.CloudPayments.PublicID,
		cfg.CloudPayments.Secret,
		cfg.App.BaseURL,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, mailerService, cfg)
	userHandler := handlers.NewUserHandler(userService, donationService, subscriptionService, achievementsService)
	projectsHandler := handlers.NewProjectsHandler(projectsService, cfg)
	newsHandler := handlers.NewNewsHandler(newsService, cfg)
	pricesHandler := handlers.NewPricesHandler(pricesService, cfg)
	achievementsHandler := handlers.NewAchievementsHandler(achievementsService, cfg)

	// Initialize share services and handlers
	sharesHandler := handlers.NewSharesHandler(sharesService, cfg.App.BaseURL)

	// Initialize payment services and handlers
	paymentService := payments.NewCloudPaymentsService(
		cfg.CloudPayments.PublicID,
		cfg.CloudPayments.Secret,
		cfg.App.BaseURL,
	)
	paymentsHandler := handlers.NewPaymentsHandler(paymentService)

	// Initialize subscription handlers
	subscriptionsHandler := handlers.NewSubscriptionsHandler(paymentService)

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

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Auth routes
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

		// User profile and data (requires authentication)
		me := v1.Group("/me")
		me.Use(middleware.RequireAuth(authService, cfg))
		{
			me.GET("", userHandler.Me)
			me.GET("/donations", userHandler.GetMyDonations)
			me.GET("/subscriptions", userHandler.GetMySubscriptions)
			me.GET("/achievements", userHandler.GetMyAchievements)
		}

		// Projects
		projects := v1.Group("/projects")
		{
			projects.GET("", projectsHandler.GetProjects)
			projects.GET("/:id", projectsHandler.GetProject)
		}

		news := v1.Group("/news")
		{
			news.GET("", newsHandler.GetNews)
			news.GET("/:id", newsHandler.GetNewsItem)
		}

		// Prices
		prices := v1.Group("/prices")
		{
			prices.GET("", pricesHandler.GetPrices)
			prices.GET("/:currency", pricesHandler.GetPriceByCurrency)
		}

		// Achievements
		achievements := v1.Group("/achievements")
		achievements.Use(middleware.RequireAuth(authService, cfg))
		{
			achievements.GET("", achievementsHandler.GetAchievements)
		}

		// Badges (public catalog of all achievements)
		v1.GET("/badges", achievementsHandler.GetAchievements)

		users := v1.Group("/users")
		users.Use(middleware.RequireAuth(authService, cfg))
		{
			users.GET("/leaderboard", userHandler.GetLeaderboard)
		}

		// Payments
		payments := v1.Group("/payments")
		payments.Use(middleware.RequireAuth(authService, cfg))
		{
			payments.POST("/intents", paymentsHandler.CreatePaymentIntent)
		}

		// Subscriptions
		subscriptions := v1.Group("/subscriptions")
		subscriptions.Use(middleware.RequireAuth(authService, cfg))
		{
			subscriptions.POST("/intents", subscriptionsHandler.CreateSubscriptionIntent)
		}

		// Shares
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

	// Webhooks
	router.POST("/webhooks/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		if provider == "cloudpayments" {
			// TODO: Implement CloudPayments webhook handler
			c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider"})
		}
	})

	// Admin interface
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
