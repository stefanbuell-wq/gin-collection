package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PasswordResetRepository defines the interface for password reset token data access
type PasswordResetRepository interface {
	// Create creates a new password reset token
	Create(ctx context.Context, token *models.PasswordResetToken) error

	// GetByToken retrieves a password reset token by token string
	GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error)

	// MarkAsUsed marks a token as used
	MarkAsUsed(ctx context.Context, id int64) error

	// DeleteExpired deletes all expired tokens
	DeleteExpired(ctx context.Context) error

	// DeleteByUserID deletes all tokens for a user
	DeleteByUserID(ctx context.Context, userID int64) error
}
