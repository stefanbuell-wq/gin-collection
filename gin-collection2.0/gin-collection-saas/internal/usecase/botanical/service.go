package botanical

import (
	"context"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Service handles botanical business logic
type Service struct {
	botanicalRepo repositories.BotanicalRepository
	ginRepo       repositories.GinRepository
}

// NewService creates a new botanical service
func NewService(
	botanicalRepo repositories.BotanicalRepository,
	ginRepo repositories.GinRepository,
) *Service {
	return &Service{
		botanicalRepo: botanicalRepo,
		ginRepo:       ginRepo,
	}
}

// GetAllBotanicals retrieves all available botanicals
func (s *Service) GetAllBotanicals(ctx context.Context) ([]*models.Botanical, error) {
	botanicals, err := s.botanicalRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get botanicals: %w", err)
	}

	return botanicals, nil
}

// GetGinBotanicals retrieves botanicals for a specific gin
func (s *Service) GetGinBotanicals(ctx context.Context, tenantID, ginID int64) ([]*models.GinBotanical, error) {
	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return nil, fmt.Errorf("gin does not belong to tenant")
	}

	botanicals, err := s.botanicalRepo.GetByGinID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gin botanicals: %w", err)
	}

	return botanicals, nil
}

// UpdateGinBotanicals updates the botanicals for a gin
func (s *Service) UpdateGinBotanicals(ctx context.Context, tenantID, ginID int64, botanicals []*models.GinBotanical) error {
	logger.Info("Updating gin botanicals", "tenant_id", tenantID, "gin_id", ginID, "count", len(botanicals))

	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return fmt.Errorf("gin does not belong to tenant")
	}

	// Validate botanicals
	for _, botanical := range botanicals {
		if botanical.BotanicalID == 0 {
			return fmt.Errorf("botanical_id is required")
		}

		// Verify botanical exists
		_, err := s.botanicalRepo.GetByID(ctx, botanical.BotanicalID)
		if err != nil {
			return fmt.Errorf("invalid botanical_id %d: %w", botanical.BotanicalID, err)
		}

		// Validate prominence
		if botanical.Prominence != models.ProminenceDominant &&
			botanical.Prominence != models.ProminenceNotable &&
			botanical.Prominence != models.ProminenceSubtle {
			return fmt.Errorf("invalid prominence: %s (must be dominant, notable, or subtle)", botanical.Prominence)
		}

		// Set gin_id
		botanical.GinID = ginID
	}

	// Update botanicals
	if err := s.botanicalRepo.UpdateGinBotanicals(ctx, tenantID, ginID, botanicals); err != nil {
		logger.Error("Failed to update gin botanicals", "error", err.Error())
		return fmt.Errorf("failed to update botanicals: %w", err)
	}

	logger.Info("Gin botanicals updated successfully", "gin_id", ginID, "count", len(botanicals))

	return nil
}
