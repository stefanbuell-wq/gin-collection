package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PasswordHistoryRepository defines the interface for password history data access
type PasswordHistoryRepository interface {
	// Add adds a password hash to the user's history
	Add(ctx context.Context, userID int64, passwordHash string) error

	// GetByUserID retrieves the password history for a user (most recent first)
	GetByUserID(ctx context.Context, userID int64, limit int) ([]*models.PasswordHistory, error)

	// IsPasswordUsed checks if a password hash exists in the user's history
	IsPasswordUsed(ctx context.Context, userID int64, passwordHash string) (bool, error)

	// Cleanup removes old password history entries beyond the limit
	Cleanup(ctx context.Context, userID int64, keepCount int) error
}
