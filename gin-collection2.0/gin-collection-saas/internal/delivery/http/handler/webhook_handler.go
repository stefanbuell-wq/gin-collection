package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	subscriptionUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/subscription"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// WebhookHandler handles webhook HTTP requests
type WebhookHandler struct {
	subscriptionService *subscriptionUsecase.Service
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(subscriptionService *subscriptionUsecase.Service) *WebhookHandler {
	return &WebhookHandler{
		subscriptionService: subscriptionService,
	}
}

// PayPal handles POST /api/v1/webhooks/paypal
func (h *WebhookHandler) PayPal(c *gin.Context) {
	var event subscriptionUsecase.WebhookEvent

	if err := c.ShouldBindJSON(&event); err != nil {
		logger.Error("Invalid webhook payload", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": "Invalid webhook payload",
		})
		return
	}

	logger.Info("Received PayPal webhook", "event_type", event.EventType, "event_id", event.ID)

	// Process webhook event
	if err := h.subscriptionService.HandleWebhookEvent(c.Request.Context(), &event); err != nil {
		logger.Error("Failed to process webhook", "error", err.Error(), "event_type", event.EventType)
		response.Error(c, err)
		return
	}

	// Return 200 OK to acknowledge receipt
	response.Success(c, gin.H{
		"message": "Webhook processed successfully",
	})
}
