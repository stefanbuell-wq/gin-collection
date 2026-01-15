package router

import (
	"github.com/gin-gonic/gin"
	adminHandler "github.com/yourusername/gin-collection-saas/internal/delivery/http/handler/admin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// AdminRouterConfig holds admin router configuration
type AdminRouterConfig struct {
	AdminHandler        *adminHandler.Handler
	PlatformAdminMiddle *middleware.PlatformAdminMiddleware
	AllowedOrigins      []string
}

// SetupAdminRoutes configures all admin routes on an existing router
func SetupAdminRoutes(r *gin.Engine, cfg *AdminRouterConfig) {
	// Admin API routes
	adminAPI := r.Group("/admin/api/v1")
	{
		// Public admin routes (login)
		auth := adminAPI.Group("/auth")
		{
			auth.POST("/login", cfg.AdminHandler.Login)
		}

		// Protected admin routes
		protected := adminAPI.Group("")
		protected.Use(cfg.PlatformAdminMiddle.RequirePlatformAdmin())
		{
			// Auth
			protected.GET("/auth/me", cfg.AdminHandler.Me)
			protected.POST("/auth/change-password", cfg.AdminHandler.ChangePassword)

			// Statistics
			protected.GET("/stats/overview", cfg.AdminHandler.GetStats)

			// Tenants
			tenants := protected.Group("/tenants")
			{
				tenants.GET("", cfg.AdminHandler.ListTenants)
				tenants.GET("/:id", cfg.AdminHandler.GetTenant)
				tenants.POST("/:id/suspend", cfg.AdminHandler.SuspendTenant)
				tenants.POST("/:id/activate", cfg.AdminHandler.ActivateTenant)
				tenants.PUT("/:id/tier", cfg.AdminHandler.UpdateTenantTier)
			}

			// Users
			users := protected.Group("/users")
			{
				users.GET("", cfg.AdminHandler.ListUsers)
			}

			// Health
			protected.GET("/health", cfg.AdminHandler.GetHealth)
		}
	}

	logger.Info("Admin routes configured")
}
