package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// BotanicalRepository defines botanical data access
type BotanicalRepository interface {
	// GetAll retrieves all botanicals (shared reference data)
	GetAll(ctx context.Context) ([]*models.Botanical, error)

	// GetByID retrieves a botanical by ID
	GetByID(ctx context.Context, id int64) (*models.Botanical, error)

	// GetByGinID retrieves all botanicals for a specific gin
	GetByGinID(ctx context.Context, tenantID, ginID int64) ([]*models.GinBotanical, error)

	// UpdateGinBotanicals updates botanicals for a gin (delete all + insert new)
	UpdateGinBotanicals(ctx context.Context, tenantID, ginID int64, botanicals []*models.GinBotanical) error

	// Create creates a new botanical (admin only)
	Create(ctx context.Context, botanical *models.Botanical) error

	// Update updates a botanical (admin only)
	Update(ctx context.Context, botanical *models.Botanical) error

	// Delete deletes a botanical (admin only)
	Delete(ctx context.Context, id int64) error
}
