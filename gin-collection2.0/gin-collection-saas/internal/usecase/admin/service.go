package admin

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
	"github.com/yourusername/gin-collection-saas/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAdminNotActive     = errors.New("admin account is not active")
	ErrAdminNotFound      = errors.New("admin not found")
)

// AdminJWTClaims extends standard JWT claims for platform admins
type AdminJWTClaims struct {
	AdminID         int64  `json:"admin_id"`
	Email           string `json:"email"`
	IsPlatformAdmin bool   `json:"is_platform_admin"`
	jwt.RegisteredClaims
}

// Service handles platform admin business logic
type Service struct {
	adminRepo *mysql.PlatformAdminRepository
	db        *sql.DB
	jwtSecret string
}

// NewService creates a new admin service
func NewService(adminRepo *mysql.PlatformAdminRepository, db *sql.DB, jwtSecret string) *Service {
	return &Service{
		adminRepo: adminRepo,
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Login authenticates a platform admin and returns a JWT token
func (s *Service) Login(ctx context.Context, email, password string) (*models.PlatformAdminAuthResponse, error) {
	logger.Info("Platform admin login attempt", "email", email)

	// Get admin by email
	admin, err := s.adminRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("Failed to get admin by email", "error", err.Error())
		return nil, err
	}
	if admin == nil {
		logger.Warn("Admin not found", "email", email)
		return nil, ErrInvalidCredentials
	}

	// Check if admin is active
	if !admin.IsActive {
		logger.Warn("Admin account not active", "email", email)
		return nil, ErrAdminNotActive
	}

	// Verify password
	if !utils.CheckPasswordHash(password, admin.PasswordHash) {
		logger.Warn("Invalid password", "email", email)
		return nil, ErrInvalidCredentials
	}

	// Update last login
	if err := s.adminRepo.UpdateLastLogin(ctx, admin.ID); err != nil {
		logger.Error("Failed to update last login", "error", err.Error())
		// Don't fail login for this
	}

	// Generate JWT token
	token, err := s.generateAdminToken(admin)
	if err != nil {
		logger.Error("Failed to generate token", "error", err.Error())
		return nil, err
	}

	logger.Info("Platform admin login successful", "admin_id", admin.ID, "email", email)

	return &models.PlatformAdminAuthResponse{
		Token: token,
		Admin: admin,
	}, nil
}

// generateAdminToken creates a JWT token for a platform admin
func (s *Service) generateAdminToken(admin *models.PlatformAdmin) (string, error) {
	claims := AdminJWTClaims{
		AdminID:         admin.ID,
		Email:           admin.Email,
		IsPlatformAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-collection-platform",
			Subject:   "platform-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateAdminToken validates a platform admin JWT token
func (s *Service) ValidateAdminToken(tokenString string) (*AdminJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AdminJWTClaims); ok && token.Valid {
		// Verify it's actually a platform admin token
		if !claims.IsPlatformAdmin {
			return nil, errors.New("not a platform admin token")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetAdmin retrieves a platform admin by ID
func (s *Service) GetAdmin(ctx context.Context, id int64) (*models.PlatformAdmin, error) {
	return s.adminRepo.GetByID(ctx, id)
}

// ChangePassword changes an admin's password
func (s *Service) ChangePassword(ctx context.Context, adminID int64, oldPassword, newPassword string) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrAdminNotFound
	}

	// Verify old password
	if !utils.CheckPasswordHash(oldPassword, admin.PasswordHash) {
		return ErrInvalidCredentials
	}

	// Hash new password
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.adminRepo.UpdatePassword(ctx, adminID, newHash)
}

// GetPlatformStats retrieves overall platform statistics
func (s *Service) GetPlatformStats(ctx context.Context) (*models.PlatformStats, error) {
	return s.adminRepo.GetPlatformStats(ctx)
}

// GetAllTenants retrieves all tenants with pagination
func (s *Service) GetAllTenants(ctx context.Context, page, limit int) ([]*models.TenantWithStats, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.adminRepo.GetAllTenants(ctx, limit, offset)
}

// GetAllUsers retrieves all users with pagination
func (s *Service) GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.adminRepo.GetAllUsers(ctx, limit, offset)
}

// SuspendTenant suspends a tenant
func (s *Service) SuspendTenant(ctx context.Context, tenantID int64) error {
	logger.Info("Suspending tenant", "tenant_id", tenantID)
	return s.adminRepo.UpdateTenantStatus(ctx, tenantID, models.TenantStatusSuspended)
}

// ActivateTenant activates a tenant
func (s *Service) ActivateTenant(ctx context.Context, tenantID int64) error {
	logger.Info("Activating tenant", "tenant_id", tenantID)
	return s.adminRepo.UpdateTenantStatus(ctx, tenantID, models.TenantStatusActive)
}

// UpdateTenantTier updates a tenant's subscription tier
func (s *Service) UpdateTenantTier(ctx context.Context, tenantID int64, tier models.SubscriptionTier) error {
	logger.Info("Updating tenant tier", "tenant_id", tenantID, "tier", tier)
	return s.adminRepo.UpdateTenantTier(ctx, tenantID, tier)
}

// GetSystemHealth retrieves system health status
func (s *Service) GetSystemHealth(ctx context.Context) *models.SystemHealth {
	health := &models.SystemHealth{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	// Check database
	start := time.Now()
	if err := s.db.PingContext(ctx); err != nil {
		health.Database = models.HealthStatus{
			Status: "unhealthy",
			Error:  err.Error(),
		}
		health.Status = "degraded"
	} else {
		health.Database = models.HealthStatus{
			Status:  "healthy",
			Latency: time.Since(start).String(),
		}
	}

	// Redis check would go here if we had a redis client in the service

	return health
}
