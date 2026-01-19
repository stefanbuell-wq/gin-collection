package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/usecase/auth"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
	"github.com/yourusername/gin-collection-saas/pkg/utils"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService  *auth.Service
	cookieConfig *utils.CookieConfig
	jwtExpiry    time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.Service, cookieConfig *utils.CookieConfig, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		cookieConfig: cookieConfig,
		jwtExpiry:    jwtExpiry,
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

	// Set HttpOnly cookies for authentication
	utils.SetAuthCookies(
		c,
		authResp.Token,
		authResp.RefreshToken,
		h.cookieConfig,
		h.jwtExpiry,
		30*24*time.Hour, // Refresh token: 30 days
	)

	// Return response WITHOUT tokens (they are in HttpOnly cookies now)
	response.Created(c, gin.H{
		"user":   authResp.User,
		"tenant": authResp.Tenant,
	})
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

	var authResp *models.AuthResponse
	var err error

	// Try to get tenant from context (set by tenant middleware)
	tenantID, ok := middleware.GetTenantID(c)
	if ok {
		// Login with known tenant
		authResp, err = h.authService.Login(c.Request.Context(), &req, tenantID)
	} else {
		// Fallback: Login by email only (for localhost or no subdomain)
		logger.Debug("No tenant in context, trying login by email only")
		authResp, err = h.authService.LoginByEmail(c.Request.Context(), &req)
	}

	if err != nil {
		logger.Error("Login failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	// Set HttpOnly cookies for authentication
	utils.SetAuthCookies(
		c,
		authResp.Token,
		authResp.RefreshToken,
		h.cookieConfig,
		h.jwtExpiry,
		30*24*time.Hour, // Refresh token: 30 days
	)

	// Return response WITHOUT tokens (they are in HttpOnly cookies now)
	response.Success(c, gin.H{
		"user":   authResp.User,
		"tenant": authResp.Tenant,
	})
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Try to get refresh token from HttpOnly cookie first
	refreshToken, err := utils.GetRefreshTokenFromCookie(c)
	if err != nil || refreshToken == "" {
		// Fallback: Try JSON body (backward compatibility for API clients)
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
			response.ValidationError(c, map[string]string{
				"error": "Refresh token required",
			})
			return
		}
		refreshToken = req.RefreshToken
	}

	// Generate new access token
	newToken, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Set new access token cookie
	utils.SetAccessTokenCookie(c, newToken, h.cookieConfig, h.jwtExpiry)

	response.Success(c, gin.H{
		"message": "Token refreshed successfully",
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear HttpOnly auth cookies
	utils.ClearAuthCookies(c, h.cookieConfig)

	response.Success(c, gin.H{
		"message": "Logged out successfully",
	})
}

// GetMe handles GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ValidationError(c, map[string]string{
			"error": "User not found in context",
		})
		return
	}

	user, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get current user", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateProfile handles PUT /api/v1/auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ValidationError(c, map[string]string{
			"error": "User not found in context",
		})
		return
	}

	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid profile update request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, req.FirstName, req.LastName)
	if err != nil {
		logger.Error("Profile update failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ValidationError(c, map[string]string{
			"error": "User not found in context",
		})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid password change request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		logger.Error("Password change failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Password changed successfully",
	})
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid forgot password request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Request password reset
	if err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		logger.Error("Password reset request failed", "error", err.Error())
		// Don't reveal if the operation failed for security
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Wenn ein Konto mit dieser E-Mail existiert, wurde ein Link zum Zurücksetzen des Passworts gesendet.",
	})
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("Invalid reset password request", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Reset password
	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		logger.Error("Password reset failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"success": true,
		"message": "Passwort wurde erfolgreich zurückgesetzt. Du kannst dich jetzt anmelden.",
	})
}

// ValidateResetToken handles GET /api/v1/auth/validate-reset-token
func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.ValidationError(c, map[string]string{
			"error": "Token is required",
		})
		return
	}

	valid, err := h.authService.ValidateResetToken(c.Request.Context(), token)
	if err != nil {
		logger.Error("Token validation failed", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"valid": valid,
	})
}
