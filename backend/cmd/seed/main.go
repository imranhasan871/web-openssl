package main

import (
	"fmt"
	"log"
	"time"

	"web-openssl-backend/internal/config"
	"web-openssl-backend/internal/models"
	"web-openssl-backend/pkg/auth"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize auth service
	authService := auth.NewService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// Seed users
	if err := seedUsers(db, authService); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	log.Println("Database seeded successfully!")
}

func seedUsers(db *gorm.DB, authService *auth.Service) error {
	users := []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      models.UserRole
		Plan      models.PlanType
	}{
		{
			Email:     "demo@opensslui.com",
			Password:  "demo123",
			FirstName: "Demo",
			LastName:  "User",
			Role:      models.RoleUser,
			Plan:      models.PlanFree,
		},
		{
			Email:     "admin@opensslui.com",
			Password:  "admin123",
			FirstName: "Admin",
			LastName:  "User",
			Role:      models.RoleAdmin,
			Plan:      models.PlanEnterprise,
		},
	}

	for _, u := range users {
		// Check if user already exists
		var existingUser models.User
		if err := db.Where("email = ?", u.Email).First(&existingUser).Error; err == nil {
			log.Printf("User %s already exists, skipping...", u.Email)
			continue
		}

		// Hash password
		hashedPassword, err := authService.HashPassword(u.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Email, err)
		}

		// Generate API key
		apiKey := generateAPIKey(u.Email)

		// Create user
		user := models.User{
			Email:        u.Email,
			Password:     hashedPassword,
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			Role:         u.Role,
			Plan:         u.Plan,
			IsActive:     true,
			APIKey:       apiKey,
			UsageResetAt: time.Now().AddDate(0, 1, 0),
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", u.Email, err)
		}

		log.Printf("Created user: %s (Role: %s, Plan: %s)", u.Email, u.Role, u.Plan)
	}

	return nil
}

func generateAPIKey(email string) string {
	return fmt.Sprintf("sk_test_%s_%d", email[:5], time.Now().Unix())
}