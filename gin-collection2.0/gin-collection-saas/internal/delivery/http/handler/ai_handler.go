package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/external"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiClient *external.AIClient
}

// NewAIHandler creates a new AI handler
func NewAIHandler(aiClient *external.AIClient) *AIHandler {
	return &AIHandler{
		aiClient: aiClient,
	}
}

// SuggestGinInfoRequest represents the request body for gin suggestion
type SuggestGinInfoRequest struct {
	Name  string `json:"name" binding:"required"`
	Brand string `json:"brand"`
}

// SuggestGinInfo handles POST /api/v1/ai/suggest-gin
func (h *AIHandler) SuggestGinInfo(c *gin.Context) {
	var req SuggestGinInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Name ist erforderlich")
		return
	}

	if !h.aiClient.IsEnabled() {
		response.ServiceUnavailable(c, "AI-Service ist nicht verfügbar")
		return
	}

	suggestion, err := h.aiClient.SuggestGinInfo(req.Name, req.Brand)
	if err != nil {
		logger.Error("AI suggestion failed", "error", err.Error(), "name", req.Name, "brand", req.Brand)
		response.InternalError(c, "AI-Vorschlag konnte nicht generiert werden")
		return
	}

	response.Success(c, gin.H{
		"suggestion": suggestion,
	})
}

// Status handles GET /api/v1/ai/status
func (h *AIHandler) Status(c *gin.Context) {
	response.Success(c, gin.H{
		"enabled": h.aiClient.IsEnabled(),
	})
}
