package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// TenantMiddleware extracts tenant from subdomain or JWT and validates it
type TenantMiddleware struct {
	tenantRepo repositories.TenantRepository
}

// NewTenantMiddleware creates a new tenant middleware
func NewTenantMiddleware(tenantRepo repositories.TenantRepository) *TenantMiddleware {
	return &TenantMiddleware{
		tenantRepo: tenantRepo,
	}
}

// ExtractTenant middleware extracts tenant from subdomain or JWT claim
func (tm *TenantMiddleware) ExtractTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenant *models.Tenant
		var err error

		// 1. Try to extract from subdomain
		host := c.Request.Host
		subdomain := extractSubdomain(host)

		if subdomain != "" && subdomain != "www" && subdomain != "api" {
			tenant, err = tm.tenantRepo.GetBySubdomain(c.Request.Context(), subdomain)
			if err != nil {
				logger.Debug("Tenant not found by subdomain", "subdomain", subdomain, "error", err.Error())
			}
		}

		// 2. Fallback to JWT claim (for API access without subdomain)
		if tenant == nil {
			// Check if JWT claims are set (from auth middleware)
			if claims, exists := c.Get("jwt_claims"); exists {
				if jwtClaims, ok := claims.(map[string]interface{}); ok {
					if tenantID, ok := jwtClaims["tenant_id"].(float64); ok {
						tenant, err = tm.tenantRepo.GetByID(c.Request.Context(), int64(tenantID))
						if err != nil {
							logger.Debug("Tenant not found by ID", "tenant_id", tenantID, "error", err.Error())
						}
					}
				}
			}
		}

		// 3. If still no tenant found, reject request
		if tenant == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tenant not found",
			})
			c.Abort()
			return
		}

		// 4. Validate tenant status
		if tenant.Status != models.TenantStatusActive {
			statusMsg := "Tenant account is suspended"
			if tenant.Status == models.TenantStatusCancelled {
				statusMsg = "Tenant account is cancelled"
			}
			c.JSON(http.StatusForbidden, gin.H{
				"error": statusMsg,
			})
			c.Abort()
			return
		}

		// 5. Store tenant in context
		ctx := context.WithValue(c.Request.Context(), "tenant", tenant)
		c.Request = c.Request.WithContext(ctx)
		c.Set("tenant", tenant)
		c.Set("tenant_id", tenant.ID)

		logger.Debug("Tenant extracted", "tenant_id", tenant.ID, "subdomain", tenant.Subdomain)

		c.Next()
	}
}

// extractSubdomain extracts the subdomain from a host
// Example: "customer1.ginapp.com" -> "customer1"
func extractSubdomain(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Split by dots
	parts := strings.Split(host, ".")

	// If localhost or IP, no subdomain
	if host == "localhost" || strings.Contains(host, "127.0.0.1") {
		return ""
	}

	// Need at least 3 parts for subdomain (e.g., sub.domain.com)
	if len(parts) < 3 {
		return ""
	}

	// Return first part as subdomain
	return parts[0]
}

// GetTenant helper to retrieve tenant from context
func GetTenant(c *gin.Context) (*models.Tenant, bool) {
	tenant, exists := c.Get("tenant")
	if !exists {
		return nil, false
	}
	t, ok := tenant.(*models.Tenant)
	return t, ok
}

// GetTenantID helper to retrieve tenant ID from context
func GetTenantID(c *gin.Context) (int64, bool) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return 0, false
	}
	id, ok := tenantID.(int64)
	return id, ok
}
