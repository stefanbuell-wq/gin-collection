package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/external"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Service handles user management business logic (Enterprise feature)
type Service struct {
	userRepo      repositories.UserRepository
	tenantRepo    repositories.TenantRepository
	auditLogRepo  repositories.AuditLogRepository
	emailClient   *external.EmailClient
	baseURL       string
}

// NewService creates a new user management service
func NewService(
	userRepo repositories.UserRepository,
	tenantRepo repositories.TenantRepository,
	auditLogRepo repositories.AuditLogRepository,
	emailClient *external.EmailClient,
	baseURL string,
) *Service {
	return &Service{
		userRepo:     userRepo,
		tenantRepo:   tenantRepo,
		auditLogRepo: auditLogRepo,
		emailClient:  emailClient,
		baseURL:      baseURL,
	}
}

// ListUsers lists all users in a tenant (Enterprise only)
func (s *Service) ListUsers(ctx context.Context, tenantID, requesterUserID int64) ([]*models.User, error) {
	logger.Info("Listing users", "tenant_id", tenantID)

	// Verify tenant is Enterprise
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenant.Tier != "enterprise" {
		return nil, errors.ErrMultiUserNotAllowed
	}

	// Get all users
	users, err := s.userRepo.List(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	logger.Info("Users listed successfully", "tenant_id", tenantID, "count", len(users))

	return users, nil
}

// InviteUser invites a new user to the tenant (Enterprise only)
func (s *Service) InviteUser(ctx context.Context, tenantID, inviterUserID int64, email, firstName, lastName string, role models.UserRole) (*models.User, error) {
	logger.Info("Inviting user", "tenant_id", tenantID, "email", email, "role", role)

	// Verify tenant is Enterprise
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenant.Tier != "enterprise" {
		return nil, errors.ErrMultiUserNotAllowed
	}

	// Check if email already exists in tenant
	existingUser, _ := s.userRepo.GetByEmail(ctx, tenantID, email)
	if existingUser != nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Generate a temporary password (user should reset on first login)
	tempPassword := generateTempPassword()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		TenantID:     tenantID,
		Email:        email,
		PasswordHash: string(passwordHash),
		FirstName:    &firstName,
		LastName:     &lastName,
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create audit log
	changes, _ := json.Marshal(map[string]interface{}{
		"email":      email,
		"role":       role,
		"invited_by": inviterUserID,
	})
	changesStr := string(changes)
	auditLog := &models.AuditLog{
		TenantID:   tenantID,
		UserID:     &inviterUserID,
		Action:     string(models.AuditActionInviteUser),
		EntityType: string(models.EntityTypeUser),
		EntityID:   &user.ID,
		Changes:    &changesStr,
	}
	s.auditLogRepo.Create(ctx, auditLog)

	logger.Info("User invited successfully", "user_id", user.ID, "tenant_id", tenantID)

	// Get inviter info for email
	inviter, _ := s.userRepo.GetByID(ctx, inviterUserID)
	inviterName := "Ein Administrator"
	if inviter != nil {
		if inviter.FirstName != nil && inviter.LastName != nil {
			inviterName = fmt.Sprintf("%s %s", *inviter.FirstName, *inviter.LastName)
		} else {
			inviterName = inviter.Email
		}
	}

	// Send invitation email
	recipientName := ""
	if firstName != "" {
		recipientName = firstName
	}

	roleLabels := map[models.UserRole]string{
		models.RoleOwner:  "Inhaber",
		models.RoleAdmin:  "Administrator",
		models.RoleMember: "Mitglied",
		models.RoleViewer: "Betrachter",
	}

	inviteLink := fmt.Sprintf("%s/login?email=%s&invite=true", s.baseURL, email)

	if s.emailClient != nil {
		emailData := &external.UserInvitationData{
			RecipientName:  recipientName,
			RecipientEmail: email,
			InviterName:    inviterName,
			TenantName:     tenant.Name,
			Role:           roleLabels[role],
			InviteLink:     inviteLink,
			ExpiresIn:      "7 Tage",
		}

		if err := s.emailClient.SendUserInvitation(emailData); err != nil {
			// Log error but don't fail the invitation
			logger.Error("Failed to send invitation email", "error", err.Error(), "email", email)
		}
	}

	return user, nil
}

// UpdateUser updates user information (Enterprise only)
func (s *Service) UpdateUser(ctx context.Context, tenantID, requesterUserID, targetUserID int64, email string, firstName, lastName *string, role models.UserRole, isActive bool) (*models.User, error) {
	logger.Info("Updating user", "tenant_id", tenantID, "target_user_id", targetUserID)

	// Verify tenant is Enterprise
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenant.Tier != "enterprise" {
		return nil, errors.ErrMultiUserNotAllowed
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to tenant
	if user.TenantID != tenantID {
		return nil, errors.ErrUserNotInTenant
	}

	// Prevent owner from being deactivated or demoted
	if user.Role == models.RoleOwner {
		return nil, fmt.Errorf("cannot modify owner account")
	}

	// Track changes for audit
	changes := map[string]interface{}{}
	if user.Email != email {
		changes["email"] = map[string]string{"old": user.Email, "new": email}
		user.Email = email
	}
	if role != user.Role {
		changes["role"] = map[string]string{"old": string(user.Role), "new": string(role)}
		user.Role = role
	}
	if isActive != user.IsActive {
		changes["is_active"] = map[string]bool{"old": user.IsActive, "new": isActive}
		user.IsActive = isActive
	}
	if firstName != nil {
		user.FirstName = firstName
	}
	if lastName != nil {
		user.LastName = lastName
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	if len(changes) > 0 {
		changesJSON, _ := json.Marshal(changes)
		changesStr := string(changesJSON)
		auditLog := &models.AuditLog{
			TenantID:   tenantID,
			UserID:     &requesterUserID,
			Action:     string(models.AuditActionUpdateUser),
			EntityType: string(models.EntityTypeUser),
			EntityID:   &targetUserID,
			Changes:    &changesStr,
		}
		s.auditLogRepo.Create(ctx, auditLog)
	}

	logger.Info("User updated successfully", "user_id", targetUserID, "tenant_id", tenantID)

	return user, nil
}

// DeleteUser deletes a user from the tenant (Enterprise only)
func (s *Service) DeleteUser(ctx context.Context, tenantID, requesterUserID, targetUserID int64) error {
	logger.Info("Deleting user", "tenant_id", tenantID, "target_user_id", targetUserID)

	// Verify tenant is Enterprise
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenant.Tier != "enterprise" {
		return errors.ErrMultiUserNotAllowed
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to tenant
	if user.TenantID != tenantID {
		return errors.ErrUserNotInTenant
	}

	// Prevent owner from being deleted
	if user.Role == models.RoleOwner {
		return fmt.Errorf("cannot delete owner account")
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, targetUserID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Create audit log
	changes, _ := json.Marshal(map[string]interface{}{
		"email":      user.Email,
		"role":       user.Role,
		"deleted_by": requesterUserID,
	})
	changesStr := string(changes)
	auditLog := &models.AuditLog{
		TenantID:   tenantID,
		UserID:     &requesterUserID,
		Action:     string(models.AuditActionDeleteUser),
		EntityType: string(models.EntityTypeUser),
		EntityID:   &targetUserID,
		Changes:    &changesStr,
	}
	s.auditLogRepo.Create(ctx, auditLog)

	logger.Info("User deleted successfully", "user_id", targetUserID, "tenant_id", tenantID)

	return nil
}

// GenerateAPIKey generates an API key for a user (Enterprise only)
func (s *Service) GenerateAPIKey(ctx context.Context, tenantID, requesterUserID, targetUserID int64) (string, error) {
	logger.Info("Generating API key", "tenant_id", tenantID, "target_user_id", targetUserID)

	// Verify tenant is Enterprise
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenant.Tier != "enterprise" {
		return "", errors.ErrAPIAccessNotAllowed
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to tenant
	if user.TenantID != tenantID {
		return "", errors.ErrUserNotInTenant
	}

	// Generate API key
	apiKey, err := s.userRepo.GenerateAPIKey(ctx, targetUserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}

	// Create audit log
	auditLog := &models.AuditLog{
		TenantID:   tenantID,
		UserID:     &requesterUserID,
		Action:     string(models.AuditActionGenerateAPIKey),
		EntityType: string(models.EntityTypeUser),
		EntityID:   &targetUserID,
	}
	s.auditLogRepo.Create(ctx, auditLog)

	logger.Info("API key generated successfully", "user_id", targetUserID, "tenant_id", tenantID)

	return apiKey, nil
}

// RevokeAPIKey revokes the API key for a user (Enterprise only)
func (s *Service) RevokeAPIKey(ctx context.Context, tenantID, requesterUserID, targetUserID int64) error {
	logger.Info("Revoking API key", "tenant_id", tenantID, "target_user_id", targetUserID)

	// Verify tenant is Enterprise
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenant.Tier != "enterprise" {
		return errors.ErrAPIAccessNotAllowed
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify user belongs to tenant
	if user.TenantID != tenantID {
		return errors.ErrUserNotInTenant
	}

	// Revoke API key
	if err := s.userRepo.RevokeAPIKey(ctx, targetUserID); err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	// Create audit log
	auditLog := &models.AuditLog{
		TenantID:   tenantID,
		UserID:     &requesterUserID,
		Action:     string(models.AuditActionRevokeAPIKey),
		EntityType: string(models.EntityTypeUser),
		EntityID:   &targetUserID,
	}
	s.auditLogRepo.Create(ctx, auditLog)

	logger.Info("API key revoked successfully", "user_id", targetUserID, "tenant_id", tenantID)

	return nil
}

// generateTempPassword generates a secure temporary password for new users
func generateTempPassword() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a default password if random generation fails
		return "TempPass123!"
	}
	return base64.URLEncoding.EncodeToString(bytes)[:16]
}
