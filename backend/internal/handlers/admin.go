package handlers

import (
	"net/http"
	"strconv"

	"web-openssl-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get all users (Admin)
// @Description Get paginated list of all users (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Count total users
	var total int64
	h.DB.Model(&models.User{}).Count(&total)

	// Get users
	var users []models.User
	if err := h.DB.Select("id, email, first_name, last_name, role, plan, is_active, usage_count, created_at, updated_at").
		Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	response := map[string]interface{}{
		"users": users,
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

// @Summary Get platform statistics (Admin)
// @Description Get platform-wide statistics (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	var stats struct {
		TotalUsers       int64 `json:"totalUsers"`
		ActiveUsers      int64 `json:"activeUsers"`
		TotalOperations  int64 `json:"totalOperations"`
		FreeUsers        int64 `json:"freeUsers"`
		ProUsers         int64 `json:"proUsers"`
		EnterpriseUsers  int64 `json:"enterpriseUsers"`
		TodayOperations  int64 `json:"todayOperations"`
		MonthOperations  int64 `json:"monthOperations"`
	}

	// Total users
	h.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// Active users
	h.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)

	// Total operations
	h.DB.Model(&models.Operation{}).Count(&stats.TotalOperations)

	// Users by plan
	h.DB.Model(&models.User{}).Where("plan = ?", models.PlanFree).Count(&stats.FreeUsers)
	h.DB.Model(&models.User{}).Where("plan = ?", models.PlanPro).Count(&stats.ProUsers)
	h.DB.Model(&models.User{}).Where("plan = ?", models.PlanEnterprise).Count(&stats.EnterpriseUsers)

	// Today's operations
	h.DB.Model(&models.Operation{}).Where("DATE(created_at) = CURRENT_DATE").Count(&stats.TodayOperations)

	// This month's operations
	h.DB.Model(&models.Operation{}).Where("EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE)").Count(&stats.MonthOperations)

	// Get operation types breakdown
	var operationTypes []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	h.DB.Model(&models.Operation{}).Select("type, count(*) as count").Group("type").Scan(&operationTypes)

	response := map[string]interface{}{
		"overview":        stats,
		"operationTypes":  operationTypes,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update user plan (Admin)
// @Description Update a user's subscription plan (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body map[string]string true "Plan update request"
// @Success 200 {object} map[string]string
// @Router /api/v1/admin/users/{id}/plan [post]
func (h *Handler) UpdateUserPlan(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Plan string `json:"plan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate plan
	var planType models.PlanType
	switch req.Plan {
	case "free":
		planType = models.PlanFree
	case "pro":
		planType = models.PlanPro
	case "enterprise":
		planType = models.PlanEnterprise
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan type"})
		return
	}

	// Update user plan
	result := h.DB.Model(&models.User{}).Where("id = ?", userID).Update("plan", planType)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user plan"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User plan updated successfully"})
}