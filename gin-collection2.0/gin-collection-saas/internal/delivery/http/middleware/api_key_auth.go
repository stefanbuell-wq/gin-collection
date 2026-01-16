package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// APIKeyAuthMiddleware handles API key authentication for Enterprise tier
type APIKeyAuthMiddleware struct {
	userRepo   repositories.UserRepository
	tenantRepo repositories.TenantRepository
}

// NewAPIKeyAuthMiddleware creates a new API key authentication middleware
func NewAPIKeyAuthMiddleware(userRepo repositories.UserRepository, tenantRepo repositories.TenantRepository) *APIKeyAuthMiddleware {
	return &APIKeyAuthMiddleware{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
	}
}

// Authenticate validates the API key and sets user context
func (m *APIKeyAuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from Authorization header: "Bearer sk_..."
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "API key required - provide Authorization: Bearer sk_xxx"})
			c.Abort()
			return
		}

		// Extract the API key
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid Authorization header format - use: Bearer sk_xxx"})
			c.Abort()
			return
		}

		apiKey := parts[1]

		// Validate API key prefix
		if !strings.HasPrefix(apiKey, "sk_") {
			c.JSON(401, gin.H{"error": "Invalid API key format"})
			c.Abort()
			return
		}

		// Lookup user by API key
		user, err := m.userRepo.GetByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			logger.Debug("Invalid API key attempt", "api_key_prefix", apiKey[:10])
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			logger.Warn("Inactive user attempted API access", "user_id", user.ID, "tenant_id", user.TenantID)
			c.JSON(403, gin.H{"error": "User account is inactive"})
			c.Abort()
			return
		}

		// Get tenant to verify Enterprise tier
		tenant, err := m.tenantRepo.GetByID(c.Request.Context(), user.TenantID)
		if err != nil {
			logger.Error("Failed to get tenant for API key auth", "tenant_id", user.TenantID, "error", err.Error())
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Verify tenant has Pro or Enterprise tier (API access)
		if tenant.Tier != "enterprise" && tenant.Tier != "pro" {
			logger.Warn("Non-Pro/Enterprise tenant attempted API access", "tenant_id", user.TenantID, "tier", tenant.Tier)
			c.JSON(403, gin.H{
				"error":            "API access requires Pro or Enterprise subscription",
				"upgrade_required": true,
			})
			c.Abort()
			return
		}

		// Verify tenant is active
		if tenant.Status != "active" {
			logger.Warn("Inactive tenant attempted API access", "tenant_id", user.TenantID, "status", tenant.Status)
			c.JSON(403, gin.H{"error": "Tenant account is not active"})
			c.Abort()
			return
		}

		// Set user and tenant in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("tenant", tenant)
		c.Set("tenant_id", user.TenantID)

		logger.Debug("API key authenticated successfully", "user_id", user.ID, "tenant_id", user.TenantID)

		c.Next()
	}
}

// OptionalAPIKey allows both JWT and API key authentication
// This middleware tries API key first, then falls back to JWT
func (m *APIKeyAuthMiddleware) OptionalAPIKey(authMiddleware *AuthMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// If no Authorization header, skip API key auth
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if it's an API key (starts with "Bearer sk_")
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" && strings.HasPrefix(parts[1], "sk_") {
			// Use API key authentication
			m.Authenticate()(c)
			return
		}

		// Otherwise, use JWT authentication
		authMiddleware.RequireAuth()(c)
	}
}
