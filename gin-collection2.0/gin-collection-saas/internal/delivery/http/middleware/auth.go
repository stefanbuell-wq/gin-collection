package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
	"github.com/yourusername/gin-collection-saas/pkg/utils"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtSecret      string
	userRepo       repositories.UserRepository
	tokenBlacklist *utils.TokenBlacklist
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtSecret string, userRepo repositories.UserRepository, blacklist *utils.TokenBlacklist) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:      jwtSecret,
		userRepo:       userRepo,
		tokenBlacklist: blacklist,
	}
}

// RequireAuth middleware validates JWT token from cookie OR header
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var tokenSource string

		// Priority 1: Authorization header (for API keys and backward compatibility)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
				tokenSource = "header"
			}
		}

		// Priority 2: HttpOnly cookie
		if tokenString == "" {
			cookie, err := c.Cookie(utils.AccessTokenCookieName)
			if err == nil && cookie != "" {
				tokenString = cookie
				tokenSource = "cookie"
			}
		}

		// No token found
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(tokenString, am.jwtSecret)
		if err != nil {
			logger.Debug("Invalid JWT token", "error", err.Error(), "source", tokenSource)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if token is blacklisted
		if am.tokenBlacklist != nil {
			// Check specific token (by JTI)
			if claims.ID != "" && am.tokenBlacklist.IsRevoked(c.Request.Context(), claims.ID) {
				logger.Debug("Revoked token used", "jti", claims.ID, "user_id", claims.UserID)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token has been revoked",
				})
				c.Abort()
				return
			}

			// Check if all user tokens before certain time are revoked (password change)
			if claims.IssuedAt != nil && am.tokenBlacklist.IsUserTokenRevoked(c.Request.Context(), claims.UserID, claims.IssuedAt.Time) {
				logger.Debug("User tokens revoked", "user_id", claims.UserID, "issued_at", claims.IssuedAt.Time)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token has been revoked",
				})
				c.Abort()
				return
			}
		}

		// Store claims in context
		c.Set("jwt_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_source", tokenSource)

		logger.Debug("User authenticated", "user_id", claims.UserID, "email", claims.Email, "source", tokenSource)

		c.Next()
	}
}

// OptionalAuth middleware validates JWT token if present but doesn't require it
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try header first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// Try cookie if no header token
		if tokenString == "" {
			cookie, _ := c.Cookie(utils.AccessTokenCookieName)
			tokenString = cookie
		}

		if tokenString == "" {
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(tokenString, am.jwtSecret)
		if err != nil {
			c.Next()
			return
		}

		// Check if token is blacklisted (silently skip if revoked)
		if am.tokenBlacklist != nil {
			if claims.ID != "" && am.tokenBlacklist.IsRevoked(c.Request.Context(), claims.ID) {
				c.Next()
				return
			}
			if claims.IssuedAt != nil && am.tokenBlacklist.IsUserTokenRevoked(c.Request.Context(), claims.UserID, claims.IssuedAt.Time) {
				c.Next()
				return
			}
		}

		c.Set("jwt_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found in context",
			})
			c.Abort()
			return
		}

		role := models.UserRole(userRole.(string))

		// Check if user has one of the allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		c.Abort()
	}
}

// RequirePermission middleware checks if user has specific permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found in context",
			})
			c.Abort()
			return
		}

		role := models.UserRole(userRole.(string))
		if !role.HasPermission(permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions for action: " + permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyAuth middleware authenticates using API key (Enterprise only)
func (am *AuthMiddleware) APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required. Use X-API-Key header.",
			})
			c.Abort()
			return
		}

		// Retrieve user by API key
		user, err := am.userRepo.GetByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			logger.Debug("Invalid API key", "error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", user.ID)
		c.Set("tenant_id", user.TenantID)
		c.Set("user_email", user.Email)
		c.Set("user_role", string(user.Role))

		logger.Debug("API key authenticated", "user_id", user.ID, "email", user.Email)

		c.Next()
	}
}

// GetUserID helper to retrieve user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUserRole helper to retrieve user role from context
func GetUserRole(c *gin.Context) (models.UserRole, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := userRole.(string)
	return models.UserRole(role), ok
}
