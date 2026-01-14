package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// CocktailRepository implements cocktail data access
type CocktailRepository struct {
	db *sql.DB
}

// NewCocktailRepository creates a new cocktail repository
func NewCocktailRepository(db *sql.DB) *CocktailRepository {
	return &CocktailRepository{db: db}
}

// GetAll retrieves all cocktails (shared reference data)
func (r *CocktailRepository) GetAll(ctx context.Context) ([]*models.Cocktail, error) {
	query := `
		SELECT id, name, description, instructions, glass_type, ice_type, difficulty, prep_time, created_at
		FROM cocktails
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cocktails: %w", err)
	}
	defer rows.Close()

	var cocktails []*models.Cocktail

	for rows.Next() {
		cocktail := &models.Cocktail{}
		var description, instructions, glassType, iceType sql.NullString
		var prepTime sql.NullInt64

		err := rows.Scan(
			&cocktail.ID,
			&cocktail.Name,
			&description,
			&instructions,
			&glassType,
			&iceType,
			&cocktail.Difficulty,
			&prepTime,
			&cocktail.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan cocktail: %w", err)
		}

		if description.Valid {
			cocktail.Description = &description.String
		}
		if instructions.Valid {
			cocktail.Instructions = &instructions.String
		}
		if glassType.Valid {
			cocktail.GlassType = &glassType.String
		}
		if iceType.Valid {
			cocktail.IceType = &iceType.String
		}
		if prepTime.Valid {
			pt := int(prepTime.Int64)
			cocktail.PrepTime = &pt
		}

		cocktails = append(cocktails, cocktail)
	}

	return cocktails, nil
}

// GetByID retrieves a cocktail by ID with ingredients
func (r *CocktailRepository) GetByID(ctx context.Context, id int64) (*models.Cocktail, error) {
	query := `
		SELECT id, name, description, instructions, glass_type, ice_type, difficulty, prep_time, created_at
		FROM cocktails
		WHERE id = ?
	`

	cocktail := &models.Cocktail{}
	var description, instructions, glassType, iceType sql.NullString
	var prepTime sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cocktail.ID,
		&cocktail.Name,
		&description,
		&instructions,
		&glassType,
		&iceType,
		&cocktail.Difficulty,
		&prepTime,
		&cocktail.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cocktail not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cocktail: %w", err)
	}

	if description.Valid {
		cocktail.Description = &description.String
	}
	if instructions.Valid {
		cocktail.Instructions = &instructions.String
	}
	if glassType.Valid {
		cocktail.GlassType = &glassType.String
	}
	if iceType.Valid {
		cocktail.IceType = &iceType.String
	}
	if prepTime.Valid {
		pt := int(prepTime.Int64)
		cocktail.PrepTime = &pt
	}

	// Get ingredients
	ingredients, err := r.GetIngredientsForCocktail(ctx, id)
	if err != nil {
		return nil, err
	}

	cocktail.Ingredients = ingredients

	return cocktail, nil
}

// GetIngredientsForCocktail retrieves ingredients for a cocktail
func (r *CocktailRepository) GetIngredientsForCocktail(ctx context.Context, cocktailID int64) ([]*models.CocktailIngredient, error) {
	query := `
		SELECT id, cocktail_id, ingredient, amount, unit, is_gin
		FROM cocktail_ingredients
		WHERE cocktail_id = ?
		ORDER BY is_gin DESC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, cocktailID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cocktail ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []*models.CocktailIngredient

	for rows.Next() {
		ingredient := &models.CocktailIngredient{}
		var amount, unit sql.NullString

		err := rows.Scan(
			&ingredient.ID,
			&ingredient.CocktailID,
			&ingredient.Ingredient,
			&amount,
			&unit,
			&ingredient.IsGin,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}

		if amount.Valid {
			ingredient.Amount = &amount.String
		}
		if unit.Valid {
			ingredient.Unit = &unit.String
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

// GetCocktailsForGin retrieves cocktails that use a specific gin
func (r *CocktailRepository) GetCocktailsForGin(ctx context.Context, tenantID, ginID int64) ([]*models.Cocktail, error) {
	query := `
		SELECT
			c.id, c.name, c.description, c.instructions,
			c.glass_type, c.ice_type, c.difficulty, c.prep_time, c.created_at
		FROM cocktails c
		INNER JOIN gin_cocktails gc ON c.id = gc.cocktail_id
		WHERE gc.tenant_id = ? AND gc.gin_id = ?
		ORDER BY c.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gin cocktails: %w", err)
	}
	defer rows.Close()

	var cocktails []*models.Cocktail

	for rows.Next() {
		cocktail := &models.Cocktail{}
		var description, instructions, glassType, iceType sql.NullString
		var prepTime sql.NullInt64

		err := rows.Scan(
			&cocktail.ID,
			&cocktail.Name,
			&description,
			&instructions,
			&glassType,
			&iceType,
			&cocktail.Difficulty,
			&prepTime,
			&cocktail.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan cocktail: %w", err)
		}

		if description.Valid {
			cocktail.Description = &description.String
		}
		if instructions.Valid {
			cocktail.Instructions = &instructions.String
		}
		if glassType.Valid {
			cocktail.GlassType = &glassType.String
		}
		if iceType.Valid {
			cocktail.IceType = &iceType.String
		}
		if prepTime.Valid {
			pt := int(prepTime.Int64)
			cocktail.PrepTime = &pt
		}

		// Get ingredients
		ingredients, _ := r.GetIngredientsForCocktail(ctx, cocktail.ID)
		cocktail.Ingredients = ingredients

		cocktails = append(cocktails, cocktail)
	}

	return cocktails, nil
}

// LinkCocktailToGin links a cocktail to a gin
func (r *CocktailRepository) LinkCocktailToGin(ctx context.Context, tenantID, ginID, cocktailID int64) error {
	query := `
		INSERT IGNORE INTO gin_cocktails (tenant_id, gin_id, cocktail_id)
		VALUES (?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, ginID, cocktailID)
	if err != nil {
		return fmt.Errorf("failed to link cocktail to gin: %w", err)
	}

	return nil
}

// UnlinkCocktailFromGin unlinks a cocktail from a gin
func (r *CocktailRepository) UnlinkCocktailFromGin(ctx context.Context, tenantID, ginID, cocktailID int64) error {
	query := `
		DELETE FROM gin_cocktails
		WHERE tenant_id = ? AND gin_id = ? AND cocktail_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, tenantID, ginID, cocktailID)
	if err != nil {
		return fmt.Errorf("failed to unlink cocktail from gin: %w", err)
	}

	return nil
}

// Create creates a new cocktail (admin only)
func (r *CocktailRepository) Create(ctx context.Context, cocktail *models.Cocktail) error {
	query := `
		INSERT INTO cocktails (name, description, instructions, glass_type, ice_type, difficulty, prep_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		cocktail.Name,
		cocktail.Description,
		cocktail.Instructions,
		cocktail.GlassType,
		cocktail.IceType,
		cocktail.Difficulty,
		cocktail.PrepTime,
	)

	if err != nil {
		return fmt.Errorf("failed to create cocktail: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get cocktail ID: %w", err)
	}

	cocktail.ID = id

	return nil
}

// Update updates a cocktail (admin only)
func (r *CocktailRepository) Update(ctx context.Context, cocktail *models.Cocktail) error {
	query := `
		UPDATE cocktails
		SET name = ?, description = ?, instructions = ?,
		    glass_type = ?, ice_type = ?, difficulty = ?, prep_time = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		cocktail.Name,
		cocktail.Description,
		cocktail.Instructions,
		cocktail.GlassType,
		cocktail.IceType,
		cocktail.Difficulty,
		cocktail.PrepTime,
		cocktail.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update cocktail: %w", err)
	}

	return nil
}

// Delete deletes a cocktail (admin only)
func (r *CocktailRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM cocktails
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cocktail: %w", err)
	}

	return nil
}
