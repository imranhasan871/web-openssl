package handlers

import (
	"io"
	"net/http"
	"time"

	"web-openssl-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get pricing plans
// @Description Get all available pricing plans
// @Tags billing
// @Produce json
// @Success 200 {array} billing.Plan
// @Router /api/v1/billing/plans [get]
func (h *Handler) GetPlans(c *gin.Context) {
	plans := h.BillingService.GetPlans()
	c.JSON(http.StatusOK, plans)
}

// @Summary Create subscription
// @Description Create a new subscription for the user
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body billing.CreateSubscriptionRequest true "Subscription request"
// @Success 200 {object} billing.SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/billing/subscribe [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create Stripe customer if not exists
	if user.StripeCustomerID == "" {
		customer, err := h.BillingService.CreateCustomer(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
			return
		}

		user.StripeCustomerID = customer.ID
		h.DB.Save(&user)
	}

	var req struct {
		PriceID       string `json:"priceId" binding:"required"`
		PaymentMethod string `json:"paymentMethod" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create subscription
	subscription, err := h.BillingService.CreateSubscription(user.StripeCustomerID, req.PriceID, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user plan based on price ID
	var plan models.PlanType
	switch req.PriceID {
	case "price_pro": // This would be the actual Stripe price ID
		plan = models.PlanPro
	case "price_enterprise":
		plan = models.PlanEnterprise
	default:
		plan = models.PlanFree
	}

	user.Plan = plan
	h.DB.Save(&user)

	// Create subscription record
	sub := models.Subscription{
		UserID:               user.ID,
		StripeSubscriptionID: subscription.ID,
		Plan:                 plan,
		Status:               string(subscription.Status),
		CurrentPeriodStart:   time.Unix(subscription.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(subscription.CurrentPeriodEnd, 0),
	}
	h.DB.Create(&sub)

	response := map[string]interface{}{
		"id":                 subscription.ID,
		"status":             subscription.Status,
		"currentPeriodStart": subscription.CurrentPeriodStart,
		"currentPeriodEnd":   subscription.CurrentPeriodEnd,
	}

	if subscription.LatestInvoice != nil {
		response["latestInvoice"] = subscription.LatestInvoice.ID
		if subscription.LatestInvoice.PaymentIntent != nil {
			response["clientSecret"] = subscription.LatestInvoice.PaymentIntent.ClientSecret
		}
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get current subscription
// @Description Get user's current subscription details
// @Tags billing
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Router /api/v1/billing/subscription [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var subscription models.Subscription
	if err := h.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subscription found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// @Summary Cancel subscription
// @Description Cancel user's current subscription
// @Tags billing
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/billing/cancel [post]
func (h *Handler) CancelSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var subscription models.Subscription
	if err := h.DB.Where("user_id = ? AND status = ?", userID, "active").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active subscription found"})
		return
	}

	// Cancel in Stripe
	_, err := h.BillingService.CancelSubscription(subscription.StripeSubscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	// Update subscription status
	subscription.Status = "canceled"
	h.DB.Save(&subscription)

	// Update user plan back to free
	var user models.User
	h.DB.First(&user, userID)
	user.Plan = models.PlanFree
	h.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled successfully"})
}

// @Summary Get invoices
// @Description Get user's billing invoices
// @Tags billing
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/billing/invoices [get]
func (h *Handler) GetInvoices(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.StripeCustomerID == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Get invoices from Stripe
	invoices := []map[string]interface{}{}
	iter := h.BillingService.GetInvoices(user.StripeCustomerID, 10)

	for iter.Next() {
		invoice := iter.Invoice()
		invoices = append(invoices, map[string]interface{}{
			"id":            invoice.ID,
			"amount":        invoice.AmountPaid,
			"currency":      invoice.Currency,
			"status":        invoice.Status,
			"created":       invoice.Created,
			"periodStart":   invoice.PeriodStart,
			"periodEnd":     invoice.PeriodEnd,
			"hostedInvoiceURL": invoice.HostedInvoiceURL,
			"invoicePDF":    invoice.InvoicePDF,
		})
	}

	c.JSON(http.StatusOK, invoices)
}

// @Summary Get usage statistics
// @Description Get user's current usage statistics
// @Tags billing
// @Produce json
// @Security BearerAuth
// @Success 200 {object} billing.UsageResponse
// @Router /api/v1/billing/usage [get]
func (h *Handler) GetUsage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	limits := models.GetPlanLimits(user.Plan)
	monthlyLimit := limits["operations_per_month"].(int)

	response := map[string]interface{}{
		"currentUsage":   user.UsageCount,
		"monthlyLimit":   monthlyLimit,
		"resetDate":      user.UsageResetAt.Format("2006-01-02"),
		"plan":           user.Plan,
		"overageAllowed": user.Plan == models.PlanEnterprise,
		"percentage":     0.0,
	}

	if monthlyLimit > 0 {
		percentage := float64(user.UsageCount) / float64(monthlyLimit) * 100
		response["percentage"] = percentage
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Handle Stripe webhooks
// @Description Handle incoming Stripe webhook events
// @Tags billing
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Stripe signature"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /webhooks/stripe [post]
func (h *Handler) HandleStripeWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	if sigHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe signature"})
		return
	}

	err = h.BillingService.HandleWebhook(payload, sigHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}