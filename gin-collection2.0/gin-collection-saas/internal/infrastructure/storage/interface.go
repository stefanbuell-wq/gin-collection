package storage

import "context"

// Storage defines the interface for photo storage
type Storage interface {
	UploadPhoto(ctx context.Context, tenantID int64, ginID int64, filename string, data []byte, contentType string) (*UploadResult, error)
	DeletePhoto(ctx context.Context, key string) error
	DownloadPhoto(ctx context.Context, key string) ([]byte, error)
	CheckExists(ctx context.Context, key string) (bool, error)
}
