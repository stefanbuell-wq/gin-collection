package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PhotoRepository defines photo data access
type PhotoRepository interface {
	// Create creates a new photo record
	Create(ctx context.Context, photo *models.GinPhoto) error

	// GetByID retrieves a photo by ID
	GetByID(ctx context.Context, tenantID, id int64) (*models.GinPhoto, error)

	// GetByGinID retrieves all photos for a specific gin
	GetByGinID(ctx context.Context, tenantID, ginID int64) ([]*models.GinPhoto, error)

	// Update updates a photo record
	Update(ctx context.Context, photo *models.GinPhoto) error

	// Delete deletes a photo record
	Delete(ctx context.Context, tenantID, id int64) error

	// SetPrimary sets a photo as primary (and unsets all others for that gin)
	SetPrimary(ctx context.Context, tenantID, ginID, photoID int64) error

	// CountByGinID counts photos for a gin
	CountByGinID(ctx context.Context, tenantID, ginID int64) (int, error)

	// GetTotalStorageUsage gets total storage usage for a tenant in KB
	GetTotalStorageUsage(ctx context.Context, tenantID int64) (int, error)
}
