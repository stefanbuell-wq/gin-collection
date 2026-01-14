package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/usecase/auth"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid registration request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Register new tenant and user
	authResp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Registration failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Created(c, authResp)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid login request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Get tenant from context (set by tenant middleware)
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		logger.Error("Tenant ID not found in context")
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	// Login user
	authResp, err := h.authService.Login(c.Request.Context(), &req, tenantID)
	if err != nil {
		logger.Error("Login failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, authResp)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Generate new access token
	newToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// For JWT, logout is handled client-side by removing the token
	// Server-side logout would require a token blacklist (Redis)
	// For now, just return success
	response.Success(c, gin.H{
		"message": "Logged out successfully",
	})
}
