package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Password    string    `json:"-" gorm:"not null"`
	FirstName   string    `json:"firstName" gorm:"not null"`
	LastName    string    `json:"lastName" gorm:"not null"`
	Role        UserRole  `json:"role" gorm:"default:'user'"`
	Plan        PlanType  `json:"plan" gorm:"default:'free'"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	EmailVerified bool    `json:"emailVerified" gorm:"default:false"`
	StripeCustomerID string `json:"stripeCustomerId" gorm:"index"`
	APIKey      string    `json:"apiKey" gorm:"uniqueIndex"`
	UsageCount  int       `json:"usageCount" gorm:"default:0"`
	UsageResetAt time.Time `json:"usageResetAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Organizations []OrganizationMember `json:"organizations"`
	Operations    []Operation          `json:"operations"`
	Subscriptions []Subscription       `json:"subscriptions"`
}

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

type PlanType string

const (
	PlanFree       PlanType = "free"
	PlanPro        PlanType = "pro"
	PlanEnterprise PlanType = "enterprise"
)

type Organization struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Plan        PlanType  `json:"plan" gorm:"default:'free'"`
	StripeCustomerID string `json:"stripeCustomerId" gorm:"index"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	Settings    string    `json:"settings" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Members []OrganizationMember `json:"members"`
}

type OrganizationMember struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	UserID         uint         `json:"userId" gorm:"not null"`
	OrganizationID uint         `json:"organizationId" gorm:"not null"`
	Role           MemberRole   `json:"role" gorm:"default:'member'"`
	JoinedAt       time.Time    `json:"joinedAt"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`

	// Relationships
	User         User         `json:"user"`
	Organization Organization `json:"organization"`
}

type MemberRole string

const (
	MemberRoleMember MemberRole = "member"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleOwner  MemberRole = "owner"
)

type Operation struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"userId" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	Command     string    `json:"command" gorm:"not null"`
	Input       string    `json:"input" gorm:"type:text"`
	Output      string    `json:"output" gorm:"type:text"`
	Status      OpStatus  `json:"status" gorm:"default:'pending'"`
	Error       string    `json:"error" gorm:"type:text"`
	Duration    int       `json:"duration"` // milliseconds
	IPAddress   string    `json:"ipAddress"`
	UserAgent   string    `json:"userAgent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Relationships
	User User `json:"user"`
}

type OpStatus string

const (
	OpStatusPending   OpStatus = "pending"
	OpStatusRunning   OpStatus = "running"
	OpStatusCompleted OpStatus = "completed"
	OpStatusFailed    OpStatus = "failed"
)

type Subscription struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	UserID             uint      `json:"userId" gorm:"not null"`
	StripeSubscriptionID string  `json:"stripeSubscriptionId" gorm:"uniqueIndex;not null"`
	Plan               PlanType  `json:"plan" gorm:"not null"`
	Status             string    `json:"status" gorm:"not null"`
	CurrentPeriodStart time.Time `json:"currentPeriodStart"`
	CurrentPeriodEnd   time.Time `json:"currentPeriodEnd"`
	CanceledAt         *time.Time `json:"canceledAt"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`

	// Relationships
	User User `json:"user"`
}

func GetPlanLimits(plan PlanType) map[string]interface{} {
	limits := make(map[string]interface{})

	switch plan {
	case PlanFree:
		limits["operations_per_month"] = 50
		limits["max_file_size"] = 1024 * 1024 // 1MB
		limits["api_access"] = false
		limits["priority_support"] = false
	case PlanPro:
		limits["operations_per_month"] = 5000
		limits["max_file_size"] = 10 * 1024 * 1024 // 10MB
		limits["api_access"] = true
		limits["priority_support"] = true
	case PlanEnterprise:
		limits["operations_per_month"] = -1 // unlimited
		limits["max_file_size"] = 100 * 1024 * 1024 // 100MB
		limits["api_access"] = true
		limits["priority_support"] = true
		limits["white_label"] = true
		limits["sso"] = true
	}

	return limits
}