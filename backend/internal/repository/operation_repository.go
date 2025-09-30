package repository

import (
	"context"
	"web-openssl-backend/internal/models"

	"gorm.io/gorm"
)

// operationRepository implements OperationRepository interface
// Single Responsibility: Only handles operation data persistence
type operationRepository struct {
	db *gorm.DB
}

// NewOperationRepository creates a new operation repository instance
func NewOperationRepository(db *gorm.DB) OperationRepository {
	return &operationRepository{db: db}
}

func (r *operationRepository) Create(ctx context.Context, operation *models.Operation) error {
	return r.db.WithContext(ctx).Create(operation).Error
}

func (r *operationRepository) GetByID(ctx context.Context, id uint) (*models.Operation, error) {
	var operation models.Operation
	err := r.db.WithContext(ctx).First(&operation, id).Error
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

func (r *operationRepository) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.Operation, error) {
	var operations []*models.Operation
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&operations).Error
	return operations, err
}

func (r *operationRepository) GetStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total operations count
	var totalCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Operation{}).
		Where("user_id = ?", userID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["totalOperations"] = totalCount

	// Certificates generated count
	var certCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Operation{}).
		Where("user_id = ? AND type = ?", userID, "certificate").
		Count(&certCount).Error; err != nil {
		return nil, err
	}
	stats["certificatesGenerated"] = certCount

	// Encryption operations count
	var encryptCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.Operation{}).
		Where("user_id = ? AND type = ?", userID, "encryption").
		Count(&encryptCount).Error; err != nil {
		return nil, err
	}
	stats["encryptionOperations"] = encryptCount

	// This month's usage
	stats["usageThisMonth"] = totalCount // Simplified for now

	return stats, nil
}

func (r *operationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Operation{}, id).Error
}