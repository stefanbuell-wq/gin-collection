package photo

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/storage"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// Service handles photo business logic
type Service struct {
	photoRepo        repositories.PhotoRepository
	ginRepo          repositories.GinRepository
	usageMetricsRepo repositories.UsageMetricsRepository
	tenantRepo       repositories.TenantRepository
	storage          storage.Storage
}

// NewService creates a new photo service
func NewService(
	photoRepo repositories.PhotoRepository,
	ginRepo repositories.GinRepository,
	usageMetricsRepo repositories.UsageMetricsRepository,
	tenantRepo repositories.TenantRepository,
	storageClient storage.Storage,
) *Service {
	return &Service{
		photoRepo:        photoRepo,
		ginRepo:          ginRepo,
		usageMetricsRepo: usageMetricsRepo,
		tenantRepo:       tenantRepo,
		storage:          storageClient,
	}
}

// UploadPhoto uploads a photo for a gin
func (s *Service) UploadPhoto(ctx context.Context, tenantID, ginID int64, filename string, data []byte, photoType models.PhotoType, caption *string) (*models.GinPhoto, error) {
	logger.Info("Uploading photo", "tenant_id", tenantID, "gin_id", ginID, "filename", filename, "size", len(data))

	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return nil, fmt.Errorf("gin does not belong to tenant")
	}

	// Get tenant for limits
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	limits := tenant.GetLimits()

	// Check photo count limit
	currentCount, err := s.photoRepo.CountByGinID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to count photos: %w", err)
	}

	if currentCount >= limits.MaxPhotosPerGin {
		return nil, fmt.Errorf("photo limit reached (%d/%d) - upgrade to add more photos", currentCount, limits.MaxPhotosPerGin)
	}

	// Check storage limit (if applicable)
	fileSizeKB := len(data) / 1024
	if limits.StorageLimitMB != nil {
		currentStorageMB, err := s.usageMetricsRepo.GetMetric(ctx, tenantID, "storage_mb")
		if err != nil {
			return nil, fmt.Errorf("failed to get storage usage: %w", err)
		}

		fileSizeMB := fileSizeKB / 1024
		if currentStorageMB+fileSizeMB > *limits.StorageLimitMB {
			return nil, fmt.Errorf("storage limit would be exceeded (%dMB/%dMB)", currentStorageMB+fileSizeMB, *limits.StorageLimitMB)
		}
	}

	// Validate file type
	contentType := detectContentType(filename)
	if !isValidImageType(contentType) {
		return nil, fmt.Errorf("invalid file type: %s (allowed: jpg, jpeg, png, webp)", contentType)
	}

	// Upload to S3
	uploadResult, err := s.storage.UploadPhoto(ctx, tenantID, ginID, filename, data, contentType)
	if err != nil {
		logger.Error("Failed to upload to S3", "error", err.Error())
		return nil, fmt.Errorf("failed to upload photo: %w", err)
	}

	// Check if this is the first photo (make it primary)
	isPrimary := currentCount == 0

	// Create photo record
	photo := &models.GinPhoto{
		TenantID:   tenantID,
		GinID:      ginID,
		PhotoURL:   uploadResult.URL,
		PhotoType:  photoType,
		Caption:    caption,
		IsPrimary:  isPrimary,
		StorageKey: &uploadResult.Key,
		FileSizeKB: &fileSizeKB,
	}

	if err := s.photoRepo.Create(ctx, photo); err != nil {
		// Rollback: delete from S3
		s.storage.DeletePhoto(ctx, uploadResult.Key)
		return nil, fmt.Errorf("failed to create photo record: %w", err)
	}

	// Update storage metrics
	storageMB := fileSizeKB / 1024
	if storageMB > 0 {
		if err := s.usageMetricsRepo.IncrementMetric(ctx, tenantID, "storage_mb", storageMB); err != nil {
			logger.Error("Failed to update storage metrics", "error", err.Error())
		}
	}

	logger.Info("Photo uploaded successfully", "photo_id", photo.ID, "gin_id", ginID)

	return photo, nil
}

// GetPhotosByGinID retrieves all photos for a gin
func (s *Service) GetPhotosByGinID(ctx context.Context, tenantID, ginID int64) ([]*models.GinPhoto, error) {
	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return nil, fmt.Errorf("gin does not belong to tenant")
	}

	photos, err := s.photoRepo.GetByGinID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}

	return photos, nil
}

// DeletePhoto deletes a photo
func (s *Service) DeletePhoto(ctx context.Context, tenantID, photoID int64) error {
	logger.Info("Deleting photo", "tenant_id", tenantID, "photo_id", photoID)

	// Get photo
	photo, err := s.photoRepo.GetByID(ctx, tenantID, photoID)
	if err != nil {
		return fmt.Errorf("photo not found: %w", err)
	}

	// Delete from S3
	if photo.StorageKey != nil {
		if err := s.storage.DeletePhoto(ctx, *photo.StorageKey); err != nil {
			logger.Error("Failed to delete from S3", "error", err.Error())
			// Continue anyway to delete DB record
		}
	}

	// Delete from database
	if err := s.photoRepo.Delete(ctx, tenantID, photoID); err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	// Update storage metrics
	if photo.FileSizeKB != nil {
		storageMB := *photo.FileSizeKB / 1024
		if storageMB > 0 {
			if err := s.usageMetricsRepo.DecrementMetric(ctx, tenantID, "storage_mb", storageMB); err != nil {
				logger.Error("Failed to update storage metrics", "error", err.Error())
			}
		}
	}

	logger.Info("Photo deleted successfully", "photo_id", photoID)

	return nil
}

// SetPrimaryPhoto sets a photo as primary
func (s *Service) SetPrimaryPhoto(ctx context.Context, tenantID, ginID, photoID int64) error {
	logger.Info("Setting primary photo", "tenant_id", tenantID, "gin_id", ginID, "photo_id", photoID)

	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return fmt.Errorf("gin not found: %w", err)
	}

	if gin.TenantID != tenantID {
		return fmt.Errorf("gin does not belong to tenant")
	}

	// Set primary
	if err := s.photoRepo.SetPrimary(ctx, tenantID, ginID, photoID); err != nil {
		return fmt.Errorf("failed to set primary photo: %w", err)
	}

	logger.Info("Primary photo set successfully", "gin_id", ginID, "photo_id", photoID)

	return nil
}

// detectContentType detects content type from filename
func detectContentType(filename string) string {
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(filename), "."))

	// Extract extension
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		ext = strings.ToLower(parts[len(parts)-1])
	}

	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	case "gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// isValidImageType checks if content type is a valid image
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/gif",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}

	return false
}
