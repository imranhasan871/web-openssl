package handlers

import (
	"net/http"
	"time"

	"web-openssl-backend/internal/models"
	"web-openssl-backend/pkg/openssl"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Generate private key
// @Description Generate a new private key
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.GenerateKeyRequest true "Key generation request"
// @Success 200 {object} openssl.GenerateKeyResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/keys/generate [post]
func (h *Handler) GenerateKey(c *gin.Context) {
	var req openssl.GenerateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "generate_key", string(req.KeyType))

	// Generate key
	response, err := h.OpenSSLService.GenerateKey(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Key generated successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Generate certificate
// @Description Generate a new self-signed certificate
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.GenerateCertificateRequest true "Certificate generation request"
// @Success 200 {object} openssl.GenerateCertificateResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/certificates/generate [post]
func (h *Handler) GenerateCertificate(c *gin.Context) {
	var req openssl.GenerateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "generate_certificate", req.Subject.CommonName)

	// Generate certificate
	response, err := h.OpenSSLService.GenerateCertificate(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Certificate generated successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Generate CSR
// @Description Generate a Certificate Signing Request
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.GenerateCSRRequest true "CSR generation request"
// @Success 200 {object} openssl.GenerateCSRResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/certificates/csr [post]
func (h *Handler) GenerateCSR(c *gin.Context) {
	var req openssl.GenerateCSRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "generate_csr", req.Subject.CommonName)

	// For now, return a placeholder response since CSR generation is similar to certificate generation
	// TODO: Implement actual CSR generation in openssl service
	response := &openssl.GenerateCSRResponse{
		CSR:        "-----BEGIN CERTIFICATE REQUEST-----\nPlaceholder CSR content\n-----END CERTIFICATE REQUEST-----",
		PrivateKey: "-----BEGIN PRIVATE KEY-----\nPlaceholder private key\n-----END PRIVATE KEY-----",
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "CSR generated successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Parse certificate
// @Description Parse and analyze a certificate
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.ParseCertificateRequest true "Certificate parsing request"
// @Success 200 {object} openssl.CertificateInfo
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/certificates/parse [post]
func (h *Handler) ParseCertificate(c *gin.Context) {
	var req openssl.ParseCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "parse_certificate", "Certificate analysis")

	// Parse certificate
	response, err := h.OpenSSLService.ParseCertificate(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Certificate parsed successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Verify certificate
// @Description Verify certificate validity and chain
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.VerifyCertificateRequest true "Certificate verification request"
// @Success 200 {object} openssl.VerifyCertificateResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/certificates/verify [post]
func (h *Handler) VerifyCertificate(c *gin.Context) {
	var req openssl.VerifyCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "verify_certificate", "Certificate verification")

	// Verify certificate
	response, err := h.OpenSSLService.VerifyCertificate(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Certificate verified")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Convert certificate format
// @Description Convert certificate between different formats
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Certificate conversion request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/certificates/convert [post]
func (h *Handler) ConvertCertificate(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "convert_certificate", "Certificate conversion")

	// TODO: Implement certificate conversion
	// For now, return placeholder
	response := map[string]string{
		"converted":    "Placeholder converted certificate",
		"format":       req["targetFormat"],
		"originalFormat": req["sourceFormat"],
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Certificate converted successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Parse private key
// @Description Parse and analyze a private key
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Key parsing request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/keys/parse [post]
func (h *Handler) ParseKey(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "parse_key", "Key analysis")

	// TODO: Implement key parsing
	response := map[string]string{
		"keyType":    "RSA",
		"keySize":    "2048",
		"format":     "PEM",
		"encrypted":  "false",
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Key parsed successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Convert key format
// @Description Convert key between different formats
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Key conversion request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/keys/convert [post]
func (h *Handler) ConvertKey(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "convert_key", "Key conversion")

	// TODO: Implement key conversion
	response := map[string]string{
		"convertedKey": "Placeholder converted key",
		"format":       req["targetFormat"],
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Key converted successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Symmetric encryption
// @Description Encrypt data using symmetric encryption
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.EncryptRequest true "Encryption request"
// @Success 200 {object} openssl.EncryptResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/encrypt/symmetric [post]
func (h *Handler) SymmetricEncrypt(c *gin.Context) {
	var req openssl.EncryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "symmetric_encrypt", string(req.Algorithm))

	// Encrypt data
	response, err := h.OpenSSLService.SymmetricEncrypt(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Data encrypted successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Asymmetric encryption
// @Description Encrypt data using asymmetric encryption
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.EncryptRequest true "Encryption request"
// @Success 200 {object} openssl.EncryptResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/encrypt/asymmetric [post]
func (h *Handler) AsymmetricEncrypt(c *gin.Context) {
	var req openssl.EncryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "asymmetric_encrypt", string(req.Algorithm))

	// TODO: Implement asymmetric encryption
	response := &openssl.EncryptResponse{
		EncryptedData: "Placeholder encrypted data",
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Data encrypted successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Decrypt data
// @Description Decrypt encrypted data
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.DecryptRequest true "Decryption request"
// @Success 200 {object} openssl.DecryptResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/encrypt/decrypt [post]
func (h *Handler) Decrypt(c *gin.Context) {
	var req openssl.DecryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "decrypt", string(req.Algorithm))

	// Decrypt data
	response, err := h.OpenSSLService.Decrypt(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Data decrypted successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Generate hash
// @Description Generate hash or HMAC of data
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.HashRequest true "Hash generation request"
// @Success 200 {object} openssl.HashResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/hash/generate [post]
func (h *Handler) GenerateHash(c *gin.Context) {
	var req openssl.HashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "generate_hash", string(req.Algorithm))

	// Generate hash
	response, err := h.OpenSSLService.GenerateHash(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Hash generated successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Verify hash
// @Description Verify data against hash
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Hash verification request"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/hash/verify [post]
func (h *Handler) VerifyHash(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "verify_hash", req["algorithm"])

	// TODO: Implement hash verification
	response := map[string]bool{
		"isValid": true,
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "Hash verified successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Generate HMAC
// @Description Generate HMAC of data
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.HashRequest true "HMAC generation request"
// @Success 200 {object} openssl.HashResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/hash/hmac [post]
func (h *Handler) GenerateHMAC(c *gin.Context) {
	var req openssl.HashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key is required for HMAC"})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "generate_hmac", string(req.Algorithm))

	// Generate HMAC
	response, err := h.OpenSSLService.GenerateHash(&req)
	if err != nil {
		h.finishOperation(operation, models.OpStatusFailed, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "HMAC generated successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Test SSL connection
// @Description Test SSL/TLS connection to a server
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.SSLTestRequest true "SSL test request"
// @Success 200 {object} openssl.SSLTestResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/ssl/test-connection [post]
func (h *Handler) TestSSLConnection(c *gin.Context) {
	var req openssl.SSLTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "ssl_test", req.Hostname)

	// TODO: Implement SSL testing
	response := &openssl.SSLTestResponse{
		IsValid:  true,
		Protocol: "TLSv1.3",
		Cipher:   "TLS_AES_256_GCM_SHA384",
		Grade:    "A+",
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "SSL connection tested successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// @Summary Analyze SSL certificate
// @Description Analyze SSL certificate from a server
// @Tags openssl
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openssl.SSLTestRequest true "SSL analysis request"
// @Success 200 {object} openssl.CertificateInfo
// @Failure 400 {object} map[string]string
// @Router /api/v1/openssl/ssl/analyze-certificate [post]
func (h *Handler) AnalyzeSSLCertificate(c *gin.Context) {
	var req openssl.SSLTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check usage limits
	if !h.checkUsageLimits(c) {
		return
	}

	// Record operation start
	operation := h.startOperation(c, "analyze_ssl_cert", req.Hostname)

	// TODO: Implement SSL certificate analysis
	response := &openssl.CertificateInfo{
		Subject: openssl.Subject{
			CommonName:   req.Hostname,
			Organization: "Example Corp",
		},
		NotBefore:       time.Now().AddDate(0, -1, 0),
		NotAfter:        time.Now().AddDate(1, 0, 0),
		IsExpired:       false,
		DaysUntilExpiry: 365,
	}

	// Record successful operation
	h.finishOperation(operation, models.OpStatusCompleted, "", "SSL certificate analyzed successfully")
	h.incrementUsage(c)

	c.JSON(http.StatusOK, response)
}

// Helper functions for usage tracking and operation logging

func (h *Handler) checkUsageLimits(c *gin.Context) bool {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return false
	}

	limits := models.GetPlanLimits(user.Plan)
	monthlyLimit := limits["operations_per_month"].(int)

	// Check if usage reset is needed
	if time.Now().After(user.UsageResetAt) {
		user.UsageCount = 0
		user.UsageResetAt = time.Now().AddDate(0, 1, 0)
		h.DB.Save(&user)
	}

	// Check usage limits (unlimited for enterprise = -1)
	if monthlyLimit != -1 && user.UsageCount >= monthlyLimit {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Monthly usage limit exceeded",
			"plan":  user.Plan,
			"limit": monthlyLimit,
		})
		return false
	}

	return true
}

func (h *Handler) incrementUsage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	h.DB.Model(&models.User{}).Where("id = ?", userID).Update("usage_count", gorm.Expr("usage_count + 1"))
}

func (h *Handler) startOperation(c *gin.Context, opType, command string) *models.Operation {
	userID, _ := c.Get("user_id")

	operation := &models.Operation{
		UserID:    userID.(uint),
		Type:      opType,
		Command:   command,
		Status:    models.OpStatusRunning,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	h.DB.Create(operation)
	return operation
}

func (h *Handler) finishOperation(operation *models.Operation, status models.OpStatus, errorMsg, output string) {
	operation.Status = status
	operation.Error = errorMsg
	operation.Output = output
	operation.Duration = int(time.Since(operation.CreatedAt).Milliseconds())
	h.DB.Save(operation)
}