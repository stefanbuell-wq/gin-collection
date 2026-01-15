package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	cocktailUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/cocktail"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// CocktailHandler handles cocktail HTTP requests
type CocktailHandler struct {
	cocktailService *cocktailUsecase.Service
}

// NewCocktailHandler creates a new cocktail handler
func NewCocktailHandler(cocktailService *cocktailUsecase.Service) *CocktailHandler {
	return &CocktailHandler{
		cocktailService: cocktailService,
	}
}

// GetAll handles GET /api/v1/cocktails
func (h *CocktailHandler) GetAll(c *gin.Context) {
	cocktails, err := h.cocktailService.GetAllCocktails(c.Request.Context())
	if err != nil {
		logger.Error("Failed to get cocktails", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"cocktails": cocktails,
		"count":     len(cocktails),
	})
}

// GetByID handles GET /api/v1/cocktails/:id
func (h *CocktailHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid cocktail ID"})
		return
	}

	cocktail, err := h.cocktailService.GetCocktailByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, cocktail)
}

// GetGinCocktails handles GET /api/v1/gins/:gin_id/cocktails
func (h *CocktailHandler) GetGinCocktails(c *gin.Context) {
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

	cocktails, err := h.cocktailService.GetCocktailsForGin(c.Request.Context(), tenantID, ginID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"cocktails": cocktails,
		"count":     len(cocktails),
	})
}

// LinkCocktail handles POST /api/v1/gins/:gin_id/cocktails/:cocktail_id
func (h *CocktailHandler) LinkCocktail(c *gin.Context) {
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

	cocktailIDStr := c.Param("cocktail_id")
	cocktailID, err := strconv.ParseInt(cocktailIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid cocktail ID"})
		return
	}

	// Link cocktail to gin
	if err := h.cocktailService.LinkCocktailToGin(c.Request.Context(), tenantID, ginID, cocktailID); err != nil {
		logger.Error("Failed to link cocktail to gin", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Cocktail linked successfully",
	})
}

// UnlinkCocktail handles DELETE /api/v1/gins/:gin_id/cocktails/:cocktail_id
func (h *CocktailHandler) UnlinkCocktail(c *gin.Context) {
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

	cocktailIDStr := c.Param("cocktail_id")
	cocktailID, err := strconv.ParseInt(cocktailIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid cocktail ID"})
		return
	}

	// Unlink cocktail from gin
	if err := h.cocktailService.UnlinkCocktailFromGin(c.Request.Context(), tenantID, ginID, cocktailID); err != nil {
		logger.Error("Failed to unlink cocktail from gin", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Cocktail unlinked successfully",
	})
}
