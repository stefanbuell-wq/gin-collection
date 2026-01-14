package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	userUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/user"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// UserHandler handles user management HTTP requests (Enterprise only)
type UserHandler struct {
	userService *userUsecase.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *userUsecase.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// List handles GET /api/v1/users
func (h *UserHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	users, err := h.userService.ListUsers(c.Request.Context(), tenantID, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"users": users,
		"count": len(users),
	})
}

// Invite handles POST /api/v1/users/invite
func (h *UserHandler) Invite(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Email     string          `json:"email" binding:"required,email"`
		FirstName string          `json:"first_name"`
		LastName  string          `json:"last_name"`
		Role      models.UserRole `json:"role" binding:"required,oneof=admin member viewer"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Prevent inviting as owner (only one owner per tenant)
	if req.Role == models.RoleOwner {
		response.ValidationError(c, map[string]string{
			"error": "Cannot invite user as owner - each tenant has exactly one owner",
		})
		return
	}

	user, err := h.userService.InviteUser(
		c.Request.Context(),
		tenantID,
		userID,
		req.Email,
		req.FirstName,
		req.LastName,
		req.Role,
	)

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, user)
}

// Update handles PUT /api/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	requesterUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Email     string          `json:"email" binding:"required,email"`
		FirstName *string         `json:"first_name"`
		LastName  *string         `json:"last_name"`
		Role      models.UserRole `json:"role" binding:"required,oneof=owner admin member viewer"`
		IsActive  bool            `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(
		c.Request.Context(),
		tenantID,
		requesterUserID,
		targetUserID,
		req.Email,
		req.FirstName,
		req.LastName,
		req.Role,
		req.IsActive,
	)

	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	requesterUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), tenantID, requesterUserID, targetUserID); err != nil {
		logger.Error("Failed to delete user", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "User deleted successfully",
	})
}

// GenerateAPIKey handles POST /api/v1/users/:id/api-key
func (h *UserHandler) GenerateAPIKey(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	requesterUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	apiKey, err := h.userService.GenerateAPIKey(c.Request.Context(), tenantID, requesterUserID, targetUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"api_key": apiKey,
		"message": "API key generated successfully - store this securely, it will not be shown again",
	})
}

// RevokeAPIKey handles DELETE /api/v1/users/:id/api-key
func (h *UserHandler) RevokeAPIKey(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	requesterUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	targetUserIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.RevokeAPIKey(c.Request.Context(), tenantID, requesterUserID, targetUserID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "API key revoked successfully",
	})
}
