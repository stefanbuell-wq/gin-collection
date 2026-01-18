package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainErrors "github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// PhotoRepository implements photo data access
type PhotoRepository struct {
	db *sql.DB
}

// NewPhotoRepository creates a new photo repository
func NewPhotoRepository(db *sql.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

// Create creates a new photo record
func (r *PhotoRepository) Create(ctx context.Context, photo *models.GinPhoto) error {
	query := `
		INSERT INTO gin_photos (
			tenant_id, gin_id, photo_url, photo_type, caption,
			is_primary, storage_key, file_size_kb, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		photo.TenantID,
		photo.GinID,
		photo.PhotoURL,
		photo.PhotoType,
		photo.Caption,
		photo.IsPrimary,
		photo.StorageKey,
		photo.FileSizeKB,
	)

	if err != nil {
		return fmt.Errorf("failed to create photo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get photo ID: %w", err)
	}

	photo.ID = id
	photo.CreatedAt = time.Now()

	return nil
}

// GetByID retrieves a photo by ID
func (r *PhotoRepository) GetByID(ctx context.Context, tenantID, id int64) (*models.GinPhoto, error) {
	query := `
		SELECT id, tenant_id, gin_id, photo_url, photo_type, caption,
		       is_primary, storage_key, file_size_kb, created_at
		FROM gin_photos
		WHERE tenant_id = ? AND id = ?
	`

	photo := &models.GinPhoto{}
	var caption, storageKey sql.NullString
	var fileSizeKB sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, tenantID, id).Scan(
		&photo.ID,
		&photo.TenantID,
		&photo.GinID,
		&photo.PhotoURL,
		&photo.PhotoType,
		&caption,
		&photo.IsPrimary,
		&storageKey,
		&fileSizeKB,
		&photo.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domainErrors.ErrPhotoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get photo: %w", err)
	}

	if caption.Valid {
		photo.Caption = &caption.String
	}
	if storageKey.Valid {
		photo.StorageKey = &storageKey.String
	}
	if fileSizeKB.Valid {
		size := int(fileSizeKB.Int64)
		photo.FileSizeKB = &size
	}

	return photo, nil
}

// GetByGinID retrieves all photos for a specific gin
func (r *PhotoRepository) GetByGinID(ctx context.Context, tenantID, ginID int64) ([]*models.GinPhoto, error) {
	query := `
		SELECT id, tenant_id, gin_id, photo_url, photo_type, caption,
		       is_primary, storage_key, file_size_kb, created_at
		FROM gin_photos
		WHERE tenant_id = ? AND gin_id = ?
		ORDER BY is_primary DESC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	defer rows.Close()

	var photos []*models.GinPhoto

	for rows.Next() {
		photo := &models.GinPhoto{}
		var caption, storageKey sql.NullString
		var fileSizeKB sql.NullInt64

		err := rows.Scan(
			&photo.ID,
			&photo.TenantID,
			&photo.GinID,
			&photo.PhotoURL,
			&photo.PhotoType,
			&caption,
			&photo.IsPrimary,
			&storageKey,
			&fileSizeKB,
			&photo.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}

		if caption.Valid {
			photo.Caption = &caption.String
		}
		if storageKey.Valid {
			photo.StorageKey = &storageKey.String
		}
		if fileSizeKB.Valid {
			size := int(fileSizeKB.Int64)
			photo.FileSizeKB = &size
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

// Update updates a photo record
func (r *PhotoRepository) Update(ctx context.Context, photo *models.GinPhoto) error {
	query := `
		UPDATE gin_photos
		SET photo_url = ?, photo_type = ?, caption = ?, is_primary = ?
		WHERE tenant_id = ? AND id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		photo.PhotoURL,
		photo.PhotoType,
		photo.Caption,
		photo.IsPrimary,
		photo.TenantID,
		photo.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}

	return nil
}

// Delete deletes a photo record
func (r *PhotoRepository) Delete(ctx context.Context, tenantID, id int64) error {
	query := `
		DELETE FROM gin_photos
		WHERE tenant_id = ? AND id = ?
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	return nil
}

// SetPrimary sets a photo as primary (and unsets all others for that gin)
func (r *PhotoRepository) SetPrimary(ctx context.Context, tenantID, ginID, photoID int64) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Unset all primary photos for this gin
	unsetQuery := `
		UPDATE gin_photos
		SET is_primary = FALSE
		WHERE tenant_id = ? AND gin_id = ?
	`

	_, err = tx.ExecContext(ctx, unsetQuery, tenantID, ginID)
	if err != nil {
		return fmt.Errorf("failed to unset primary photos: %w", err)
	}

	// Set the new primary photo
	setQuery := `
		UPDATE gin_photos
		SET is_primary = TRUE
		WHERE tenant_id = ? AND id = ? AND gin_id = ?
	`

	result, err := tx.ExecContext(ctx, setQuery, tenantID, photoID, ginID)
	if err != nil {
		return fmt.Errorf("failed to set primary photo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domainErrors.ErrPhotoNotFound
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CountByGinID counts photos for a gin
func (r *PhotoRepository) CountByGinID(ctx context.Context, tenantID, ginID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM gin_photos
		WHERE tenant_id = ? AND gin_id = ?
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID, ginID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count photos: %w", err)
	}

	return count, nil
}

// GetTotalStorageUsage gets total storage usage for a tenant in KB
func (r *PhotoRepository) GetTotalStorageUsage(ctx context.Context, tenantID int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(file_size_kb), 0)
		FROM gin_photos
		WHERE tenant_id = ?
	`

	var totalKB int
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&totalKB)
	if err != nil {
		return 0, fmt.Errorf("failed to get total storage usage: %w", err)
	}

	return totalKB, nil
}
