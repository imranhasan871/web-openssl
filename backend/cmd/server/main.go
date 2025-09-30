package main

import (
	"fmt"
	"log"
	"net/http"

	"web-openssl-backend/internal/config"
	"web-openssl-backend/internal/handlers"
	"web-openssl-backend/internal/middleware"
	"web-openssl-backend/internal/models"
	"web-openssl-backend/pkg/auth"
	"web-openssl-backend/pkg/billing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title OpenSSL UI API
// @version 1.0
// @description Production-grade OpenSSL operations API
// @termsOfService https://your-domain.com/terms

// @contact.name API Support
// @contact.url https://your-domain.com/support
// @contact.email support@your-domain.com

// @license.name Commercial License
// @license.url https://your-domain.com/license

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	authService := auth.NewService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	billingService := billing.NewService(cfg.Stripe.SecretKey)

	// Initialize handlers
	h := handlers.NewHandler(db, cfg, authService, billingService)

	// Setup router
	router := setupRouter(cfg, h)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Env)
	log.Printf("Swagger docs available at: http://localhost:%s/swagger/index.html", cfg.Server.Port)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.Operation{},
		&models.Subscription{},
	)
}

func setupRouter(cfg *config.Config, h *handlers.Handler) *gin.Engine {
	// Set gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Rate limiting middleware
	router.Use(middleware.RateLimit(cfg.RateLimit.Requests, cfg.RateLimit.Window))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"version": "1.0.0",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/forgot-password", h.ForgotPassword)
			auth.POST("/reset-password", h.ResetPassword)
			auth.POST("/verify-email", h.VerifyEmail)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(h.AuthService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", h.GetCurrentUser)
				users.PUT("/me", h.UpdateCurrentUser)
				users.DELETE("/me", h.DeleteCurrentUser)
				users.POST("/api-key", h.GenerateAPIKey)
			}

			// OpenSSL routes
			openssl := protected.Group("/openssl")
			{
				// Certificate operations
				certs := openssl.Group("/certificates")
				{
					certs.POST("/generate", h.GenerateCertificate)
					certs.POST("/csr", h.GenerateCSR)
					certs.POST("/parse", h.ParseCertificate)
					certs.POST("/verify", h.VerifyCertificate)
					certs.POST("/convert", h.ConvertCertificate)
				}

				// Key operations
				keys := openssl.Group("/keys")
				{
					keys.POST("/generate", h.GenerateKey)
					keys.POST("/parse", h.ParseKey)
					keys.POST("/convert", h.ConvertKey)
				}

				// Encryption operations
				encrypt := openssl.Group("/encrypt")
				{
					encrypt.POST("/symmetric", h.SymmetricEncrypt)
					encrypt.POST("/asymmetric", h.AsymmetricEncrypt)
					encrypt.POST("/decrypt", h.Decrypt)
				}

				// Hash operations
				hash := openssl.Group("/hash")
				{
					hash.POST("/generate", h.GenerateHash)
					hash.POST("/verify", h.VerifyHash)
					hash.POST("/hmac", h.GenerateHMAC)
				}

				// SSL/TLS testing
				ssl := openssl.Group("/ssl")
				{
					ssl.POST("/test-connection", h.TestSSLConnection)
					ssl.POST("/analyze-certificate", h.AnalyzeSSLCertificate)
				}
			}

			// Billing routes
			billing := protected.Group("/billing")
			{
				billing.GET("/plans", h.GetPlans)
				billing.POST("/subscribe", h.CreateSubscription)
				billing.GET("/subscription", h.GetSubscription)
				billing.POST("/cancel", h.CancelSubscription)
				billing.GET("/invoices", h.GetInvoices)
				billing.GET("/usage", h.GetUsage)
			}

			// Operations history
			operations := protected.Group("/operations")
			{
				operations.GET("/", h.GetOperations)
				operations.GET("/:id", h.GetOperation)
				operations.DELETE("/:id", h.DeleteOperation)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", h.GetAllUsers)
				admin.GET("/stats", h.GetStats)
				admin.POST("/users/:id/plan", h.UpdateUserPlan)
			}
		}
	}

	// Stripe webhook (unprotected)
	router.POST("/webhooks/stripe", h.HandleStripeWebhook)

	// Swagger documentation
	if cfg.Server.Env != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}