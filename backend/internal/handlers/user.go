package handlers

import (
	"net/http"

	"web-openssl-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get current user
// @Description Get current authenticated user information
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Router /api/v1/users/me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.Preload("Subscriptions").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// @Summary Update current user
// @Description Update current authenticated user information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Update user request"
// @Success 200 {object} models.User
// @Router /api/v1/users/me [put]
func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update allowed fields
	if firstName, ok := req["firstName"]; ok {
		user.FirstName = firstName
	}
	if lastName, ok := req["lastName"]; ok {
		user.LastName = lastName
	}
	if email, ok := req["email"]; ok {
		user.Email = email
	}

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// @Summary Delete current user
// @Description Delete current authenticated user account
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/v1/users/me [delete]
func (h *Handler) DeleteCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Soft delete the user
	if err := h.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User account deleted successfully"})
}

// @Summary Generate new API key
// @Description Generate a new API key for the current user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/v1/users/api-key [post]
func (h *Handler) GenerateAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate new API key
	apiKey, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	user.APIKey = apiKey
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"apiKey": apiKey})
}