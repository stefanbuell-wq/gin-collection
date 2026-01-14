package cocktail

import (
	"context"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Service handles cocktail business logic
type Service struct {
	cocktailRepo repositories.CocktailRepository
	ginRepo      repositories.GinRepository
}

// NewService creates a new cocktail service
func NewService(
	cocktailRepo repositories.CocktailRepository,
	ginRepo repositories.GinRepository,
) *Service {
	return &Service{
		cocktailRepo: cocktailRepo,
		ginRepo:      ginRepo,
	}
}

// GetAllCocktails retrieves all available cocktails
func (s *Service) GetAllCocktails(ctx context.Context) ([]*models.Cocktail, error) {
	cocktails, err := s.cocktailRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cocktails: %w", err)
	}

	// Get ingredients for each cocktail
	for _, cocktail := range cocktails {
		ingredients, err := s.cocktailRepo.GetIngredientsForCocktail(ctx, cocktail.ID)
		if err != nil {
			logger.Error("Failed to get ingredients for cocktail", "cocktail_id", cocktail.ID, "error", err.Error())
			continue
		}
		cocktail.Ingredients = ingredients
	}

	return cocktails, nil
}

// GetCocktailByID retrieves a cocktail by ID with ingredients
func (s *Service) GetCocktailByID(ctx context.Context, id int64) (*models.Cocktail, error) {
	cocktail, err := s.cocktailRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cocktail: %w", err)
	}

	return cocktail, nil
}

// GetCocktailsForGin retrieves cocktails for a specific gin
func (s *Service) GetCocktailsForGin(ctx context.Context, tenantID, ginID int64) ([]*models.Cocktail, error) {
	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return nil, fmt.Errorf("gin does not belong to tenant")
	}

	cocktails, err := s.cocktailRepo.GetCocktailsForGin(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gin cocktails: %w", err)
	}

	return cocktails, nil
}

// LinkCocktailToGin links a cocktail to a gin
func (s *Service) LinkCocktailToGin(ctx context.Context, tenantID, ginID, cocktailID int64) error {
	logger.Info("Linking cocktail to gin", "tenant_id", tenantID, "gin_id", ginID, "cocktail_id", cocktailID)

	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return fmt.Errorf("gin does not belong to tenant")
	}

	// Verify cocktail exists
	_, err = s.cocktailRepo.GetByID(ctx, cocktailID)
	if err != nil {
		return fmt.Errorf("cocktail not found: %w", err)
	}

	// Link cocktail to gin
	if err := s.cocktailRepo.LinkCocktailToGin(ctx, tenantID, ginID, cocktailID); err != nil {
		logger.Error("Failed to link cocktail to gin", "error", err.Error())
		return fmt.Errorf("failed to link cocktail: %w", err)
	}

	logger.Info("Cocktail linked to gin successfully", "gin_id", ginID, "cocktail_id", cocktailID)

	return nil
}

// UnlinkCocktailFromGin unlinks a cocktail from a gin
func (s *Service) UnlinkCocktailFromGin(ctx context.Context, tenantID, ginID, cocktailID int64) error {
	logger.Info("Unlinking cocktail from gin", "tenant_id", tenantID, "gin_id", ginID, "cocktail_id", cocktailID)

	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return fmt.Errorf("gin does not belong to tenant")
	}

	// Unlink cocktail from gin
	if err := s.cocktailRepo.UnlinkCocktailFromGin(ctx, tenantID, ginID, cocktailID); err != nil {
		logger.Error("Failed to unlink cocktail from gin", "error", err.Error())
		return fmt.Errorf("failed to unlink cocktail: %w", err)
	}

	logger.Info("Cocktail unlinked from gin successfully", "gin_id", ginID, "cocktail_id", cocktailID)

	return nil
}
