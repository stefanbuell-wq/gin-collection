package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PasswordHistoryRepository implements the password history repository interface
type PasswordHistoryRepository struct {
	db *sql.DB
}

// NewPasswordHistoryRepository creates a new password history repository
func NewPasswordHistoryRepository(db *sql.DB) *PasswordHistoryRepository {
	return &PasswordHistoryRepository{db: db}
}

// Add adds a password hash to the user's history
func (r *PasswordHistoryRepository) Add(ctx context.Context, userID int64, passwordHash string) error {
	query := `
		INSERT INTO password_history (user_id, password_hash, created_at)
		VALUES (?, ?, NOW())
	`

	_, err := r.db.ExecContext(ctx, query, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to add password to history: %w", err)
	}

	return nil
}

// GetByUserID retrieves the password history for a user (most recent first)
func (r *PasswordHistoryRepository) GetByUserID(ctx context.Context, userID int64, limit int) ([]*models.PasswordHistory, error) {
	query := `
		SELECT id, user_id, password_hash, created_at
		FROM password_history
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get password history: %w", err)
	}
	defer rows.Close()

	var history []*models.PasswordHistory
	for rows.Next() {
		h := &models.PasswordHistory{}
		err := rows.Scan(&h.ID, &h.UserID, &h.PasswordHash, &h.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan password history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

// IsPasswordUsed checks if a password hash exists in the user's recent history
func (r *PasswordHistoryRepository) IsPasswordUsed(ctx context.Context, userID int64, passwordHash string) (bool, error) {
	// We need to check against the actual password, not the hash directly
	// because bcrypt generates different hashes for the same password
	// This method is called after getting the history and comparing with bcrypt
	query := `
		SELECT COUNT(*)
		FROM password_history
		WHERE user_id = ? AND password_hash = ?
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID, passwordHash).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check password history: %w", err)
	}

	return count > 0, nil
}

// Cleanup removes old password history entries beyond the limit
func (r *PasswordHistoryRepository) Cleanup(ctx context.Context, userID int64, keepCount int) error {
	// Delete entries older than the most recent 'keepCount' entries
	query := `
		DELETE FROM password_history
		WHERE user_id = ? AND id NOT IN (
			SELECT id FROM (
				SELECT id FROM password_history
				WHERE user_id = ?
				ORDER BY created_at DESC
				LIMIT ?
			) AS recent
		)
	`

	_, err := r.db.ExecContext(ctx, query, userID, userID, keepCount)
	if err != nil {
		return fmt.Errorf("failed to cleanup password history: %w", err)
	}

	return nil
}
