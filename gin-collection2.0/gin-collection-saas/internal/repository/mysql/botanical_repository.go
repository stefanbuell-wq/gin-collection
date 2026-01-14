package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// BotanicalRepository implements botanical data access
type BotanicalRepository struct {
	db *sql.DB
}

// NewBotanicalRepository creates a new botanical repository
func NewBotanicalRepository(db *sql.DB) *BotanicalRepository {
	return &BotanicalRepository{db: db}
}

// GetAll retrieves all botanicals (shared reference data)
func (r *BotanicalRepository) GetAll(ctx context.Context) ([]*models.Botanical, error) {
	query := `
		SELECT id, name, category, description
		FROM botanicals
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get botanicals: %w", err)
	}
	defer rows.Close()

	var botanicals []*models.Botanical

	for rows.Next() {
		botanical := &models.Botanical{}
		var description sql.NullString

		err := rows.Scan(
			&botanical.ID,
			&botanical.Name,
			&botanical.Category,
			&description,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan botanical: %w", err)
		}

		if description.Valid {
			botanical.Description = &description.String
		}

		botanicals = append(botanicals, botanical)
	}

	return botanicals, nil
}

// GetByID retrieves a botanical by ID
func (r *BotanicalRepository) GetByID(ctx context.Context, id int64) (*models.Botanical, error) {
	query := `
		SELECT id, name, category, description
		FROM botanicals
		WHERE id = ?
	`

	botanical := &models.Botanical{}
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&botanical.ID,
		&botanical.Name,
		&botanical.Category,
		&description,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("botanical not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get botanical: %w", err)
	}

	if description.Valid {
		botanical.Description = &description.String
	}

	return botanical, nil
}

// GetByGinID retrieves all botanicals for a specific gin
func (r *BotanicalRepository) GetByGinID(ctx context.Context, tenantID, ginID int64) ([]*models.GinBotanical, error) {
	query := `
		SELECT
			gb.id, gb.gin_id, gb.botanical_id, gb.prominence,
			b.name, b.category, b.description
		FROM gin_botanicals gb
		INNER JOIN botanicals b ON gb.botanical_id = b.id
		WHERE gb.tenant_id = ? AND gb.gin_id = ?
		ORDER BY
			CASE gb.prominence
				WHEN 'dominant' THEN 1
				WHEN 'notable' THEN 2
				WHEN 'subtle' THEN 3
				ELSE 4
			END,
			b.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gin botanicals: %w", err)
	}
	defer rows.Close()

	var ginBotanicals []*models.GinBotanical

	for rows.Next() {
		gb := &models.GinBotanical{
			Botanical: &models.Botanical{},
		}
		var description sql.NullString

		err := rows.Scan(
			&gb.ID,
			&gb.GinID,
			&gb.BotanicalID,
			&gb.Prominence,
			&gb.Botanical.Name,
			&gb.Botanical.Category,
			&description,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan gin botanical: %w", err)
		}

		gb.Botanical.ID = gb.BotanicalID

		if description.Valid {
			gb.Botanical.Description = &description.String
		}

		ginBotanicals = append(ginBotanicals, gb)
	}

	return ginBotanicals, nil
}

// UpdateGinBotanicals updates botanicals for a gin (delete all + insert new)
func (r *BotanicalRepository) UpdateGinBotanicals(ctx context.Context, tenantID, ginID int64, botanicals []*models.GinBotanical) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing botanicals for this gin
	deleteQuery := `
		DELETE FROM gin_botanicals
		WHERE tenant_id = ? AND gin_id = ?
	`

	_, err = tx.ExecContext(ctx, deleteQuery, tenantID, ginID)
	if err != nil {
		return fmt.Errorf("failed to delete existing botanicals: %w", err)
	}

	// Insert new botanicals
	if len(botanicals) > 0 {
		insertQuery := `
			INSERT INTO gin_botanicals (tenant_id, gin_id, botanical_id, prominence)
			VALUES (?, ?, ?, ?)
		`

		for _, botanical := range botanicals {
			_, err := tx.ExecContext(ctx, insertQuery,
				tenantID,
				ginID,
				botanical.BotanicalID,
				botanical.Prominence,
			)
			if err != nil {
				return fmt.Errorf("failed to insert botanical: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Create creates a new botanical (admin only)
func (r *BotanicalRepository) Create(ctx context.Context, botanical *models.Botanical) error {
	query := `
		INSERT INTO botanicals (name, category, description)
		VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		botanical.Name,
		botanical.Category,
		botanical.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to create botanical: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get botanical ID: %w", err)
	}

	botanical.ID = id

	return nil
}

// Update updates a botanical (admin only)
func (r *BotanicalRepository) Update(ctx context.Context, botanical *models.Botanical) error {
	query := `
		UPDATE botanicals
		SET name = ?, category = ?, description = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		botanical.Name,
		botanical.Category,
		botanical.Description,
		botanical.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update botanical: %w", err)
	}

	return nil
}

// Delete deletes a botanical (admin only)
func (r *BotanicalRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM botanicals
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete botanical: %w", err)
	}

	return nil
}
