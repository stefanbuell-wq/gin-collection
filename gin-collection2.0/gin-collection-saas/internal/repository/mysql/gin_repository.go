package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/internal/domain/errors"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// GinRepository implements the gin repository interface
type GinRepository struct {
	db *sql.DB
}

// NewGinRepository creates a new gin repository
func NewGinRepository(db *sql.DB) *GinRepository {
	return &GinRepository{db: db}
}

// Create creates a new gin
func (r *GinRepository) Create(ctx context.Context, gin *models.Gin) error {
	query := `
		INSERT INTO gins (
			tenant_id, uuid, name, brand, country, region, gin_type, abv,
			bottle_size, fill_level, price, current_market_value, purchase_date,
			purchase_location, barcode, rating, nose_notes, palate_notes,
			finish_notes, general_notes, description, photo_url, is_finished,
			recommended_tonic, recommended_garnish, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		gin.TenantID,
		uuid.New().String(),
		gin.Name,
		gin.Brand,
		gin.Country,
		gin.Region,
		gin.GinType,
		gin.ABV,
		gin.BottleSize,
		gin.FillLevel,
		gin.Price,
		gin.CurrentMarketValue,
		gin.PurchaseDate,
		gin.PurchaseLocation,
		gin.Barcode,
		gin.Rating,
		gin.NoseNotes,
		gin.PalateNotes,
		gin.FinishNotes,
		gin.GeneralNotes,
		gin.Description,
		gin.PhotoURL,
		gin.IsFinished,
		gin.RecommendedTonic,
		gin.RecommendedGarnish,
	)

	if err != nil {
		return fmt.Errorf("failed to create gin: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	gin.ID = id
	return nil
}

// GetByID retrieves a gin by ID with tenant scoping
func (r *GinRepository) GetByID(ctx context.Context, tenantID, id int64) (*models.Gin, error) {
	query := `
		SELECT id, tenant_id, uuid, name, brand, country, region, gin_type, abv,
		       bottle_size, fill_level, price, current_market_value, purchase_date,
		       purchase_location, barcode, rating, nose_notes, palate_notes,
		       finish_notes, general_notes, description, photo_url, is_finished,
		       recommended_tonic, recommended_garnish, created_at, updated_at
		FROM gins
		WHERE tenant_id = ? AND id = ?
	`

	gin := &models.Gin{}
	err := r.db.QueryRowContext(ctx, query, tenantID, id).Scan(
		&gin.ID,
		&gin.TenantID,
		&gin.UUID,
		&gin.Name,
		&gin.Brand,
		&gin.Country,
		&gin.Region,
		&gin.GinType,
		&gin.ABV,
		&gin.BottleSize,
		&gin.FillLevel,
		&gin.Price,
		&gin.CurrentMarketValue,
		&gin.PurchaseDate,
		&gin.PurchaseLocation,
		&gin.Barcode,
		&gin.Rating,
		&gin.NoseNotes,
		&gin.PalateNotes,
		&gin.FinishNotes,
		&gin.GeneralNotes,
		&gin.Description,
		&gin.PhotoURL,
		&gin.IsFinished,
		&gin.RecommendedTonic,
		&gin.RecommendedGarnish,
		&gin.CreatedAt,
		&gin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrGinNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get gin: %w", err)
	}

	return gin, nil
}

// GetByUUID retrieves a gin by UUID with tenant scoping
func (r *GinRepository) GetByUUID(ctx context.Context, tenantID int64, uuidStr string) (*models.Gin, error) {
	query := `
		SELECT id, tenant_id, uuid, name, brand, country, region, gin_type, abv,
		       bottle_size, fill_level, price, current_market_value, purchase_date,
		       purchase_location, barcode, rating, nose_notes, palate_notes,
		       finish_notes, general_notes, description, photo_url, is_finished,
		       recommended_tonic, recommended_garnish, created_at, updated_at
		FROM gins
		WHERE tenant_id = ? AND uuid = ?
	`

	gin := &models.Gin{}
	err := r.db.QueryRowContext(ctx, query, tenantID, uuidStr).Scan(
		&gin.ID,
		&gin.TenantID,
		&gin.UUID,
		&gin.Name,
		&gin.Brand,
		&gin.Country,
		&gin.Region,
		&gin.GinType,
		&gin.ABV,
		&gin.BottleSize,
		&gin.FillLevel,
		&gin.Price,
		&gin.CurrentMarketValue,
		&gin.PurchaseDate,
		&gin.PurchaseLocation,
		&gin.Barcode,
		&gin.Rating,
		&gin.NoseNotes,
		&gin.PalateNotes,
		&gin.FinishNotes,
		&gin.GeneralNotes,
		&gin.Description,
		&gin.PhotoURL,
		&gin.IsFinished,
		&gin.RecommendedTonic,
		&gin.RecommendedGarnish,
		&gin.CreatedAt,
		&gin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrGinNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get gin by UUID: %w", err)
	}

	return gin, nil
}

// List retrieves gins with filtering and pagination
func (r *GinRepository) List(ctx context.Context, filter *models.GinFilter) ([]*models.Gin, error) {
	query := `
		SELECT id, tenant_id, uuid, name, brand, country, region, gin_type, abv,
		       bottle_size, fill_level, price, current_market_value, purchase_date,
		       purchase_location, barcode, rating, nose_notes, palate_notes,
		       finish_notes, general_notes, description, photo_url, is_finished,
		       recommended_tonic, recommended_garnish, created_at, updated_at
		FROM gins
		WHERE tenant_id = ?
	`

	args := []interface{}{filter.TenantID}

	// Apply filters
	if filter.IsFinished != nil {
		query += " AND is_finished = ?"
		args = append(args, *filter.IsFinished)
	}

	if filter.GinType != nil && *filter.GinType != "" {
		query += " AND gin_type = ?"
		args = append(args, *filter.GinType)
	}

	if filter.Country != nil && *filter.Country != "" {
		query += " AND country = ?"
		args = append(args, *filter.Country)
	}

	if filter.MinRating != nil {
		query += " AND rating >= ?"
		args = append(args, *filter.MinRating)
	}

	if filter.MaxRating != nil {
		query += " AND rating <= ?"
		args = append(args, *filter.MaxRating)
	}

	// Sorting
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, strings.ToUpper(sortOrder))

	// Pagination
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list gins: %w", err)
	}
	defer rows.Close()

	var gins []*models.Gin
	for rows.Next() {
		gin := &models.Gin{}
		err := rows.Scan(
			&gin.ID,
			&gin.TenantID,
			&gin.UUID,
			&gin.Name,
			&gin.Brand,
			&gin.Country,
			&gin.Region,
			&gin.GinType,
			&gin.ABV,
			&gin.BottleSize,
			&gin.FillLevel,
			&gin.Price,
			&gin.CurrentMarketValue,
			&gin.PurchaseDate,
			&gin.PurchaseLocation,
			&gin.Barcode,
			&gin.Rating,
			&gin.NoseNotes,
			&gin.PalateNotes,
			&gin.FinishNotes,
			&gin.GeneralNotes,
			&gin.Description,
			&gin.PhotoURL,
			&gin.IsFinished,
			&gin.RecommendedTonic,
			&gin.RecommendedGarnish,
			&gin.CreatedAt,
			&gin.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gin: %w", err)
		}
		gins = append(gins, gin)
	}

	return gins, nil
}

// Update updates a gin
func (r *GinRepository) Update(ctx context.Context, gin *models.Gin) error {
	query := `
		UPDATE gins SET
			name = ?, brand = ?, country = ?, region = ?, gin_type = ?, abv = ?,
			bottle_size = ?, fill_level = ?, price = ?, current_market_value = ?,
			purchase_date = ?, purchase_location = ?, barcode = ?, rating = ?,
			nose_notes = ?, palate_notes = ?, finish_notes = ?, general_notes = ?,
			description = ?, photo_url = ?, is_finished = ?, recommended_tonic = ?,
			recommended_garnish = ?, updated_at = NOW()
		WHERE tenant_id = ? AND id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		gin.Name,
		gin.Brand,
		gin.Country,
		gin.Region,
		gin.GinType,
		gin.ABV,
		gin.BottleSize,
		gin.FillLevel,
		gin.Price,
		gin.CurrentMarketValue,
		gin.PurchaseDate,
		gin.PurchaseLocation,
		gin.Barcode,
		gin.Rating,
		gin.NoseNotes,
		gin.PalateNotes,
		gin.FinishNotes,
		gin.GeneralNotes,
		gin.Description,
		gin.PhotoURL,
		gin.IsFinished,
		gin.RecommendedTonic,
		gin.RecommendedGarnish,
		gin.TenantID,
		gin.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update gin: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrGinNotFound
	}

	return nil
}

