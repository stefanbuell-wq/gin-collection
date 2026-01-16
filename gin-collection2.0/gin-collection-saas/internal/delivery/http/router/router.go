package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/handler"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// RouterConfig holds router configuration
type RouterConfig struct {
	AuthHandler          *handler.AuthHandler
	GinHandler           *handler.GinHandler
	GinReferenceHandler  *handler.GinReferenceHandler
	SubscriptionHandler  *handler.SubscriptionHandler
	WebhookHandler       *handler.WebhookHandler
	BotanicalHandler     *handler.BotanicalHandler
	CocktailHandler      *handler.CocktailHandler
	PhotoHandler         *handler.PhotoHandler
	UserHandler          *handler.UserHandler
	TenantHandler        *handler.TenantHandler
	AuthMiddleware       *middleware.AuthMiddleware
	TenantMiddleware     *middleware.TenantMiddleware
	TierEnforcement      *middleware.TierEnforcementMiddleware
	AllowedOrigins       []string
}

// Setup configures all routes
func Setup(cfg *RouterConfig) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// Health check endpoints (no auth required)
	r.GET("/health", healthCheck)
	r.GET("/ready", readyCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (tenant middleware only, no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)

			// Login with optional tenant (allows login via email only for localhost)
			auth.POST("/login", cfg.TenantMiddleware.OptionalTenant(), cfg.AuthHandler.Login)

			// Refresh token
			auth.POST("/refresh", cfg.AuthHandler.RefreshToken)

			// Logout (requires auth)
			auth.POST("/logout", cfg.AuthMiddleware.RequireAuth(), cfg.AuthHandler.Logout)

			// Protected auth routes (require tenant + auth)
			authProtected := auth.Group("")
			authProtected.Use(cfg.TenantMiddleware.ExtractTenant())
			authProtected.Use(cfg.AuthMiddleware.RequireAuth())
			{
				authProtected.GET("/me", cfg.AuthHandler.GetMe)
				authProtected.PUT("/profile", cfg.AuthHandler.UpdateProfile)
				authProtected.POST("/change-password", cfg.AuthHandler.ChangePassword)
			}
		}

		// Protected routes (require tenant + auth)
		protected := v1.Group("")
		protected.Use(cfg.TenantMiddleware.ExtractTenant())
		protected.Use(cfg.AuthMiddleware.RequireAuth())
		{
			// Tenants
			tenants := protected.Group("/tenants")
			{
				tenants.GET("/current", cfg.TenantHandler.GetCurrent)
				tenants.PUT("/current", cfg.TenantHandler.UpdateCurrent)
				tenants.GET("/usage", cfg.TenantHandler.GetUsage)
			}

			// Subscriptions
			subscriptions := protected.Group("/subscriptions")
			{
				subscriptions.GET("/current", cfg.SubscriptionHandler.GetCurrent)
				subscriptions.GET("/plans", cfg.SubscriptionHandler.GetPlans)
				subscriptions.POST("/upgrade", cfg.SubscriptionHandler.Upgrade)
				subscriptions.POST("/activate", cfg.SubscriptionHandler.Activate)
				subscriptions.POST("/cancel", middleware.RequirePermission("manage_billing"), cfg.SubscriptionHandler.Cancel)
			}

			// Gins
			gins := protected.Group("/gins")
			{
				gins.GET("", cfg.GinHandler.List)
				gins.POST("", cfg.TierEnforcement.CheckGinLimit(), cfg.GinHandler.Create)
				gins.GET("/search", cfg.GinHandler.Search)
				gins.GET("/stats", cfg.GinHandler.Stats)
				gins.POST("/export", cfg.TierEnforcement.RequireFeature("export"), cfg.GinHandler.Export)
				gins.POST("/import", cfg.TierEnforcement.RequireFeature("import"), cfg.GinHandler.Import)
				gins.GET("/:id", cfg.GinHandler.Get)
				gins.PUT("/:id", cfg.GinHandler.Update)
				gins.DELETE("/:id", middleware.RequirePermission("delete"), cfg.GinHandler.Delete)
				gins.GET("/:id/suggestions", cfg.TierEnforcement.RequireFeature("ai_suggestions"), cfg.GinHandler.Suggestions)

				// Gin Botanicals (Pro+ feature)
				gins.GET("/:id/botanicals", cfg.TierEnforcement.RequireFeature("botanicals"), cfg.BotanicalHandler.GetGinBotanicals)
				gins.PUT("/:id/botanicals", cfg.TierEnforcement.RequireFeature("botanicals"), cfg.BotanicalHandler.UpdateGinBotanicals)

				// Gin Cocktails (Pro+ feature)
				gins.GET("/:id/cocktails", cfg.TierEnforcement.RequireFeature("cocktails"), cfg.CocktailHandler.GetGinCocktails)
				gins.POST("/:id/cocktails/:cocktail_id", cfg.TierEnforcement.RequireFeature("cocktails"), cfg.CocktailHandler.LinkCocktail)
				gins.DELETE("/:id/cocktails/:cocktail_id", cfg.TierEnforcement.RequireFeature("cocktails"), cfg.CocktailHandler.UnlinkCocktail)

				// Gin Photos
				gins.GET("/:id/photos", cfg.PhotoHandler.GetPhotos)
				gins.POST("/:id/photos", cfg.PhotoHandler.Upload)
				gins.DELETE("/:id/photos/:photo_id", cfg.PhotoHandler.Delete)
				gins.PUT("/:id/photos/:photo_id/primary", cfg.PhotoHandler.SetPrimary)
			}

			// Botanicals (reference data, available to all)
			botanicals := protected.Group("/botanicals")
			{
				botanicals.GET("", cfg.BotanicalHandler.GetAll)
			}

			// Cocktails (reference data, available to all)
			cocktails := protected.Group("/cocktails")
			{
				cocktails.GET("", cfg.CocktailHandler.GetAll)
				cocktails.GET("/:id", cfg.CocktailHandler.GetByID)
			}

			// Users (Enterprise only)
			users := protected.Group("/users")
			users.Use(middleware.RequireRole(models.RoleOwner, models.RoleAdmin))
			{
				users.GET("", cfg.UserHandler.List)
				users.POST("/invite", cfg.UserHandler.Invite)
				users.PUT("/:id", cfg.UserHandler.Update)
				users.DELETE("/:id", cfg.UserHandler.Delete)
				users.POST("/:id/api-key", cfg.UserHandler.GenerateAPIKey)
				users.DELETE("/:id/api-key", cfg.UserHandler.RevokeAPIKey)
			}
		}

		// Gin References (public catalog, requires auth but no tenant)
		ginRefs := v1.Group("/gin-references")
		ginRefs.Use(cfg.AuthMiddleware.RequireAuth())
		{
			ginRefs.GET("", cfg.GinReferenceHandler.Search)
			ginRefs.GET("/filters", cfg.GinReferenceHandler.GetFilters)
			ginRefs.GET("/:id", cfg.GinReferenceHandler.GetByID)
		}

		// Webhooks (no auth, validated by signature)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/paypal", cfg.WebhookHandler.PayPal)
		}
	}

	return r
}

// healthCheck returns OK if server is running
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// readyCheck checks if dependencies (DB, Redis) are ready
func readyCheck(c *gin.Context) {
	// TODO: Implement actual health checks
	c.JSON(200, gin.H{
		"status": "ready",
	})
}

// placeholderHandler is a temporary handler for routes not yet implemented
func placeholderHandler(c *gin.Context) {
	c.JSON(501, gin.H{
		"error": "Not yet implemented",
	})
}
