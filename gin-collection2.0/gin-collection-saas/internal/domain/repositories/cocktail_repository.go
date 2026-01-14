package repositories

import (
	"context"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// CocktailRepository defines cocktail data access
type CocktailRepository interface {
	// GetAll retrieves all cocktails (shared reference data)
	GetAll(ctx context.Context) ([]*models.Cocktail, error)

	// GetByID retrieves a cocktail by ID with ingredients
	GetByID(ctx context.Context, id int64) (*models.Cocktail, error)

	// GetIngredientsForCocktail retrieves ingredients for a cocktail
	GetIngredientsForCocktail(ctx context.Context, cocktailID int64) ([]*models.CocktailIngredient, error)

	// GetCocktailsForGin retrieves cocktails that use a specific gin
	GetCocktailsForGin(ctx context.Context, tenantID, ginID int64) ([]*models.Cocktail, error)

	// LinkCocktailToGin links a cocktail to a gin
	LinkCocktailToGin(ctx context.Context, tenantID, ginID, cocktailID int64) error

	// UnlinkCocktailFromGin unlinks a cocktail from a gin
	UnlinkCocktailFromGin(ctx context.Context, tenantID, ginID, cocktailID int64) error

	// Create creates a new cocktail (admin only)
	Create(ctx context.Context, cocktail *models.Cocktail) error

	// Update updates a cocktail (admin only)
	Update(ctx context.Context, cocktail *models.Cocktail) error

	// Delete deletes a cocktail (admin only)
	Delete(ctx context.Context, id int64) error
}
