package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/domain/repositories"
)

type ginReferenceRepository struct {
	db *sql.DB
}

// NewGinReferenceRepository creates a new gin reference repository
func NewGinReferenceRepository(db *sql.DB) repositories.GinReferenceRepository {
	return &ginReferenceRepository{db: db}
}

// Search searches gin references by name, brand, or other fields
func (r *ginReferenceRepository) Search(ctx context.Context, params *models.GinReferenceSearchParams) ([]*models.GinReference, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE clause
	if params.Query != "" {
		searchTerm := "%" + strings.ToLower(params.Query) + "%"
		conditions = append(conditions, "(LOWER(name) LIKE ? OR LOWER(brand) LIKE ? OR LOWER(country) LIKE ?)")
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if params.Country != "" {
		conditions = append(conditions, "country = ?")
		args = append(args, params.Country)
	}

	if params.GinType != "" {
		conditions = append(conditions, "gin_type = ?")
		args = append(args, params.GinType)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM gin_references %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count gin references: %w", err)
	}

	// Apply pagination
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT id, name, brand, country, region, gin_type, abv, bottle_size,
		       description, nose_notes, palate_notes, finish_notes,
		       recommended_tonic, recommended_garnish, image_url, barcode
		FROM gin_references
		%s
		ORDER BY name ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search gin references: %w", err)
	}
	defer rows.Close()

	var gins []*models.GinReference
	for rows.Next() {
		g := &models.GinReference{}
		err := rows.Scan(
			&g.ID, &g.Name, &g.Brand, &g.Country, &g.Region, &g.GinType,
			&g.ABV, &g.BottleSize, &g.Description, &g.NoseNotes, &g.PalateNotes,
			&g.FinishNotes, &g.RecommendedTonic, &g.RecommendedGarnish,
			&g.ImageURL, &g.Barcode,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan gin reference: %w", err)
		}
		gins = append(gins, g)
	}

	return gins, total, nil
}

// GetByID retrieves a single gin reference by ID
func (r *ginReferenceRepository) GetByID(ctx context.Context, id int64) (*models.GinReference, error) {
	query := `
		SELECT id, name, brand, country, region, gin_type, abv, bottle_size,
		       description, nose_notes, palate_notes, finish_notes,
		       recommended_tonic, recommended_garnish, image_url, barcode
		FROM gin_references
		WHERE id = ?
	`

	g := &models.GinReference{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&g.ID, &g.Name, &g.Brand, &g.Country, &g.Region, &g.GinType,
		&g.ABV, &g.BottleSize, &g.Description, &g.NoseNotes, &g.PalateNotes,
		&g.FinishNotes, &g.RecommendedTonic, &g.RecommendedGarnish,
		&g.ImageURL, &g.Barcode,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get gin reference: %w", err)
	}

	return g, nil
}

// GetCountries returns list of unique countries
func (r *ginReferenceRepository) GetCountries(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT country FROM gin_references WHERE country IS NOT NULL ORDER BY country`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get countries: %w", err)
	}
	defer rows.Close()

	var countries []string
	for rows.Next() {
		var country string
		if err := rows.Scan(&country); err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, country)
	}

	return countries, nil
}

// GetGinTypes returns list of unique gin types
func (r *ginReferenceRepository) GetGinTypes(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT gin_type FROM gin_references WHERE gin_type IS NOT NULL ORDER BY gin_type`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get gin types: %w", err)
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var ginType string
		if err := rows.Scan(&ginType); err != nil {
			return nil, fmt.Errorf("failed to scan gin type: %w", err)
		}
		types = append(types, ginType)
	}

	return types, nil
}

// GetBrands returns list of unique brands
func (r *ginReferenceRepository) GetBrands(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT brand FROM gin_references WHERE brand IS NOT NULL ORDER BY brand`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get brands: %w", err)
	}
	defer rows.Close()

	var brands []string
	for rows.Next() {
		var brand string
		if err := rows.Scan(&brand); err != nil {
			return nil, fmt.Errorf("failed to scan brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// Ensure argIndex is used (for future parameterized queries)
var _ = func() int { return 1 }()
