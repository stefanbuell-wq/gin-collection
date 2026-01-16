package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
)

// GinReferenceHandler handles gin reference catalog requests
type GinReferenceHandler struct {
	repo repositories.GinReferenceRepository
}

// NewGinReferenceHandler creates a new gin reference handler
func NewGinReferenceHandler(repo repositories.GinReferenceRepository) *GinReferenceHandler {
	return &GinReferenceHandler{repo: repo}
}

// Search searches gin references
// GET /api/v1/gin-references?q=hendricks&country=Scotland&type=New%20Western&limit=20&offset=0
func (h *GinReferenceHandler) Search(c *gin.Context) {
	params := &models.GinReferenceSearchParams{
		Query:   c.Query("q"),
		Country: c.Query("country"),
		GinType: c.Query("type"),
	}

	// Parse pagination
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		params.Limit = limit
	} else {
		params.Limit = 20
	}

	if offset, err := strconv.Atoi(c.Query("offset")); err == nil {
		params.Offset = offset
	}

	gins, total, err := h.repo.Search(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search gin references",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"gins":   gins,
			"total":  total,
			"limit":  params.Limit,
			"offset": params.Offset,
		},
		"success": true,
	})
}

// GetByID retrieves a single gin reference
// GET /api/v1/gin-references/:id
func (h *GinReferenceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"success": false,
		})
		return
	}

	gin_ref, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get gin reference",
			"success": false,
		})
		return
	}

	if gin_ref == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Gin reference not found",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    gin_ref,
		"success": true,
	})
}

// GetFilters returns available filter options (countries, types, brands)
// GET /api/v1/gin-references/filters
func (h *GinReferenceHandler) GetFilters(c *gin.Context) {
	countries, err := h.repo.GetCountries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get countries",
			"success": false,
		})
		return
	}

	ginTypes, err := h.repo.GetGinTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get gin types",
			"success": false,
		})
		return
	}

	brands, err := h.repo.GetBrands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get brands",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"countries": countries,
			"gin_types": ginTypes,
			"brands":    brands,
		},
		"success": true,
	})
}
