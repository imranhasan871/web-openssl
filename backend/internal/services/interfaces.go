package services

import (
	"context"
	"web-openssl-backend/internal/models"
)

// UserService defines business logic operations for users
// Following Interface Segregation Principle
type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
	Authenticate(ctx context.Context, email, password string) (*models.User, error)
	GetProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) error
	DeleteAccount(ctx context.Context, userID uint) error
	GenerateAPIKey(ctx context.Context, userID uint) (string, error)
}

// OperationService defines business logic for operations
type OperationService interface {
	RecordOperation(ctx context.Context, userID uint, opType, description string, metadata map[string]interface{}) error
	GetOperations(ctx context.Context, userID uint, limit int) ([]*models.Operation, error)
	GetStats(ctx context.Context, userID uint) (map[string]interface{}, error)
	DeleteOperation(ctx context.Context, userID, operationID uint) error
}

// CertificateService defines business logic for certificate operations
type CertificateService interface {
	Generate(ctx context.Context, req GenerateCertRequest) (*CertificateResponse, error)
	GenerateCSR(ctx context.Context, req GenerateCSRRequest) (*CSRResponse, error)
	Parse(ctx context.Context, certPEM string) (*CertificateInfo, error)
	Verify(ctx context.Context, certPEM, caPEM string) (bool, error)
	Convert(ctx context.Context, input string, fromFormat, toFormat string) (string, error)
}

// EncryptionService defines business logic for encryption operations
type EncryptionService interface {
	EncryptSymmetric(ctx context.Context, data, password string, algorithm string) (string, error)
	EncryptAsymmetric(ctx context.Context, data, publicKeyPEM string) (string, error)
	Decrypt(ctx context.Context, encryptedData, key string) (string, error)
	GenerateHash(ctx context.Context, data, algorithm string) (string, error)
	GenerateHMAC(ctx context.Context, data, key, algorithm string) (string, error)
}

// DTOs (Data Transfer Objects)
type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type UpdateProfileRequest struct {
	FirstName string
	LastName  string
	Email     string
}

type GenerateCertRequest struct {
	CommonName   string
	Organization string
	Country      string
	ValidDays    int
	KeySize      int
}

type CertificateResponse struct {
	Certificate string
	PrivateKey  string
	PublicKey   string
}

type GenerateCSRRequest struct {
	CommonName   string
	Organization string
	Country      string
	KeySize      int
}

type CSRResponse struct {
	CSR        string
	PrivateKey string
}

type CertificateInfo struct {
	Subject    string
	Issuer     string
	NotBefore  string
	NotAfter   string
	SerialNum  string
	PublicKey  string
}