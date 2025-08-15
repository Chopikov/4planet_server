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
	"github.com/4planet/backend/pkg/auth"
	"github.com/4planet/backend/pkg/donations"
	"github.com/4planet/backend/pkg/mailer"
	"github.com/4planet/backend/pkg/payments"
	"github.com/4planet/backend/pkg/projects"
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
	userHandler := handlers.NewUserHandler(userService, donationService, subscriptionService)
	projectsHandler := handlers.NewProjectsHandler(projectsService, cfg)

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
		}

		// Projects
		projects := v1.Group("/projects")
		{
			projects.GET("", projectsHandler.GetProjects)
			projects.GET("/:id", projectsHandler.GetProject)
		}

		// Donations & Payments
		v1.POST("/payments/intents", middleware.RequireAuth(authService, cfg), func(c *gin.Context) {
			// TODO: Implement payment intent handler
			c.JSON(http.StatusOK, gin.H{"message": "Payment intent endpoint"})
		})

		// Subscriptions
		v1.POST("/subscriptions/intents", middleware.RequireAuth(authService, cfg), func(c *gin.Context) {
			// TODO: Implement subscription intent handler
			c.JSON(http.StatusOK, gin.H{"message": "Subscription intent endpoint"})
		})

		// News
		v1.GET("/news", func(c *gin.Context) {
			// TODO: Implement news handler
			c.JSON(http.StatusOK, gin.H{"message": "News endpoint"})
		})

		// Prices
		v1.GET("/prices", func(c *gin.Context) {
			// TODO: Implement prices handler
			c.JSON(http.StatusOK, gin.H{"message": "Prices endpoint"})
		})

		// Leaderboard
		v1.GET("/leaderboard", func(c *gin.Context) {
			// TODO: Implement leaderboard handler
			c.JSON(http.StatusOK, gin.H{"message": "Leaderboard endpoint"})
		})

		// Achievements
		v1.GET("/achievements", middleware.RequireAuth(authService, cfg), func(c *gin.Context) {
			// TODO: Implement achievements handler
			c.JSON(http.StatusOK, gin.H{"message": "Achievements endpoint"})
		})

		v1.GET("/badges", func(c *gin.Context) {
			// TODO: Implement badges handler
			c.JSON(http.StatusOK, gin.H{"message": "Badges endpoint"})
		})

		// Shares
		v1.POST("/shares/profile", middleware.RequireAuth(authService, cfg), func(c *gin.Context) {
			// TODO: Implement profile share handler
			c.JSON(http.StatusOK, gin.H{"message": "Profile share endpoint"})
		})

		v1.POST("/shares/donation/:donationId", middleware.RequireAuth(authService, cfg), func(c *gin.Context) {
			// TODO: Implement donation share handler
			c.JSON(http.StatusOK, gin.H{"message": "Donation share endpoint"})
		})

		v1.GET("/shares/resolve/:slug", func(c *gin.Context) {
			// TODO: Implement share resolver handler
			c.JSON(http.StatusOK, gin.H{"message": "Share resolver endpoint"})
		})
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

	// Load HTML templates
	router.LoadHTMLGlob("web/**/*.html")

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
