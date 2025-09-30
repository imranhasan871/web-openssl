package handlers

import (
	"web-openssl-backend/internal/config"
	"web-openssl-backend/pkg/auth"
	"web-openssl-backend/pkg/billing"
	"web-openssl-backend/pkg/openssl"

	"gorm.io/gorm"
)

type Handler struct {
	DB             *gorm.DB
	Config         *config.Config
	AuthService    *auth.Service
	BillingService *billing.Service
	OpenSSLService *openssl.Service
}

func NewHandler(db *gorm.DB, cfg *config.Config, authService *auth.Service, billingService *billing.Service) *Handler {
	opensslService := openssl.NewService(cfg.OpenSSL.BinaryPath)

	return &Handler{
		DB:             db,
		Config:         cfg,
		AuthService:    authService,
		BillingService: billingService,
		OpenSSLService: opensslService,
	}
}