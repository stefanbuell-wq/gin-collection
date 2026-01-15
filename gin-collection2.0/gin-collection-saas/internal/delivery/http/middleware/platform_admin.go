package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	adminUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/admin"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// PlatformAdminMiddleware handles platform admin authentication
type PlatformAdminMiddleware struct {
	adminService *adminUsecase.Service
}

// NewPlatformAdminMiddleware creates a new platform admin middleware
func NewPlatformAdminMiddleware(adminService *adminUsecase.Service) *PlatformAdminMiddleware {
	return &PlatformAdminMiddleware{
		adminService: adminService,
	}
}

// RequirePlatformAdmin validates the request has a valid platform admin token
func (m *PlatformAdminMiddleware) RequirePlatformAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Debug("No authorization header")
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Debug("Invalid authorization format")
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.adminService.ValidateAdminToken(tokenString)
		if err != nil {
			logger.Debug("Invalid admin token", "error", err.Error())
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Verify it's a platform admin token
		if !claims.IsPlatformAdmin {
			logger.Warn("Token is not a platform admin token")
			c.JSON(403, gin.H{"error": "Platform admin access required"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("admin_claims", claims)
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_email", claims.Email)
		c.Set("is_platform_admin", true)

		c.Next()
	}
}

// GetAdminID retrieves the admin ID from the context
func GetAdminID(c *gin.Context) (int64, bool) {
	id, exists := c.Get("admin_id")
	if !exists {
		return 0, false
	}
	adminID, ok := id.(int64)
	return adminID, ok
}

// GetAdminEmail retrieves the admin email from the context
func GetAdminEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("admin_email")
	if !exists {
		return "", false
	}
	return email.(string), true
}

// IsPlatformAdmin checks if the current request is from a platform admin
func IsPlatformAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("is_platform_admin")
	if !exists {
		return false
	}
	return isAdmin.(bool)
}
