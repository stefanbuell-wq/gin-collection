package repositories

import (
	"context"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	// Create creates a new tenant
	Create(ctx context.Context, tenant *models.Tenant) error

	// GetByID retrieves a tenant by ID
	GetByID(ctx context.Context, id int64) (*models.Tenant, error)

	// GetBySubdomain retrieves a tenant by subdomain
	GetBySubdomain(ctx context.Context, subdomain string) (*models.Tenant, error)

	// GetByUUID retrieves a tenant by UUID
	GetByUUID(ctx context.Context, uuid string) (*models.Tenant, error)

	// Update updates a tenant
	Update(ctx context.Context, tenant *models.Tenant) error

	// UpdateStatus updates tenant status
	UpdateStatus(ctx context.Context, id int64, status models.TenantStatus) error

	// UpdateTier updates tenant subscription tier
	UpdateTier(ctx context.Context, id int64, tier models.SubscriptionTier) error
}
