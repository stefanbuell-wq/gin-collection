package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
	"github.com/yourusername/gin-collection-saas/pkg/utils"
)

// Service handles authentication business logic
type Service struct {
	userRepo   repositories.UserRepository
	tenantRepo repositories.TenantRepository
	jwtSecret  string
	jwtExpiration time.Duration
}

// NewService creates a new auth service
func NewService(
	userRepo repositories.UserRepository,
	tenantRepo repositories.TenantRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
) *Service {
	return &Service{
		userRepo:      userRepo,
		tenantRepo:    tenantRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// Register registers a new tenant with an owner user
func (s *Service) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	logger.Info("Registering new tenant", "subdomain", req.Subdomain, "email", req.Email)

	// Validate subdomain is not taken
	existingTenant, err := s.tenantRepo.GetBySubdomain(ctx, req.Subdomain)
	if err == nil && existingTenant != nil {
		return nil, errors.ErrSubdomainTaken
	}

	// Validate email is not already used
	// Note: We can't check across tenants easily here, but we'll check when tenant is created

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", "error", err.Error())
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create tenant
	tenant := &models.Tenant{
		UUID:      uuid.New().String(),
		Name:      req.TenantName,
		Subdomain: req.Subdomain,
		Tier:      models.TierFree, // Start with free tier
		Status:    models.TenantStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		logger.Error("Failed to create tenant", "error", err.Error())
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// Create owner user
	user := &models.User{
		TenantID:     tenant.ID,
		UUID:         uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		Role:         models.RoleOwner,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user", "error", err.Error())
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		tenant.ID,
		user.Email,
		string(user.Role),
		s.jwtSecret,
		s.jwtExpiration,
	)
	if err != nil {
		logger.Error("Failed to generate token", "error", err.Error())
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		tenant.ID,
		user.Email,
		s.jwtSecret,
	)
	if err != nil {
		logger.Error("Failed to generate refresh token", "error", err.Error())
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	logger.Info("User registered successfully", "user_id", user.ID, "tenant_id", tenant.ID)

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
		Tenant:       tenant,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, req *models.LoginRequest, tenantID int64) (*models.AuthResponse, error) {
	logger.Info("User login attempt", "email", req.Email, "tenant_id", tenantID)

	// Get user by email and tenant
	user, err := s.userRepo.GetByEmail(ctx, tenantID, req.Email)
	if err != nil {
		logger.Debug("User not found", "email", req.Email, "tenant_id", tenantID)
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		logger.Debug("Inactive user attempted login", "user_id", user.ID)
		return nil, errors.ErrForbidden
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		logger.Debug("Invalid password", "user_id", user.ID)
		return nil, errors.ErrInvalidCredentials
	}

	// Get tenant
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		logger.Error("Failed to get tenant", "tenant_id", tenantID, "error", err.Error())
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Check tenant status
	if tenant.Status != models.TenantStatusActive {
		logger.Debug("Suspended tenant attempted login", "tenant_id", tenant.ID)
		return nil, errors.ErrTenantSuspended
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		tenant.ID,
		user.Email,
		string(user.Role),
		s.jwtSecret,
		s.jwtExpiration,
	)
	if err != nil {
		logger.Error("Failed to generate token", "error", err.Error())
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		tenant.ID,
		user.Email,
		s.jwtSecret,
	)
	if err != nil {
		logger.Error("Failed to generate refresh token", "error", err.Error())
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		logger.Error("Failed to update last login", "user_id", user.ID, "error", err.Error())
		// Don't fail the login, just log the error
	}

	logger.Info("User logged in successfully", "user_id", user.ID)

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
		Tenant:       tenant,
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", errors.ErrInvalidToken
	}

	// Generate new access token
	newToken, err := utils.GenerateToken(
		claims.UserID,
		claims.TenantID,
		claims.Email,
		claims.Role,
		s.jwtSecret,
		s.jwtExpiration,
	)
	if err != nil {
		logger.Error("Failed to generate new token", "error", err.Error())
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	return newToken, nil
}
