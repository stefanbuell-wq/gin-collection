package gin

import (
	"context"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Service handles gin business logic
type Service struct {
	ginRepo   repositories.GinRepository
	usageRepo repositories.UsageMetricsRepository
}

// NewService creates a new gin service
func NewService(
	ginRepo repositories.GinRepository,
	usageRepo repositories.UsageMetricsRepository,
) *Service {
	return &Service{
		ginRepo:   ginRepo,
		usageRepo: usageRepo,
	}
}

// Create creates a new gin
func (s *Service) Create(ctx context.Context, gin *models.Gin) error {
	logger.Info("Creating gin", "tenant_id", gin.TenantID, "name", gin.Name)

	// Validate
	if gin.Name == "" {
		return errors.ErrInvalidInput
	}

	// Check rating if provided
	if gin.Rating != nil && (*gin.Rating < 1 || *gin.Rating > 5) {
		return errors.ErrInvalidRating
	}

	// Check if barcode already exists (if provided)
	if gin.Barcode != nil && *gin.Barcode != "" {
		exists, err := s.ginRepo.CheckBarcodeExists(ctx, gin.TenantID, *gin.Barcode)
		if err != nil {
			logger.Error("Failed to check barcode", "error", err.Error())
			return fmt.Errorf("failed to check barcode: %w", err)
		}
		if exists {
			return errors.ErrBarcodeAlreadyExists
		}
	}

	// Create gin
	if err := s.ginRepo.Create(ctx, gin); err != nil {
		logger.Error("Failed to create gin", "error", err.Error())
		return fmt.Errorf("failed to create gin: %w", err)
	}

	// Increment usage metric
	if err := s.usageRepo.IncrementMetric(ctx, gin.TenantID, "gin_count", 1); err != nil {
		logger.Error("Failed to increment gin count", "error", err.Error())
		// Don't fail the operation, just log
	}

	logger.Info("Gin created successfully", "gin_id", gin.ID, "tenant_id", gin.TenantID)
	return nil
}

// GetByID retrieves a gin by ID
func (s *Service) GetByID(ctx context.Context, tenantID, id int64) (*models.Gin, error) {
	gin, err := s.ginRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	return gin, nil
}

// List retrieves gins with filtering and pagination
func (s *Service) List(ctx context.Context, filter *models.GinFilter) ([]*models.Gin, error) {
	gins, err := s.ginRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list gins: %w", err)
	}

	return gins, nil
}

// Update updates a gin
func (s *Service) Update(ctx context.Context, gin *models.Gin) error {
	logger.Info("Updating gin", "gin_id", gin.ID, "tenant_id", gin.TenantID)

	// Validate
	if gin.Name == "" {
		return errors.ErrInvalidInput
	}

	// Check rating if provided
	if gin.Rating != nil && (*gin.Rating < 1 || *gin.Rating > 5) {
		return errors.ErrInvalidRating
	}

	// Update gin
	if err := s.ginRepo.Update(ctx, gin); err != nil {
		logger.Error("Failed to update gin", "error", err.Error())
		return err
	}

	logger.Info("Gin updated successfully", "gin_id", gin.ID)
	return nil
}

// Delete deletes a gin
func (s *Service) Delete(ctx context.Context, tenantID, id int64) error {
	logger.Info("Deleting gin", "gin_id", id, "tenant_id", tenantID)

	// Delete gin
	if err := s.ginRepo.Delete(ctx, tenantID, id); err != nil {
		logger.Error("Failed to delete gin", "error", err.Error())
		return err
	}

	// Decrement usage metric
	if err := s.usageRepo.DecrementMetric(ctx, tenantID, "gin_count", 1); err != nil {
		logger.Error("Failed to decrement gin count", "error", err.Error())
		// Don't fail the operation, just log
	}

	logger.Info("Gin deleted successfully", "gin_id", id)
	return nil
}

// Search searches gins by query string
func (s *Service) Search(ctx context.Context, tenantID int64, query string, limit, offset int) ([]*models.Gin, error) {
	if limit == 0 {
		limit = 20 // Default limit
	}

	gins, err := s.ginRepo.Search(ctx, tenantID, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search gins: %w", err)
	}

	return gins, nil
}

// GetStats retrieves statistics for a tenant's gin collection
func (s *Service) GetStats(ctx context.Context, tenantID int64) (*models.GinStats, error) {
	stats, err := s.ginRepo.GetStats(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}
