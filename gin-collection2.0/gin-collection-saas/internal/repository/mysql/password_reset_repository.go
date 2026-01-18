package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PasswordResetRepository implements the password reset repository interface
type PasswordResetRepository struct {
	db *sql.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	token.ID = id
	return nil
}

// GetByToken retrieves a password reset token by token string
func (r *PasswordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token = ?
	`

	resetToken := &models.PasswordResetToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.UsedAt,
		&resetToken.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	return resetToken, nil
}

// MarkAsUsed marks a token as used
func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, id int64) error {
	query := `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// DeleteExpired deletes all expired tokens
func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE expires_at < NOW()
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all tokens for a user
func (r *PasswordResetRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM password_reset_tokens
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens for user: %w", err)
	}

	return nil
}

// generateToken generates a secure random token (helper, not part of interface)
func generateToken() string {
	// This is handled in the service layer
	return ""
}

// TokenExpiry is the duration for which a token is valid
const TokenExpiry = 1 * time.Hour
