package handlers

import (
	"net/http"
	"strconv"

	"web-openssl-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// OperationHandler handles HTTP requests for operation tracking
// Single Responsibility: Only handles HTTP layer concerns
type OperationHandler struct {
	operationService services.OperationService
}

// NewOperationHandler creates a new operation handler
func NewOperationHandler(operationService services.OperationService) *OperationHandler {
	return &OperationHandler{
		operationService: operationService,
	}
}

// GetOperations godoc
// @Summary Get user operations
// @Description Get list of user's operations
// @Tags operations
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/operations/ [get]
func (h *OperationHandler) GetOperations(c *gin.Context) {
	userID := c.GetUint("user_id")

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	operations, err := h.operationService.GetOperations(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get operations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"operations": operations,
		"count":      len(operations),
	})
}

// GetStats godoc
// @Summary Get operation statistics
// @Description Get statistics about user's operations
// @Tags operations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/operations/stats [get]
func (h *OperationHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.operationService.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteOperation godoc
// @Summary Delete an operation
// @Description Delete a specific operation
// @Tags operations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/operations/{id} [delete]
func (h *OperationHandler) DeleteOperation(c *gin.Context) {
	userID := c.GetUint("user_id")

	operationIDStr := c.Param("id")
	operationID, err := strconv.ParseUint(operationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation ID"})
		return
	}

	if err := h.operationService.DeleteOperation(c.Request.Context(), userID, uint(operationID)); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete operation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operation deleted successfully"})
}