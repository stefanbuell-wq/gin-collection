package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// GinReferenceRepository defines operations for gin reference catalog
type GinReferenceRepository interface {
	// Search searches gin references by name, brand, or other fields
	Search(ctx context.Context, params *models.GinReferenceSearchParams) ([]*models.GinReference, int, error)

	// GetByID retrieves a single gin reference by ID
	GetByID(ctx context.Context, id int64) (*models.GinReference, error)

	// GetByBarcode retrieves a gin reference by barcode
	GetByBarcode(ctx context.Context, barcode string) (*models.GinReference, error)

	// GetCountries returns list of unique countries
	GetCountries(ctx context.Context) ([]string, error)

	// GetGinTypes returns list of unique gin types
	GetGinTypes(ctx context.Context) ([]string, error)

	// GetBrands returns list of unique brands
	GetBrands(ctx context.Context) ([]string, error)
}
