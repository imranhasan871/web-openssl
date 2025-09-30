package handlers

import (
	"net/http"
	"strconv"

	"web-openssl-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get operations history
// @Description Get paginated list of user's operations
// @Tags operations
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param type query string false "Filter by operation type"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/operations [get]
func (h *Handler) GetOperations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	opType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Build query
	query := h.DB.Where("user_id = ?", userID)
	if opType != "" {
		query = query.Where("type = ?", opType)
	}

	// Count total operations
	var total int64
	query.Model(&models.Operation{}).Count(&total)

	// Get operations
	var operations []models.Operation
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&operations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch operations"})
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	response := map[string]interface{}{
		"operations": operations,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
			"hasNext":    hasNext,
			"hasPrev":    hasPrev,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get single operation
// @Description Get details of a specific operation
// @Tags operations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation ID"
// @Success 200 {object} models.Operation
// @Failure 404 {object} map[string]string
// @Router /api/v1/operations/{id} [get]
func (h *Handler) GetOperation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	operationID := c.Param("id")

	var operation models.Operation
	if err := h.DB.Where("id = ? AND user_id = ?", operationID, userID).First(&operation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Operation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch operation"})
		}
		return
	}

	c.JSON(http.StatusOK, operation)
}

// @Summary Delete operation
// @Description Delete a specific operation from history
// @Tags operations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/operations/{id} [delete]
func (h *Handler) DeleteOperation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	operationID := c.Param("id")

	result := h.DB.Where("id = ? AND user_id = ?", operationID, userID).Delete(&models.Operation{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete operation"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Operation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operation deleted successfully"})
}