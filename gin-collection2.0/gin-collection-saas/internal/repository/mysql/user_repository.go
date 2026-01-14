package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// UserRepository implements the user repository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (tenant_id, uuid, email, password_hash, first_name, last_name,
		                   role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		user.TenantID,
		uuid.New().String(),
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, tenant_id, uuid, email, password_hash, first_name, last_name,
		       role, api_key, is_active, email_verified_at, last_login_at,
		       created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.TenantID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.APIKey,
		&user.IsActive,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email within a tenant
func (r *UserRepository) GetByEmail(ctx context.Context, tenantID int64, email string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, uuid, email, password_hash, first_name, last_name,
		       role, api_key, is_active, email_verified_at, last_login_at,
		       created_at, updated_at
		FROM users
		WHERE tenant_id = ? AND email = ?
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, tenantID, email).Scan(
		&user.ID,
		&user.TenantID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.APIKey,
		&user.IsActive,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetByAPIKey retrieves a user by API key
func (r *UserRepository) GetByAPIKey(ctx context.Context, apiKey string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, uuid, email, password_hash, first_name, last_name,
		       role, api_key, is_active, email_verified_at, last_login_at,
		       created_at, updated_at
		FROM users
		WHERE api_key = ? AND is_active = TRUE
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, apiKey).Scan(
		&user.ID,
		&user.TenantID,
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.APIKey,
		&user.IsActive,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrInvalidAPIKey
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by API key: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, first_name = ?, last_name = ?, role = ?,
		    is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	query := `UPDATE users SET last_login_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// List lists all users for a tenant
func (r *UserRepository) List(ctx context.Context, tenantID int64) ([]*models.User, error) {
	query := `
		SELECT id, tenant_id, uuid, email, password_hash, first_name, last_name,
		       role, api_key, is_active, email_verified_at, last_login_at,
		       created_at, updated_at
		FROM users
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.TenantID,
			&user.UUID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.APIKey,
			&user.IsActive,
			&user.EmailVerifiedAt,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GenerateAPIKey generates a new API key for a user (Enterprise only)
func (r *UserRepository) GenerateAPIKey(ctx context.Context, userID int64) (string, error) {
	// Generate a secure API key using UUID
	apiKey := "sk_" + uuid.New().String()

	query := `UPDATE users SET api_key = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, apiKey, userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}

	return apiKey, nil
}

// RevokeAPIKey revokes the API key for a user
func (r *UserRepository) RevokeAPIKey(ctx context.Context, userID int64) error {
	query := `UPDATE users SET api_key = NULL, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	return nil
}

// CountByTenant counts users in a tenant
func (r *UserRepository) CountByTenant(ctx context.Context, tenantID int64) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE tenant_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
