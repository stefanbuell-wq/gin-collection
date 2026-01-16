package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/usecase/tasting"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// TastingHandler handles tasting session HTTP requests
type TastingHandler struct {
	tastingService *tasting.Service
}

// NewTastingHandler creates a new tasting handler
func NewTastingHandler(tastingService *tasting.Service) *TastingHandler {
	return &TastingHandler{
		tastingService: tastingService,
	}
}

// CreateSessionRequest represents the request to create a tasting session
type CreateSessionRequest struct {
	Date       string  `json:"date"`
	Notes      *string `json:"notes"`
	Rating     *int    `json:"rating"`
	Tonic      *string `json:"tonic"`
	Botanicals *string `json:"botanicals"`
}

// UpdateSessionRequest represents the request to update a tasting session
type UpdateSessionRequest struct {
	Date       string  `json:"date"`
	Notes      *string `json:"notes"`
	Rating     *int    `json:"rating"`
	Tonic      *string `json:"tonic"`
	Botanicals *string `json:"botanicals"`
}

// GetSessions handles GET /api/v1/gins/:id/tastings
func (h *TastingHandler) GetSessions(c *gin.Context) {
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

	sessions, err := h.tastingService.GetSessionsForGin(c.Request.Context(), tenantID, ginID)
	if err != nil {
		logger.Error("Failed to get tasting sessions", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// CreateSession handles POST /api/v1/gins/:id/tastings
func (h *TastingHandler) CreateSession(c *gin.Context) {
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

	// Get user ID from context
	userID, _ := middleware.GetUserID(c)

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Parse date
	var date time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		date = parsedDate
	} else {
		date = time.Now()
	}

	session := &models.TastingSession{
		TenantID:   tenantID,
		GinID:      ginID,
		UserID:     &userID,
		Date:       date,
		Notes:      req.Notes,
		Rating:     req.Rating,
		Tonic:      req.Tonic,
		Botanicals: req.Botanicals,
	}

	if err := h.tastingService.CreateSession(c.Request.Context(), session); err != nil {
		logger.Error("Failed to create tasting session", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Created(c, session)
}

// GetSession handles GET /api/v1/gins/:id/tastings/:session_id
func (h *TastingHandler) GetSession(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := h.tastingService.GetSession(c.Request.Context(), tenantID, sessionID)
	if err != nil {
		response.Error(c, err)
		return
	}

	if session == nil {
		c.JSON(404, gin.H{"error": "Tasting session not found"})
		return
	}

	response.Success(c, session)
}

// UpdateSession handles PUT /api/v1/gins/:id/tastings/:session_id
func (h *TastingHandler) UpdateSession(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}

	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Parse date
	var date time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		date = parsedDate
	} else {
		date = time.Now()
	}

	session := &models.TastingSession{
		ID:         sessionID,
		TenantID:   tenantID,
		Date:       date,
		Notes:      req.Notes,
		Rating:     req.Rating,
		Tonic:      req.Tonic,
		Botanicals: req.Botanicals,
	}

	if err := h.tastingService.UpdateSession(c.Request.Context(), session); err != nil {
		logger.Error("Failed to update tasting session", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, session)
}

// DeleteSession handles DELETE /api/v1/gins/:id/tastings/:session_id
func (h *TastingHandler) DeleteSession(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := h.tastingService.DeleteSession(c.Request.Context(), tenantID, sessionID); err != nil {
		logger.Error("Failed to delete tasting session", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Tasting session deleted successfully",
	})
}

// GetRecentSessions handles GET /api/v1/tastings/recent
func (h *TastingHandler) GetRecentSessions(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	sessions, err := h.tastingService.GetRecentSessions(c.Request.Context(), tenantID, limit)
	if err != nil {
		logger.Error("Failed to get recent tasting sessions", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}
