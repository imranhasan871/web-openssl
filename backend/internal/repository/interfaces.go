package repository

import (
	"context"
	"web-openssl-backend/internal/models"
)

// UserRepository defines the interface for user data operations
// Following Interface Segregation Principle - specific to user operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
}

// OperationRepository defines the interface for operation data operations
type OperationRepository interface {
	Create(ctx context.Context, operation *models.Operation) error
	GetByID(ctx context.Context, id uint) (*models.Operation, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.Operation, error)
	GetStats(ctx context.Context, userID uint) (map[string]interface{}, error)
	Delete(ctx context.Context, id uint) error
}

// OrganizationRepository defines the interface for organization data operations
type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) error
	GetByID(ctx context.Context, id uint) (*models.Organization, error)
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, id uint) error
	AddMember(ctx context.Context, member *models.OrganizationMember) error
	RemoveMember(ctx context.Context, orgID, userID uint) error
}

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByUserID(ctx context.Context, userID uint) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Cancel(ctx context.Context, id uint) error
}