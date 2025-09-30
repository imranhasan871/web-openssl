package services

import (
	"context"
	"encoding/json"
	"errors"

	"web-openssl-backend/internal/models"
	"web-openssl-backend/internal/repository"
)

// operationService implements OperationService interface
// Single Responsibility: Handles operation tracking business logic
type operationService struct {
	operationRepo repository.OperationRepository
	userRepo      repository.UserRepository
}

// NewOperationService creates a new operation service
func NewOperationService(
	operationRepo repository.OperationRepository,
	userRepo repository.UserRepository,
) OperationService {
	return &operationService{
		operationRepo: operationRepo,
		userRepo:      userRepo,
	}
}

func (s *operationService) RecordOperation(
	ctx context.Context,
	userID uint,
	opType, description string,
	metadata map[string]interface{},
) error {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check usage limits based on plan
	if err := s.checkUsageLimits(ctx, user); err != nil {
		return err
	}

	// Convert metadata to JSON for storage in Input field
	metadataJSON, _ := json.Marshal(metadata)

	operation := &models.Operation{
		UserID:  userID,
		Type:    opType,
		Command: description,
		Input:   string(metadataJSON),
		Status:  models.OpStatusCompleted,
	}

	if err := s.operationRepo.Create(ctx, operation); err != nil {
		return err
	}

	// Increment user usage count
	user.UsageCount++
	return s.userRepo.Update(ctx, user)
}

func (s *operationService) GetOperations(ctx context.Context, userID uint, limit int) ([]*models.Operation, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}
	return s.operationRepo.GetByUserID(ctx, userID, 0, limit)
}

func (s *operationService) GetStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	return s.operationRepo.GetStats(ctx, userID)
}

func (s *operationService) DeleteOperation(ctx context.Context, userID, operationID uint) error {
	// Verify operation belongs to user
	operation, err := s.operationRepo.GetByID(ctx, operationID)
	if err != nil {
		return err
	}

	if operation.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.operationRepo.Delete(ctx, operationID)
}

func (s *operationService) checkUsageLimits(ctx context.Context, user *models.User) error {
	// Business logic for usage limits based on plan
	limits := map[models.PlanType]int{
		models.PlanFree:       1000,
		models.PlanPro:        10000,
		models.PlanEnterprise: -1, // unlimited
	}

	limit, ok := limits[user.Plan]
	if !ok {
		limit = limits[models.PlanFree]
	}

	if limit > 0 && user.UsageCount >= limit {
		return errors.New("usage limit exceeded for your plan")
	}

	return nil
}