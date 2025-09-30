package container

import (
	"web-openssl-backend/internal/config"
	"web-openssl-backend/internal/repository"
	"web-openssl-backend/internal/services"
	"web-openssl-backend/pkg/auth"
	"web-openssl-backend/pkg/billing"

	"gorm.io/gorm"
)

// Container holds all application dependencies
// Implements Dependency Injection pattern for better testability and maintainability
type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB *gorm.DB

	// Repositories (Data Access Layer)
	UserRepo      repository.UserRepository
	OperationRepo repository.OperationRepository

	// Services (Business Logic Layer)
	AuthService      *auth.Service
	BillingService   *billing.Service
	UserService      services.UserService
	OperationService services.OperationService
}

// NewContainer creates and initializes the dependency injection container
// This is the composition root - where all dependencies are wired together
func NewContainer(cfg *config.Config, db *gorm.DB) *Container {
	container := &Container{
		Config: cfg,
		DB:     db,
	}

	// Initialize infrastructure services
	container.AuthService = auth.NewService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	container.BillingService = billing.NewService(cfg.Stripe.SecretKey)

	// Initialize repositories
	container.UserRepo = repository.NewUserRepository(db)
	container.OperationRepo = repository.NewOperationRepository(db)

	// Initialize services (inject repository dependencies)
	container.UserService = services.NewUserService(
		container.UserRepo,
		container.AuthService,
	)
	container.OperationService = services.NewOperationService(
		container.OperationRepo,
		container.UserRepo,
	)

	return container
}

// Close cleans up resources
func (c *Container) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}