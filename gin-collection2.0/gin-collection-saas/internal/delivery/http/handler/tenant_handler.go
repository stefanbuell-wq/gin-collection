package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// TenantHandler handles tenant management HTTP requests
type TenantHandler struct {
	tenantRepo       repositories.TenantRepository
	usageMetricsRepo repositories.UsageMetricsRepository
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantRepo repositories.TenantRepository, usageMetricsRepo repositories.UsageMetricsRepository) *TenantHandler {
	return &TenantHandler{
		tenantRepo:       tenantRepo,
		usageMetricsRepo: usageMetricsRepo,
	}
}

// GetCurrent handles GET /api/v1/tenants/current
func (h *TenantHandler) GetCurrent(c *gin.Context) {
	tenant, ok := middleware.GetTenant(c)
	if !ok {
		response.ValidationError(c, map[string]string{
			"error": "Tenant not found in context",
		})
		return
	}

	// Get the plan limits for this tenant
	limits := tenant.GetLimits()

	response.Success(c, gin.H{
		"tenant": tenant,
		"limits": limits,
	})
}

// UpdateCurrent handles PUT /api/v1/tenants/current
func (h *TenantHandler) UpdateCurrent(c *gin.Context) {
	tenant, ok := middleware.GetTenant(c)
	if !ok {
		response.ValidationError(c, map[string]string{
			"error": "Tenant not found in context",
		})
		return
	}

	// Only owner can update tenant settings
	userRole, _ := c.Get("user_role")
	if userRole != "owner" {
		response.ValidationError(c, map[string]string{
			"error": "Only the tenant owner can update tenant settings",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,min=2,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid tenant update request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Update tenant name
	tenant.Name = req.Name

	if err := h.tenantRepo.Update(c.Request.Context(), tenant); err != nil {
		logger.Error("Failed to update tenant", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"tenant":  tenant,
		"message": "Tenant updated successfully",
	})
}

// GetUsage handles GET /api/v1/tenants/usage
func (h *TenantHandler) GetUsage(c *gin.Context) {
	tenant, ok := middleware.GetTenant(c)
	if !ok {
		response.ValidationError(c, map[string]string{
			"error": "Tenant not found in context",
		})
		return
	}

	// Get usage metrics
	ginCount, _ := h.usageMetricsRepo.GetMetric(c.Request.Context(), tenant.ID, "gin_count")
	storageUsedMB, _ := h.usageMetricsRepo.GetMetric(c.Request.Context(), tenant.ID, "storage_mb")
	photoCount, _ := h.usageMetricsRepo.GetMetric(c.Request.Context(), tenant.ID, "photo_count")

	// Get plan limits
	limits := tenant.GetLimits()

	// Calculate percentages
	var ginPercentage float64
	var storagePercentage float64

	if limits.MaxGins != nil && *limits.MaxGins > 0 {
		ginPercentage = float64(ginCount) / float64(*limits.MaxGins) * 100
	}

	if limits.StorageLimitMB != nil && *limits.StorageLimitMB > 0 {
		storagePercentage = float64(storageUsedMB) / float64(*limits.StorageLimitMB) * 100
	}

	response.Success(c, gin.H{
		"usage": gin.H{
			"gins": gin.H{
				"current":    ginCount,
				"limit":      limits.MaxGins,
				"percentage": ginPercentage,
				"unlimited":  limits.MaxGins == nil,
			},
			"storage": gin.H{
				"current_mb":  storageUsedMB,
				"limit_mb":    limits.StorageLimitMB,
				"percentage":  storagePercentage,
				"unlimited":   limits.StorageLimitMB == nil,
			},
			"photos": gin.H{
				"total":            photoCount,
				"per_gin_limit":    limits.MaxPhotosPerGin,
			},
		},
		"tier": tenant.Tier,
		"features": gin.H{
			"botanicals":     limits.HasBotanicals,
			"cocktails":      limits.HasCocktails,
			"ai_suggestions": limits.HasAISuggestions,
			"export":         limits.HasExport,
			"import":         limits.HasImport,
			"multi_user":     limits.HasMultiUser,
			"api_access":     limits.HasAPIAccess,
		},
	})
}
