package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	ginUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/gin"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// GinHandler handles gin HTTP requests
type GinHandler struct {
	ginService *ginUsecase.Service
}

// NewGinHandler creates a new gin handler
func NewGinHandler(ginService *ginUsecase.Service) *GinHandler {
	return &GinHandler{
		ginService: ginService,
	}
}

// List handles GET /api/v1/gins
func (h *GinHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	// Parse query parameters
	filter := &models.GinFilter{
		TenantID:  tenantID,
		SortBy:    c.DefaultQuery("sort", "created_at"),
		SortOrder: c.DefaultQuery("order", "desc"),
	}

	// Filter by status
	if filterParam := c.Query("filter"); filterParam != "" {
		if filterParam == "available" {
			isFinished := false
			filter.IsFinished = &isFinished
		} else if filterParam == "finished" {
			isFinished := true
			filter.IsFinished = &isFinished
		}
	}

	// Filter by type
	if ginType := c.Query("type"); ginType != "" {
		filter.GinType = &ginType
	}

	// Filter by country
	if country := c.Query("country"); country != "" {
		filter.Country = &country
	}

	// Pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	} else {
		filter.Limit = 50 // Default
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	// Get gins
	gins, err := h.ginService.List(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list gins", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"gins":  gins,
		"count": len(gins),
	})
}

// Create handles POST /api/v1/gins
func (h *GinHandler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	var ginModel models.Gin
	if err := c.ShouldBindJSON(&ginModel); err != nil {
		logger.Debug("Invalid gin creation request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Set tenant ID
	ginModel.TenantID = tenantID

	// Create gin
	if err := h.ginService.Create(c.Request.Context(), &ginModel); err != nil {
		logger.Error("Failed to create gin", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Created(c, ginModel)
}

// Get handles GET /api/v1/gins/:id
func (h *GinHandler) Get(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	ginModel, err := h.ginService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, ginModel)
}

// Update handles PUT /api/v1/gins/:id
func (h *GinHandler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	var ginModel models.Gin
	if err := c.ShouldBindJSON(&ginModel); err != nil {
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Set IDs
	ginModel.ID = id
	ginModel.TenantID = tenantID

	// Update gin
	if err := h.ginService.Update(c.Request.Context(), &ginModel); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, ginModel)
}

// Delete handles DELETE /api/v1/gins/:id
func (h *GinHandler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	// Delete gin
	if err := h.ginService.Delete(c.Request.Context(), tenantID, id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Gin deleted successfully",
	})
}

// Search handles GET /api/v1/gins/search
func (h *GinHandler) Search(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "Search query required"})
		return
	}

	// Pagination
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Search
	gins, err := h.ginService.Search(c.Request.Context(), tenantID, query, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"gins":  gins,
		"count": len(gins),
		"query": query,
	})
}

// Stats handles GET /api/v1/gins/stats
func (h *GinHandler) Stats(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	stats, err := h.ginService.GetStats(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stats)
}

// Export handles POST /api/v1/gins/export
func (h *GinHandler) Export(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	format := c.DefaultQuery("format", "json")

	if format == "json" {
		data, err := h.ginService.ExportJSON(c.Request.Context(), tenantID)
		if err != nil {
			response.Error(c, err)
			return
		}

		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=gins.json")
		c.Data(200, "application/json", data)
	} else if format == "csv" {
		data, err := h.ginService.ExportCSV(c.Request.Context(), tenantID)
		if err != nil {
			response.Error(c, err)
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=gins.csv")
		c.String(200, data)
	} else {
		c.JSON(400, gin.H{"error": "Invalid format. Use 'json' or 'csv'"})
	}
}

// Import handles POST /api/v1/gins/import
func (h *GinHandler) Import(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	// Read JSON body
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	// Import
	imported, err := h.ginService.ImportJSON(c.Request.Context(), tenantID, data)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message":  "Import completed",
		"imported": imported,
	})
}

// Suggestions handles GET /api/v1/gins/:id/suggestions
func (h *GinHandler) Suggestions(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	// Get limit from query
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get similar gins
	similarGins, err := h.ginService.GetSimilarGins(c.Request.Context(), tenantID, id, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"suggestions": similarGins,
		"count":       len(similarGins),
	})
}
