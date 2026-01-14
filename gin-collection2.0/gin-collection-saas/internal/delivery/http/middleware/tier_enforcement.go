package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// TierEnforcementMiddleware enforces subscription tier limits
type TierEnforcementMiddleware struct {
	usageRepo repositories.UsageMetricsRepository
}

// NewTierEnforcementMiddleware creates a new tier enforcement middleware
func NewTierEnforcementMiddleware(usageRepo repositories.UsageMetricsRepository) *TierEnforcementMiddleware {
	return &TierEnforcementMiddleware{
		usageRepo: usageRepo,
	}
}

// RequireFeature checks if a feature is available in the current tier
func (tem *TierEnforcementMiddleware) RequireFeature(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, ok := GetTenant(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant not found in context",
			})
			c.Abort()
			return
		}

		limits := tenant.GetLimits()

		var hasFeature bool
		var featureName string

		switch feature {
		case "botanicals":
			hasFeature = limits.HasBotanicals
			featureName = "Botanicals tracking"
		case "cocktails":
			hasFeature = limits.HasCocktails
			featureName = "Cocktail recipes"
		case "ai_suggestions":
			hasFeature = limits.HasAISuggestions
			featureName = "AI suggestions"
		case "export":
			hasFeature = limits.HasExport
			featureName = "Data export"
		case "import":
			hasFeature = limits.HasImport
			featureName = "Data import"
		case "multi_user":
			hasFeature = limits.HasMultiUser
			featureName = "Multi-user support"
		case "api_access":
			hasFeature = limits.HasAPIAccess
			featureName = "API access"
		default:
			hasFeature = true // Unknown features are allowed by default
		}

		if !hasFeature {
			logger.Debug("Feature not available", "tenant_id", tenant.ID, "tier", tenant.Tier, "feature", feature)
			c.JSON(http.StatusForbidden, gin.H{
				"error":            featureName + " requires " + getRequiredTier(feature) + " tier",
				"upgrade_required": true,
				"current_tier":     tenant.Tier,
				"feature":          feature,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckGinLimit checks if tenant can create more gins
func (tem *TierEnforcementMiddleware) CheckGinLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only check on POST requests
		if c.Request.Method != "POST" {
			c.Next()
			return
		}

		tenant, ok := GetTenant(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant not found in context",
			})
			c.Abort()
			return
		}

		limits := tenant.GetLimits()

		// If unlimited, allow
		if limits.MaxGins == nil {
			c.Next()
			return
		}

		// Get current gin count
		currentCount, err := tem.usageRepo.GetMetric(c.Request.Context(), tenant.ID, "gin_count")
		if err != nil {
			logger.Error("Failed to get gin count", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check gin limit",
			})
			c.Abort()
			return
		}

		// Check if limit reached
		if currentCount >= *limits.MaxGins {
			logger.Debug("Gin limit reached", "tenant_id", tenant.ID, "current", currentCount, "limit", *limits.MaxGins)
			c.JSON(http.StatusForbidden, gin.H{
				"error":            "Gin limit reached. Please upgrade to add more gins.",
				"upgrade_required": true,
				"current_tier":     tenant.Tier,
				"limit":            *limits.MaxGins,
				"current_count":    currentCount,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckPhotoLimit checks if tenant can upload more photos for a gin
func (tem *TierEnforcementMiddleware) CheckPhotoLimit(currentPhotoCount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, ok := GetTenant(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant not found in context",
			})
			c.Abort()
			return
		}

		limits := tenant.GetLimits()

		if currentPhotoCount >= limits.MaxPhotosPerGin {
			logger.Debug("Photo limit reached", "tenant_id", tenant.ID, "current", currentPhotoCount, "limit", limits.MaxPhotosPerGin)
			c.JSON(http.StatusForbidden, gin.H{
				"error":            "Photo limit reached for this gin. Please upgrade to add more photos.",
				"upgrade_required": true,
				"current_tier":     tenant.Tier,
				"limit":            limits.MaxPhotosPerGin,
				"current_count":    currentPhotoCount,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckStorageLimit checks if tenant has storage available
func (tem *TierEnforcementMiddleware) CheckStorageLimit(fileSizeKB int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, ok := GetTenant(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant not found in context",
			})
			c.Abort()
			return
		}

		limits := tenant.GetLimits()

		// If unlimited storage, allow
		if limits.StorageLimitMB == nil {
			c.Next()
			return
		}

		// Get current storage usage
		currentStorageMB, err := tem.usageRepo.GetMetric(c.Request.Context(), tenant.ID, "storage_mb")
		if err != nil {
			logger.Error("Failed to get storage usage", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check storage limit",
			})
			c.Abort()
			return
		}

		// Convert KB to MB
		fileSizeMB := fileSizeKB / 1024

		// Check if limit would be exceeded
		if currentStorageMB+fileSizeMB > *limits.StorageLimitMB {
			logger.Debug("Storage limit would be exceeded", "tenant_id", tenant.ID, "current", currentStorageMB, "limit", *limits.StorageLimitMB)
			c.JSON(http.StatusForbidden, gin.H{
				"error":            "Storage limit would be exceeded. Please upgrade for more storage.",
				"upgrade_required": true,
				"current_tier":     tenant.Tier,
				"limit_mb":         *limits.StorageLimitMB,
				"current_mb":       currentStorageMB,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getRequiredTier returns the minimum tier required for a feature
func getRequiredTier(feature string) string {
	switch feature {
	case "botanicals", "cocktails", "ai_suggestions", "export", "import":
		return "Pro"
	case "multi_user", "api_access":
		return "Enterprise"
	default:
		return "Pro"
	}
}

// UpgradeResponse is the response structure for upgrade suggestions
type UpgradeResponse struct {
	Error           string                   `json:"error"`
	UpgradeRequired bool                     `json:"upgrade_required"`
	CurrentTier     models.SubscriptionTier  `json:"current_tier"`
	RequiredTier    string                   `json:"required_tier"`
	Feature         string                   `json:"feature"`
}
