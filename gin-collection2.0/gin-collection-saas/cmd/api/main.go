package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/gin-collection-saas/internal/delivery/http/handler"
	adminHandler "github.com/yourusername/gin-collection-saas/internal/delivery/http/handler/admin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/router"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/cache"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/database"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/external"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/storage"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
	adminUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/admin"
	"github.com/yourusername/gin-collection-saas/internal/usecase/auth"
	botanicalUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/botanical"
	cocktailUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/cocktail"
	ginUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/gin"
	photoUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/photo"
	subscriptionUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/subscription"
	tastingUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/tasting"
	userUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/user"
	"github.com/yourusername/gin-collection-saas/pkg/config"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.App.LogLevel)
	logger.Info("Starting Gin Collection SaaS API", "version", "1.0.0", "env", cfg.App.Env)

	// Connect to MySQL
	db, err := database.NewMySQL(
		cfg.Database.DSN(),
		cfg.Database.MaxConns,
		cfg.Database.MaxIdle,
	)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err.Error())
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Connected to MySQL database", "host", cfg.Database.Host, "database", cfg.Database.Name)

	// Connect to Redis
	redisClient, err := cache.NewRedisClient(&cache.RedisConfig{
		URL:      cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		logger.Warn("Failed to connect to Redis - rate limiting disabled", "error", err.Error())
		redisClient = nil
	} else {
		defer redisClient.Close()
		logger.Info("Connected to Redis", "url", cfg.Redis.URL)
	}

	// Initialize Tenant Router (for hybrid multi-tenancy)
	tenantRouter := database.NewTenantRouter(db)
	defer tenantRouter.Close()

	// Initialize repositories
	tenantRepo := mysql.NewTenantRepository(db)
	userRepo := mysql.NewUserRepository(db)
	ginRepo := mysql.NewGinRepository(db)
	ginReferenceRepo := mysql.NewGinReferenceRepository(db)
	usageMetricsRepo := mysql.NewUsageMetricsRepository(db)
	subscriptionRepo := mysql.NewSubscriptionRepository(db)
	botanicalRepo := mysql.NewBotanicalRepository(db)
	cocktailRepo := mysql.NewCocktailRepository(db)
	photoRepo := mysql.NewPhotoRepository(db)
	auditLogRepo := mysql.NewAuditLogRepository(db)
	platformAdminRepo := mysql.NewPlatformAdminRepository(db)
	tastingRepo := mysql.NewTastingSessionRepository(db)
	passwordResetRepo := mysql.NewPasswordResetRepository(db)

	logger.Info("Repositories initialized")

	// Initialize PayPal client
	paypalClient := external.NewPayPalClient(&external.PayPalConfig{
		ClientID:     cfg.PayPal.ClientID,
		ClientSecret: cfg.PayPal.ClientSecret,
		Mode:         cfg.PayPal.Mode,
	})

	// Initialize storage client (S3 or Local fallback)
	var storageClient storage.Storage
	if cfg.S3.AccessKeyID == "" || cfg.S3.SecretAccessKey == "" {
		// Use local storage when S3 is not configured
		storageBaseURL := cfg.Storage.BaseURL
		if storageBaseURL == "" {
			storageBaseURL = cfg.App.BaseURL + "/uploads"
		}
		storageClient, err = storage.NewLocalStorage(&storage.LocalStorageConfig{
			BasePath: cfg.Storage.BasePath,
			BaseURL:  storageBaseURL,
		})
		if err != nil {
			logger.Error("Failed to initialize local storage", "error", err.Error())
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		logger.Info("Using local file storage", "path", cfg.Storage.BasePath, "url", storageBaseURL)
	} else {
		// Use S3 storage
		storageClient, err = storage.NewS3Client(&storage.S3Config{
			Bucket:          cfg.S3.Bucket,
			Region:          cfg.S3.Region,
			Endpoint:        cfg.S3.Endpoint,
			AccessKeyID:     cfg.S3.AccessKeyID,
			SecretAccessKey: cfg.S3.SecretAccessKey,
		})
		if err != nil {
			logger.Error("Failed to initialize S3 client", "error", err.Error())
			log.Fatalf("Failed to initialize S3 client: %v", err)
		}
		logger.Info("Using S3 storage", "bucket", cfg.S3.Bucket, "region", cfg.S3.Region)
	}

	// Initialize Email client
	emailClient := external.NewEmailClient(&external.EmailConfig{
		Host:       cfg.SMTP.Host,
		Port:       cfg.SMTP.Port,
		Username:   cfg.SMTP.Username,
		Password:   cfg.SMTP.Password,
		FromEmail:  cfg.SMTP.FromEmail,
		FromName:   cfg.SMTP.FromName,
		TLS:        cfg.SMTP.TLS,
		SkipVerify: cfg.SMTP.SkipVerify,
	})
	logger.Info("Email client initialized", "host", cfg.SMTP.Host, "from", cfg.SMTP.FromEmail)

	// Initialize AI client (Ollama or Anthropic)
	aiClient := external.NewAIClient(&external.AIClientConfig{
		Provider:        cfg.AI.Provider,
		OllamaURL:       cfg.AI.OllamaURL,
		Model:           cfg.AI.Model,
		AnthropicAPIKey: cfg.AI.AnthropicAPIKey,
		Enabled:         cfg.AI.Enabled,
	})
	if aiClient.IsEnabled() {
		logger.Info("AI client initialized", "provider", cfg.AI.Provider, "model", cfg.AI.Model)
	} else {
		logger.Info("AI client disabled")
	}

	// Initialize use cases
	authService := auth.NewService(
		userRepo,
		tenantRepo,
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
	)
	// Set optional dependencies for password reset
	authService.SetPasswordResetRepo(passwordResetRepo)
	authService.SetEmailClient(emailClient)
	authService.SetBaseURL(cfg.App.BaseURL)

	ginService := ginUsecase.NewService(
		ginRepo,
		usageMetricsRepo,
	)

	subscriptionService := subscriptionUsecase.NewService(
		subscriptionRepo,
		tenantRepo,
		paypalClient,
		cfg.App.BaseURL,
	)

	botanicalService := botanicalUsecase.NewService(
		botanicalRepo,
		ginRepo,
	)

	cocktailService := cocktailUsecase.NewService(
		cocktailRepo,
		ginRepo,
	)

	photoService := photoUsecase.NewService(
		photoRepo,
		ginRepo,
		usageMetricsRepo,
		tenantRepo,
		storageClient,
	)

	userService := userUsecase.NewService(
		userRepo,
		tenantRepo,
		auditLogRepo,
		emailClient,
		cfg.App.BaseURL,
	)

	tastingService := tastingUsecase.NewService(
		tastingRepo,
		ginRepo,
	)

	// Initialize Platform Admin Service
	adminService := adminUsecase.NewService(
		platformAdminRepo,
		db,
		cfg.JWT.Secret,
	)

	logger.Info("Services initialized")

	// Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(authService)
	ginHandler := handler.NewGinHandler(ginService)
	ginReferenceHandler := handler.NewGinReferenceHandler(ginReferenceRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	webhookHandler := handler.NewWebhookHandler(subscriptionService)
	botanicalHandler := handler.NewBotanicalHandler(botanicalService)
	cocktailHandler := handler.NewCocktailHandler(cocktailService)
	photoHandler := handler.NewPhotoHandler(photoService)
	userHandler := handler.NewUserHandler(userService)
	tenantHandler := handler.NewTenantHandler(tenantRepo, usageMetricsRepo)
	aiHandler := handler.NewAIHandler(aiClient)
	tastingHandler := handler.NewTastingHandler(tastingService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, userRepo)
	tenantMiddleware := middleware.NewTenantMiddleware(tenantRepo)
	tierEnforcement := middleware.NewTierEnforcementMiddleware(usageMetricsRepo, ginRepo)
	platformAdminMiddleware := middleware.NewPlatformAdminMiddleware(adminService)

	// Initialize rate limiting middleware (optional - requires Redis)
	var rateLimitMiddleware *middleware.RateLimitMiddleware
	if redisClient != nil {
		rateLimitMiddleware = middleware.NewRateLimitMiddleware(redisClient)
		logger.Info("Rate limiting enabled")
	} else {
		logger.Warn("Rate limiting disabled - Redis not available")
	}

	// Initialize Admin handlers
	platformAdminHandler := adminHandler.NewHandler(adminService)

	// Initialize Server handler for deployment management
	// Only enable in production when PROJECT_PATH is set
	var serverHandler *adminHandler.ServerHandler
	if projectPath := os.Getenv("PROJECT_PATH"); projectPath != "" {
		serverHandler = adminHandler.NewServerHandler(projectPath)
		logger.Info("Server management enabled", "project_path", projectPath)
	}

	logger.Info("Handlers and middleware initialized")

	// Setup router
	routerCfg := &router.RouterConfig{
		AuthHandler:         authHandler,
		GinHandler:          ginHandler,
		GinReferenceHandler: ginReferenceHandler,
		SubscriptionHandler: subscriptionHandler,
		WebhookHandler:      webhookHandler,
		BotanicalHandler:    botanicalHandler,
		CocktailHandler:     cocktailHandler,
		PhotoHandler:        photoHandler,
		UserHandler:         userHandler,
		TenantHandler:       tenantHandler,
		AIHandler:           aiHandler,
		TastingHandler:      tastingHandler,
		AuthMiddleware:      authMiddleware,
		TenantMiddleware:    tenantMiddleware,
		TierEnforcement:     tierEnforcement,
		RateLimitMiddleware: rateLimitMiddleware,
		AllowedOrigins:      cfg.App.AllowedOrigins,
	}

	r := router.Setup(routerCfg)

	// Setup static file serving for local uploads (when using local storage)
	if cfg.S3.AccessKeyID == "" || cfg.S3.SecretAccessKey == "" {
		r.Static("/uploads", cfg.Storage.BasePath)
		logger.Info("Static file serving enabled", "path", "/uploads", "directory", cfg.Storage.BasePath)
	}

	// Setup Admin routes
	adminRouterCfg := &router.AdminRouterConfig{
		AdminHandler:        platformAdminHandler,
		ServerHandler:       serverHandler,
		PlatformAdminMiddle: platformAdminMiddleware,
		AllowedOrigins:      cfg.App.AllowedOrigins,
	}
	router.SetupAdminRoutes(r, adminRouterCfg)

	logger.Info("Admin routes initialized")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info("Starting HTTP server", "addr", addr, "env", cfg.App.Env)

	fmt.Printf("\n")
	fmt.Printf("╔═══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                                                               ║\n")
	fmt.Printf("║           🍸 Gin Collection SaaS API v1.0.0                   ║\n")
	fmt.Printf("║                                                               ║\n")
	fmt.Printf("╠═══════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║  Server:    http://localhost%s                          ║\n", addr)
	fmt.Printf("║  Health:    http://localhost%s/health                   ║\n", addr)
	fmt.Printf("║  Admin:     http://localhost%s/admin/api/v1             ║\n", addr)
	fmt.Printf("║  Env:       %-49s ║\n", cfg.App.Env)
	fmt.Printf("║  Database:  %-49s ║\n", cfg.Database.Name)
	fmt.Printf("╚═══════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("\n")

	if err := r.Run(addr); err != nil {
		logger.Error("Failed to start server", "error", err.Error())
		log.Fatalf("Failed to start server: %v", err)
	}
}

// init loads environment variables from .env file if present
func init() {
	// Try to load .env file (optional)
	if _, err := os.Stat(".env"); err == nil {
		// .env exists, load it
		// Note: We're using simple env variables, no external lib needed
		logger.Debug(".env file found")
	}
}
