package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// TastingSessionRepository handles database operations for tasting sessions
type TastingSessionRepository struct {
	db *sql.DB
}

// NewTastingSessionRepository creates a new tasting session repository
func NewTastingSessionRepository(db *sql.DB) *TastingSessionRepository {
	return &TastingSessionRepository{db: db}
}

// Create creates a new tasting session
func (r *TastingSessionRepository) Create(ctx context.Context, session *models.TastingSession) error {
	query := `
		INSERT INTO tasting_sessions (tenant_id, gin_id, user_id, date, notes, rating)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		session.TenantID,
		session.GinID,
		session.UserID,
		session.Date,
		session.Notes,
		session.Rating,
	)
	if err != nil {
		return fmt.Errorf("failed to create tasting session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	session.ID = id
	session.CreatedAt = time.Now()

	return nil
}

// GetByID retrieves a tasting session by ID
func (r *TastingSessionRepository) GetByID(ctx context.Context, tenantID, id int64) (*models.TastingSession, error) {
	query := `
		SELECT id, tenant_id, gin_id, user_id, date, notes, rating, created_at
		FROM tasting_sessions
		WHERE tenant_id = ? AND id = ?
	`

	session := &models.TastingSession{}
	err := r.db.QueryRowContext(ctx, query, tenantID, id).Scan(
		&session.ID,
		&session.TenantID,
		&session.GinID,
		&session.UserID,
		&session.Date,
		&session.Notes,
		&session.Rating,
		&session.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tasting session: %w", err)
	}

	return session, nil
}

// GetByGinID retrieves all tasting sessions for a specific gin
func (r *TastingSessionRepository) GetByGinID(ctx context.Context, tenantID, ginID int64) ([]*models.TastingSession, error) {
	query := `
		SELECT ts.id, ts.tenant_id, ts.gin_id, ts.user_id, ts.date, ts.notes, ts.rating, ts.created_at,
		       u.first_name, u.last_name
		FROM tasting_sessions ts
		LEFT JOIN users u ON ts.user_id = u.id
		WHERE ts.tenant_id = ? AND ts.gin_id = ?
		ORDER BY ts.date DESC, ts.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasting sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.TastingSession
	for rows.Next() {
		session := &models.TastingSession{}
		var firstName, lastName sql.NullString

		err := rows.Scan(
			&session.ID,
			&session.TenantID,
			&session.GinID,
			&session.UserID,
			&session.Date,
			&session.Notes,
			&session.Rating,
			&session.CreatedAt,
			&firstName,
			&lastName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tasting session: %w", err)
		}

		// Set user name if available
		if firstName.Valid || lastName.Valid {
			name := ""
			if firstName.Valid {
				name = firstName.String
			}
			if lastName.Valid {
				if name != "" {
					name += " "
				}
				name += lastName.String
			}
			session.UserName = &name
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// Update updates a tasting session
func (r *TastingSessionRepository) Update(ctx context.Context, session *models.TastingSession) error {
	query := `
		UPDATE tasting_sessions
		SET date = ?, notes = ?, rating = ?
		WHERE tenant_id = ? AND id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		session.Date,
		session.Notes,
		session.Rating,
		session.TenantID,
		session.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update tasting session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("tasting session not found")
	}

	return nil
}

// Delete deletes a tasting session
func (r *TastingSessionRepository) Delete(ctx context.Context, tenantID, id int64) error {
	query := `DELETE FROM tasting_sessions WHERE tenant_id = ? AND id = ?`

	result, err := r.db.ExecContext(ctx, query, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to delete tasting session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("tasting session not found")
	}

	return nil
}

// GetRecentByTenant retrieves recent tasting sessions for a tenant
func (r *TastingSessionRepository) GetRecentByTenant(ctx context.Context, tenantID int64, limit int) ([]*models.TastingSessionWithGin, error) {
	query := `
		SELECT ts.id, ts.tenant_id, ts.gin_id, ts.user_id, ts.date, ts.notes, ts.rating, ts.created_at,
		       g.name as gin_name, g.brand as gin_brand,
		       u.first_name, u.last_name
		FROM tasting_sessions ts
		JOIN gins g ON ts.gin_id = g.id
		LEFT JOIN users u ON ts.user_id = u.id
		WHERE ts.tenant_id = ?
		ORDER BY ts.date DESC, ts.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent tasting sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.TastingSessionWithGin
	for rows.Next() {
		session := &models.TastingSessionWithGin{}
		var ginBrand, firstName, lastName sql.NullString

		err := rows.Scan(
			&session.ID,
			&session.TenantID,
			&session.GinID,
			&session.UserID,
			&session.Date,
			&session.Notes,
			&session.Rating,
			&session.CreatedAt,
			&session.GinName,
			&ginBrand,
			&firstName,
			&lastName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tasting session: %w", err)
		}

		if ginBrand.Valid {
			session.GinBrand = &ginBrand.String
		}

		if firstName.Valid || lastName.Valid {
			name := ""
			if firstName.Valid {
				name = firstName.String
			}
			if lastName.Valid {
				if name != "" {
					name += " "
				}
				name += lastName.String
			}
			session.UserName = &name
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
