package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	subscriptionUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/subscription"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// SubscriptionHandler handles subscription HTTP requests
type SubscriptionHandler struct {
	subscriptionService *subscriptionUsecase.Service
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionService *subscriptionUsecase.Service) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// GetCurrent handles GET /api/v1/subscriptions/current
func (h *SubscriptionHandler) GetCurrent(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	subscription, err := h.subscriptionService.GetCurrentSubscription(c.Request.Context(), tenantID)
	if err != nil {
		logger.Error("Failed to get current subscription", "error", err.Error())
		response.Error(c, err)
		return
	}

	// If no subscription, tenant is on Free tier
	if subscription == nil {
		response.Success(c, gin.H{
			"subscription": nil,
			"tier":         models.TierFree,
			"message":      "No active subscription - using Free tier",
		})
		return
	}

	response.Success(c, gin.H{
		"subscription": subscription,
	})
}

// GetPlans handles GET /api/v1/subscriptions/plans
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans := h.subscriptionService.GetAvailablePlans()

	response.Success(c, gin.H{
		"plans": plans,
	})
}

// Upgrade handles POST /api/v1/subscriptions/upgrade
func (h *SubscriptionHandler) Upgrade(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	var req struct {
		PlanID       string                 `json:"plan_id" binding:"required"`
		BillingCycle models.BillingCycle    `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid upgrade request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Initiate upgrade
	upgradeResp, err := h.subscriptionService.InitiateUpgrade(
		c.Request.Context(),
		tenantID,
		req.PlanID,
		req.BillingCycle,
	)

	if err != nil {
		logger.Error("Failed to initiate upgrade", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, upgradeResp)
}

// Activate handles POST /api/v1/subscriptions/activate
func (h *SubscriptionHandler) Activate(c *gin.Context) {
	var req struct {
		PayPalSubscriptionID string `json:"paypal_subscription_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid activate request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Activate subscription
	if err := h.subscriptionService.ActivateSubscription(c.Request.Context(), req.PayPalSubscriptionID); err != nil {
		logger.Error("Failed to activate subscription", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Subscription activated successfully",
	})
}

// Cancel handles POST /api/v1/subscriptions/cancel
func (h *SubscriptionHandler) Cancel(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	c.ShouldBindJSON(&req)

	if req.Reason == "" {
		req.Reason = "User requested cancellation"
	}

	// Cancel subscription
	if err := h.subscriptionService.CancelSubscription(c.Request.Context(), tenantID, req.Reason); err != nil {
		logger.Error("Failed to cancel subscription", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Subscription cancelled successfully",
	})
}
