package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	botanicalUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/botanical"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// BotanicalHandler handles botanical HTTP requests
type BotanicalHandler struct {
	botanicalService *botanicalUsecase.Service
}

// NewBotanicalHandler creates a new botanical handler
func NewBotanicalHandler(botanicalService *botanicalUsecase.Service) *BotanicalHandler {
	return &BotanicalHandler{
		botanicalService: botanicalService,
	}
}

// GetAll handles GET /api/v1/botanicals
func (h *BotanicalHandler) GetAll(c *gin.Context) {
	botanicals, err := h.botanicalService.GetAllBotanicals(c.Request.Context())
	if err != nil {
		logger.Error("Failed to get botanicals", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"botanicals": botanicals,
		"count":      len(botanicals),
	})
}

// GetGinBotanicals handles GET /api/v1/gins/:gin_id/botanicals
func (h *BotanicalHandler) GetGinBotanicals(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	ginIDStr := c.Param("id")
	ginID, err := strconv.ParseInt(ginIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	botanicals, err := h.botanicalService.GetGinBotanicals(c.Request.Context(), tenantID, ginID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"botanicals": botanicals,
		"count":      len(botanicals),
	})
}

// UpdateGinBotanicals handles PUT /api/v1/gins/:gin_id/botanicals
func (h *BotanicalHandler) UpdateGinBotanicals(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	ginIDStr := c.Param("id")
	ginID, err := strconv.ParseInt(ginIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	var req struct {
		Botanicals []*models.GinBotanical `json:"botanicals" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid botanical update request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update botanicals
	if err := h.botanicalService.UpdateGinBotanicals(c.Request.Context(), tenantID, ginID, req.Botanicals); err != nil {
		logger.Error("Failed to update gin botanicals", "error", err.Error())
		response.Error(c, err)
		return
	}

	// Get updated botanicals
	botanicals, err := h.botanicalService.GetGinBotanicals(c.Request.Context(), tenantID, ginID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message":    "Botanicals updated successfully",
		"botanicals": botanicals,
	})
}
