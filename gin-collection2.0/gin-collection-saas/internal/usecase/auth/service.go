package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/external"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
	"github.com/yourusername/gin-collection-saas/pkg/utils"
)

// Token expiry duration for password reset
const PasswordResetTokenExpiry = 1 * time.Hour

// Service handles authentication business logic
type Service struct {
	userRepo          repositories.UserRepository
	tenantRepo        repositories.TenantRepository
	passwordResetRepo repositories.PasswordResetRepository
	emailClient       *external.EmailClient
	baseURL           string
	jwtSecret         string
	jwtExpiration     time.Duration
	tokenBlacklist    *utils.TokenBlacklist
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

// SetPasswordResetRepo sets the password reset repository (optional dependency)
func (s *Service) SetPasswordResetRepo(repo repositories.PasswordResetRepository) {
	s.passwordResetRepo = repo
}

// SetEmailClient sets the email client (optional dependency)
func (s *Service) SetEmailClient(client *external.EmailClient) {
	s.emailClient = client
}

// SetBaseURL sets the base URL for password reset links
func (s *Service) SetBaseURL(baseURL string) {
	s.baseURL = baseURL
}

// SetTokenBlacklist sets the token blacklist for token invalidation
func (s *Service) SetTokenBlacklist(blacklist *utils.TokenBlacklist) {
	s.tokenBlacklist = blacklist
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

// LoginByEmail authenticates a user by email only (without knowing tenant upfront)
// This is used for login from localhost or when no subdomain is available
func (s *Service) LoginByEmail(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	logger.Info("User login attempt by email only", "email", req.Email)

	// Get user by email (across all tenants)
	user, err := s.userRepo.GetByEmailGlobal(ctx, req.Email)
	if err != nil {
		logger.Debug("User not found by email", "email", req.Email)
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
	tenant, err := s.tenantRepo.GetByID(ctx, user.TenantID)
	if err != nil {
		logger.Error("Failed to get tenant", "tenant_id", user.TenantID, "error", err.Error())
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

	logger.Info("User logged in successfully via email", "user_id", user.ID, "tenant_id", tenant.ID)

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

// UpdateProfile updates the current user's profile information
func (s *Service) UpdateProfile(ctx context.Context, userID int64, firstName, lastName *string) (*models.User, error) {
	logger.Info("Updating user profile", "user_id", userID)

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if firstName != nil {
		user.FirstName = firstName
	}
	if lastName != nil {
		user.LastName = lastName
	}

	// Save user
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user", "error", err.Error())
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("User profile updated successfully", "user_id", userID)

	return user, nil
}

// ChangePassword changes the current user's password
func (s *Service) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	logger.Info("Changing user password", "user_id", userID)

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
		logger.Debug("Invalid current password", "user_id", userID)
		return errors.ErrInvalidCredentials
	}

	// Hash new password
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.Error("Failed to hash new password", "error", err.Error())
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = newHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user password", "error", err.Error())
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all existing tokens for this user
	if s.tokenBlacklist != nil {
		if err := s.tokenBlacklist.RevokeAllUserTokens(ctx, userID, time.Now()); err != nil {
			logger.Error("Failed to revoke user tokens", "error", err.Error(), "user_id", userID)
			// Don't fail the password change, just log the error
		} else {
			logger.Info("All user tokens revoked after password change", "user_id", userID)
		}
	}

	logger.Info("User password changed successfully", "user_id", userID)

	return nil
}

// GetCurrentUser returns the current user's information
func (s *Service) GetCurrentUser(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// RequestPasswordReset initiates a password reset for a user
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	logger.Info("Password reset requested", "email", email)

	// Check if password reset repo is configured
	if s.passwordResetRepo == nil {
		return fmt.Errorf("password reset not configured")
	}

	// Find user by email (across all tenants)
	user, err := s.userRepo.GetByEmailGlobal(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not for security
		logger.Debug("User not found for password reset", "email", email)
		return nil // Return success even if user not found
	}

	// Check if user is active
	if !user.IsActive {
		logger.Debug("Inactive user requested password reset", "user_id", user.ID)
		return nil // Return success even if user is inactive
	}

	// Delete any existing tokens for this user
	if err := s.passwordResetRepo.DeleteByUserID(ctx, user.ID); err != nil {
		logger.Error("Failed to delete existing tokens", "error", err.Error())
		// Continue anyway
	}

	// Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		logger.Error("Failed to generate token", "error", err.Error())
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Create password reset token
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(PasswordResetTokenExpiry),
	}

	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		logger.Error("Failed to create password reset token", "error", err.Error())
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send email
	if s.emailClient != nil {
		recipientName := ""
		if user.FirstName != nil {
			recipientName = *user.FirstName
		}

		resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, token)

		emailData := &external.PasswordResetData{
			RecipientName: recipientName,
			ResetLink:     resetLink,
			ExpiresIn:     "1 Stunde",
		}

		if err := s.emailClient.SendPasswordReset(email, emailData); err != nil {
			logger.Error("Failed to send password reset email", "error", err.Error(), "email", email)
			// Don't fail the request, token is created
		} else {
			logger.Info("Password reset email sent", "email", email)
		}
	}

	logger.Info("Password reset token created", "user_id", user.ID)
	return nil
}

// ResetPassword resets a user's password using a valid token
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	logger.Info("Password reset attempt")

	// Check if password reset repo is configured
	if s.passwordResetRepo == nil {
		return fmt.Errorf("password reset not configured")
	}

	// Get token
	resetToken, err := s.passwordResetRepo.GetByToken(ctx, token)
	if err != nil {
		logger.Debug("Invalid password reset token")
		return errors.ErrInvalidToken
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		logger.Debug("Password reset token expired or used", "token_id", resetToken.ID)
		return errors.ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		logger.Error("Failed to get user for password reset", "error", err.Error())
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Hash new password
	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.Error("Failed to hash new password", "error", err.Error())
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = newHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user password", "error", err.Error())
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(ctx, resetToken.ID); err != nil {
		logger.Error("Failed to mark token as used", "error", err.Error())
		// Don't fail, password is already changed
	}

	// Invalidate all existing tokens for this user
	if s.tokenBlacklist != nil {
		if err := s.tokenBlacklist.RevokeAllUserTokens(ctx, user.ID, time.Now()); err != nil {
			logger.Error("Failed to revoke user tokens", "error", err.Error(), "user_id", user.ID)
			// Don't fail the password reset, just log the error
		} else {
			logger.Info("All user tokens revoked after password reset", "user_id", user.ID)
		}
	}

	logger.Info("Password reset successful", "user_id", user.ID)
	return nil
}

// ValidateResetToken checks if a password reset token is valid
func (s *Service) ValidateResetToken(ctx context.Context, token string) (bool, error) {
	if s.passwordResetRepo == nil {
		return false, fmt.Errorf("password reset not configured")
	}

	resetToken, err := s.passwordResetRepo.GetByToken(ctx, token)
	if err != nil {
		return false, nil
	}

	return resetToken.IsValid(), nil
}
