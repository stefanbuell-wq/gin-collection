package repositories

import (
	"context"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// GetByEmail retrieves a user by email within a tenant
	GetByEmail(ctx context.Context, tenantID int64, email string) (*models.User, error)

	// GetByEmailGlobal retrieves a user by email across all tenants (for login without subdomain)
	GetByEmailGlobal(ctx context.Context, email string) (*models.User, error)

	// GetByAPIKey retrieves a user by API key
	GetByAPIKey(ctx context.Context, apiKey string) (*models.User, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// UpdateLastLogin updates the last login timestamp
	UpdateLastLogin(ctx context.Context, id int64) error

	// List lists all users for a tenant
	List(ctx context.Context, tenantID int64) ([]*models.User, error)

	// Delete deletes a user
	Delete(ctx context.Context, id int64) error

	// GenerateAPIKey generates a new API key for a user (Enterprise only)
	GenerateAPIKey(ctx context.Context, userID int64) (string, error)

	// RevokeAPIKey revokes the API key for a user
	RevokeAPIKey(ctx context.Context, userID int64) error

	// CountByTenant counts users in a tenant
	CountByTenant(ctx context.Context, tenantID int64) (int, error)
}
