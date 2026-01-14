package repositories

import (
	"context"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// GinRepository defines the interface for gin data access
type GinRepository interface {
	// Create creates a new gin
	Create(ctx context.Context, gin *models.Gin) error

	// GetByID retrieves a gin by ID (with tenant scoping)
	GetByID(ctx context.Context, tenantID, id int64) (*models.Gin, error)

	// GetByUUID retrieves a gin by UUID (with tenant scoping)
	GetByUUID(ctx context.Context, tenantID int64, uuid string) (*models.Gin, error)

	// List retrieves gins with filtering and pagination
	List(ctx context.Context, filter *models.GinFilter) ([]*models.Gin, error)

	// Update updates a gin
	Update(ctx context.Context, gin *models.Gin) error

	// Delete deletes a gin
	Delete(ctx context.Context, tenantID, id int64) error

	// Count counts gins for a tenant with optional filters
	Count(ctx context.Context, tenantID int64, isFinished *bool) (int, error)

	// Search searches gins by query string
	Search(ctx context.Context, tenantID int64, query string, limit, offset int) ([]*models.Gin, error)

	// GetStats retrieves statistics for a tenant's gin collection
	GetStats(ctx context.Context, tenantID int64) (*models.GinStats, error)

	// CheckBarcodeExists checks if a barcode already exists for a tenant
	CheckBarcodeExists(ctx context.Context, tenantID int64, barcode string) (bool, error)
}

// UsageMetricsRepository defines the interface for usage metrics
type UsageMetricsRepository interface {
	// GetMetric retrieves current value for a metric
	GetMetric(ctx context.Context, tenantID int64, metricName string) (int, error)

	// IncrementMetric increments a metric by delta
	IncrementMetric(ctx context.Context, tenantID int64, metricName string, delta int) error

	// DecrementMetric decrements a metric by delta
	DecrementMetric(ctx context.Context, tenantID int64, metricName string, delta int) error

	// SetMetric sets a metric to a specific value
	SetMetric(ctx context.Context, tenantID int64, metricName string, value int) error
}