// Delete deletes a gin
func (r *GinRepository) Delete(ctx context.Context, tenantID, id int64) error {
	query := `DELETE FROM gins WHERE tenant_id = ? AND id = ?`

	result, err := r.db.ExecContext(ctx, query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete gin: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrGinNotFound
	}

	return nil
}

// Count counts gins for a tenant with optional filters
func (r *GinRepository) Count(ctx context.Context, tenantID int64, isFinished *bool) (int, error) {
	query := "SELECT COUNT(*) FROM gins WHERE tenant_id = ?"
	args := []interface{}{tenantID}

	if isFinished != nil {
		query += " AND is_finished = ?"
		args = append(args, *isFinished)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count gins: %w", err)
	}

	return count, nil
}

// Search searches gins by query string (full-text search)
func (r *GinRepository) Search(ctx context.Context, tenantID int64, query string, limit, offset int) ([]*models.Gin, error) {
	searchQuery := `
		SELECT id, tenant_id, uuid, name, brand, country, region, gin_type, abv,
		       bottle_size, fill_level, price, current_market_value, purchase_date,
		       purchase_location, barcode, rating, nose_notes, palate_notes,
		       finish_notes, general_notes, description, photo_url, is_finished,
		       recommended_tonic, recommended_garnish, created_at, updated_at
		FROM gins
		WHERE tenant_id = ?
		AND (
			name LIKE ? OR
			brand LIKE ? OR
			country LIKE ? OR
			nose_notes LIKE ? OR
			palate_notes LIKE ? OR
			finish_notes LIKE ? OR
			general_notes LIKE ?
		)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	searchTerm := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, searchQuery,
		tenantID,
		searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		limit, offset,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to search gins: %w", err)
	}
	defer rows.Close()

	var gins []*models.Gin
	for rows.Next() {
		gin := &models.Gin{}
		err := rows.Scan(
			&gin.ID,
			&gin.TenantID,
			&gin.UUID,
			&gin.Name,
			&gin.Brand,
			&gin.Country,
			&gin.Region,
			&gin.GinType,
			&gin.ABV,
			&gin.BottleSize,
			&gin.FillLevel,
			&gin.Price,
			&gin.CurrentMarketValue,
			&gin.PurchaseDate,
			&gin.PurchaseLocation,
			&gin.Barcode,
			&gin.Rating,
			&gin.NoseNotes,
			&gin.PalateNotes,
			&gin.FinishNotes,
			&gin.GeneralNotes,
			&gin.Description,
			&gin.PhotoURL,
			&gin.IsFinished,
			&gin.RecommendedTonic,
			&gin.RecommendedGarnish,
			&gin.CreatedAt,
			&gin.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gin: %w", err)
		}
		gins = append(gins, gin)
	}

	return gins, nil
}

// GetStats retrieves statistics for a tenant's gin collection
func (r *GinRepository) GetStats(ctx context.Context, tenantID int64) (*models.GinStats, error) {
	stats := &models.GinStats{
		GinsByType:            make(map[string]int),
		GinsByCountry:         make(map[string]int),
		FillLevelDistribution: make(map[string]int),
	}

	// Total, Available, Finished counts
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN is_finished = 0 THEN 1 ELSE 0 END) as available,
			SUM(CASE WHEN is_finished = 1 THEN 1 ELSE 0 END) as finished
		FROM gins
		WHERE tenant_id = ?
	`, tenantID).Scan(&stats.TotalGins, &stats.AvailableGins, &stats.FinishedGins)

	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	// Average rating
	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(rating), 0)
		FROM gins
		WHERE tenant_id = ? AND rating IS NOT NULL
	`, tenantID).Scan(&stats.AverageRating)

	if err != nil {
		return nil, fmt.Errorf("failed to get average rating: %w", err)
	}

	// Total value and market value
	err = r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(price), 0),
			COALESCE(SUM(current_market_value), 0)
		FROM gins
		WHERE tenant_id = ?
	`, tenantID).Scan(&stats.TotalValue, &stats.TotalMarketValue)

	if err != nil {
		return nil, fmt.Errorf("failed to get value stats: %w", err)
	}

	// Gins by type
	rows, err := r.db.QueryContext(ctx, `
		SELECT gin_type, COUNT(*) as count
		FROM gins
		WHERE tenant_id = ? AND gin_type IS NOT NULL
		GROUP BY gin_type
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to get gins by type: %w", err)
	}

	for rows.Next() {
		var ginType string
		var count int
		if err := rows.Scan(&ginType, &count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan gin type: %w", err)
		}
		stats.GinsByType[ginType] = count
	}
	rows.Close()

	// Gins by country
	rows, err = r.db.QueryContext(ctx, `
		SELECT country, COUNT(*) as count
		FROM gins
		WHERE tenant_id = ? AND country IS NOT NULL
		GROUP BY country
		ORDER BY count DESC
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to get gins by country: %w", err)
	}

	for rows.Next() {
		var country string
		var count int
		if err := rows.Scan(&country, &count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		stats.GinsByCountry[country] = count
	}
	rows.Close()

	// Top rated gins
	topRatedRows, err := r.db.QueryContext(ctx, `
		SELECT id, name, brand, rating, photo_url
		FROM gins
		WHERE tenant_id = ? AND rating IS NOT NULL
		ORDER BY rating DESC, name ASC
		LIMIT 5
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to get top rated gins: %w", err)
	}
	defer topRatedRows.Close()

	for topRatedRows.Next() {
		gin := &models.Gin{}
		err := topRatedRows.Scan(&gin.ID, &gin.Name, &gin.Brand, &gin.Rating, &gin.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top rated gin: %w", err)
		}
		stats.TopRatedGins = append(stats.TopRatedGins, gin)
	}

	return stats, nil
}

// CheckBarcodeExists checks if a barcode already exists for a tenant
func (r *GinRepository) CheckBarcodeExists(ctx context.Context, tenantID int64, barcode string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM gins WHERE tenant_id = ? AND barcode = ?
	`, tenantID, barcode).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check barcode: %w", err)
	}

	return count > 0, nil
}
