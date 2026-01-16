package tasting

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/repository/mysql"
)

// Service handles tasting session business logic
type Service struct {
	tastingRepo *mysql.TastingSessionRepository
	ginRepo     *mysql.GinRepository
}

// NewService creates a new tasting service
func NewService(tastingRepo *mysql.TastingSessionRepository, ginRepo *mysql.GinRepository) *Service {
	return &Service{
		tastingRepo: tastingRepo,
		ginRepo:     ginRepo,
	}
}

// CreateSession creates a new tasting session
func (s *Service) CreateSession(ctx context.Context, session *models.TastingSession) error {
	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, session.TenantID, session.GinID)
	if err != nil {
		return fmt.Errorf("failed to verify gin: %w", err)
	}
	if gin == nil {
		return fmt.Errorf("gin not found")
	}

	// Set default date if not provided
	if session.Date.IsZero() {
		session.Date = time.Now()
	}

	// Validate rating
	if session.Rating != nil && (*session.Rating < 1 || *session.Rating > 5) {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	return s.tastingRepo.Create(ctx, session)
}

// GetSession retrieves a tasting session by ID
func (s *Service) GetSession(ctx context.Context, tenantID, id int64) (*models.TastingSession, error) {
	session, err := s.tastingRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasting session: %w", err)
	}
	return session, nil
}

// GetSessionsForGin retrieves all tasting sessions for a specific gin
func (s *Service) GetSessionsForGin(ctx context.Context, tenantID, ginID int64) ([]*models.TastingSession, error) {
	// Verify gin exists and belongs to tenant
	gin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify gin: %w", err)
	}
	if gin == nil {
		return nil, fmt.Errorf("gin not found")
	}

	sessions, err := s.tastingRepo.GetByGinID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasting sessions: %w", err)
	}

	return sessions, nil
}

// UpdateSession updates a tasting session
func (s *Service) UpdateSession(ctx context.Context, session *models.TastingSession) error {
	// Verify session exists
	existing, err := s.tastingRepo.GetByID(ctx, session.TenantID, session.ID)
	if err != nil {
		return fmt.Errorf("failed to verify tasting session: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("tasting session not found")
	}

	// Validate rating
	if session.Rating != nil && (*session.Rating < 1 || *session.Rating > 5) {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	return s.tastingRepo.Update(ctx, session)
}

// DeleteSession deletes a tasting session
func (s *Service) DeleteSession(ctx context.Context, tenantID, id int64) error {
	return s.tastingRepo.Delete(ctx, tenantID, id)
}

// GetRecentSessions retrieves recent tasting sessions for a tenant
func (s *Service) GetRecentSessions(ctx context.Context, tenantID int64, limit int) ([]*models.TastingSessionWithGin, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	sessions, err := s.tastingRepo.GetRecentByTenant(ctx, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tasting sessions: %w", err)
	}

	return sessions, nil
}
