package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"web-openssl-backend/internal/models"
	"web-openssl-backend/internal/repository"
	"web-openssl-backend/pkg/auth"
)

// userService implements UserService interface
// Single Responsibility: Handles user business logic
type userService struct {
	userRepo    repository.UserRepository
	authService *auth.Service
}

// NewUserService creates a new user service
// Dependency Injection: receives dependencies through constructor
func NewUserService(userRepo repository.UserRepository, authService *auth.Service) UserService {
	return &userService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *userService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Generate API key
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Password:     hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         models.RoleUser,
		Plan:         models.PlanFree,
		IsActive:     true,
		APIKey:       apiKey,
		UsageResetAt: time.Now().AddDate(0, 1, 0),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.authService.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" && req.Email != user.Email {
		// Check if new email is available
		existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existing != nil && existing.ID != userID {
			return errors.New("email already in use")
		}
		user.Email = req.Email
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteAccount(ctx context.Context, userID uint) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) GenerateAPIKey(ctx context.Context, userID uint) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return "", err
	}

	user.APIKey = apiKey
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return apiKey, nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}