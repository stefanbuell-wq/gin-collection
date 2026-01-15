package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	adminUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/admin"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Handler handles all platform admin HTTP requests
type Handler struct {
	adminService *adminUsecase.Service
}

// NewHandler creates a new admin handler
func NewHandler(adminService *adminUsecase.Service) *Handler {
	return &Handler{
		adminService: adminService,
	}
}

// ==================== AUTH ====================

// LoginRequest represents admin login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login handles POST /admin/api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	response, err := h.adminService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		logger.Warn("Admin login failed", "email", req.Email, "error", err.Error())
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, response)
}

// Me handles GET /admin/api/v1/auth/me
func (h *Handler) Me(c *gin.Context) {
	adminID, ok := middleware.GetAdminID(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	admin, err := h.adminService.GetAdmin(c.Request.Context(), adminID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get admin"})
		return
	}

	c.JSON(200, admin)
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword handles POST /admin/api/v1/auth/change-password
func (h *Handler) ChangePassword(c *gin.Context) {
	adminID, ok := middleware.GetAdminID(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := h.adminService.ChangePassword(c.Request.Context(), adminID, req.OldPassword, req.NewPassword); err != nil {
		logger.Error("Password change failed", "error", err.Error())
		c.JSON(400, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

// ==================== STATS ====================

// GetStats handles GET /admin/api/v1/stats/overview
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetPlatformStats(c.Request.Context())
	if err != nil {
		logger.Error("Failed to get platform stats", "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(200, stats)
}

// ==================== TENANTS ====================

// ListTenants handles GET /admin/api/v1/tenants
func (h *Handler) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	tenants, total, err := h.adminService.GetAllTenants(c.Request.Context(), page, limit)
	if err != nil {
		logger.Error("Failed to list tenants", "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to list tenants"})
		return
	}

	c.JSON(200, gin.H{
		"tenants": tenants,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// GetTenant handles GET /admin/api/v1/tenants/:id
func (h *Handler) GetTenant(c *gin.Context) {
	// This would need the tenant repository injected
	// For now, we'll return the tenant from the list
	c.JSON(501, gin.H{"error": "Not implemented yet"})
}

// SuspendTenant handles POST /admin/api/v1/tenants/:id/suspend
func (h *Handler) SuspendTenant(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tenant ID"})
		return
	}

	if err := h.adminService.SuspendTenant(c.Request.Context(), tenantID); err != nil {
		logger.Error("Failed to suspend tenant", "tenant_id", tenantID, "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to suspend tenant"})
		return
	}

	c.JSON(200, gin.H{"message": "Tenant suspended successfully"})
}

// ActivateTenant handles POST /admin/api/v1/tenants/:id/activate
func (h *Handler) ActivateTenant(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tenant ID"})
		return
	}

	if err := h.adminService.ActivateTenant(c.Request.Context(), tenantID); err != nil {
		logger.Error("Failed to activate tenant", "tenant_id", tenantID, "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to activate tenant"})
		return
	}

	c.JSON(200, gin.H{"message": "Tenant activated successfully"})
}

// UpdateTenantTierRequest represents tier update request
type UpdateTenantTierRequest struct {
	Tier string `json:"tier" binding:"required,oneof=free basic pro enterprise"`
}

// UpdateTenantTier handles PUT /admin/api/v1/tenants/:id/tier
func (h *Handler) UpdateTenantTier(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req UpdateTenantTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	tier := models.SubscriptionTier(req.Tier)
	if err := h.adminService.UpdateTenantTier(c.Request.Context(), tenantID, tier); err != nil {
		logger.Error("Failed to update tenant tier", "tenant_id", tenantID, "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to update tier"})
		return
	}

	c.JSON(200, gin.H{"message": "Tenant tier updated successfully", "tier": req.Tier})
}

// ==================== USERS ====================

// ListUsers handles GET /admin/api/v1/users
func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.adminService.GetAllUsers(c.Request.Context(), page, limit)
	if err != nil {
		logger.Error("Failed to list users", "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to list users"})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// ==================== HEALTH ====================

// GetHealth handles GET /admin/api/v1/health
func (h *Handler) GetHealth(c *gin.Context) {
	health := h.adminService.GetSystemHealth(c.Request.Context())

	status := 200
	if health.Status != "healthy" {
		status = 503
	}

	c.JSON(status, health)
}
