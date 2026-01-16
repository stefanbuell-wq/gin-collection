package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// LocalStorage handles local file storage as S3 fallback
type LocalStorage struct {
	basePath string
	baseURL  string
}

// LocalStorageConfig holds local storage configuration
type LocalStorageConfig struct {
	BasePath string // e.g., "/app/uploads"
	BaseURL  string // e.g., "https://ginvault.cloud/uploads"
}

// NewLocalStorage creates a new local storage client
func NewLocalStorage(cfg *LocalStorageConfig) (*LocalStorage, error) {
	// Create base directory if not exists
	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	logger.Info("Local storage initialized", "path", cfg.BasePath, "url", cfg.BaseURL)

	return &LocalStorage{
		basePath: cfg.BasePath,
		baseURL:  cfg.BaseURL,
	}, nil
}

// UploadPhoto uploads a photo to local storage
func (s *LocalStorage) UploadPhoto(ctx context.Context, tenantID int64, ginID int64, filename string, data []byte, contentType string) (*UploadResult, error) {
	// Generate unique key
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("tenants/%d/gins/%d/%s%s", tenantID, ginID, uuid.New().String(), ext)

	// Create full path
	fullPath := filepath.Join(s.basePath, key)
	dir := filepath.Dir(fullPath)

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	logger.Info("Uploading photo to local storage", "key", key, "size", len(data))

	// Write file
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, key)

	logger.Info("Photo uploaded successfully to local storage", "key", key, "url", url)

	return &UploadResult{
		Key:       key,
		URL:       url,
		SizeBytes: int64(len(data)),
	}, nil
}

// DeletePhoto deletes a photo from local storage
func (s *LocalStorage) DeletePhoto(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.basePath, key)

	logger.Info("Deleting photo from local storage", "key", key)

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logger.Info("Photo deleted successfully from local storage", "key", key)

	return nil
}

// DownloadPhoto downloads a photo from local storage
func (s *LocalStorage) DownloadPhoto(ctx context.Context, key string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, key)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// CheckExists checks if a file exists
func (s *LocalStorage) CheckExists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(s.basePath, key)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
