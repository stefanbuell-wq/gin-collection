package gin

import (
	"context"
	"fmt"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// SimilarGinsRequest represents a request for similar gin suggestions
type SimilarGinsRequest struct {
	GinID   int64  `json:"gin_id"`
	Limit   int    `json:"limit"`
	Method  string `json:"method"` // "country", "type", "rating", "botanicals", "auto"
}

// SimilarGin represents a similar gin with match score
type SimilarGin struct {
	Gin        *models.Gin `json:"gin"`
	MatchScore float64     `json:"match_score"` // 0.0 - 1.0
	Reasons    []string    `json:"reasons"`     // Why it matches
}

// GetSimilarGins retrieves similar gins based on various criteria
func (s *Service) GetSimilarGins(ctx context.Context, tenantID, ginID int64, limit int) ([]*SimilarGin, error) {
	logger.Info("Getting similar gins", "tenant_id", tenantID, "gin_id", ginID, "limit", limit)

	// Get the source gin
	sourceGin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source gin: %w", err)
	}

	// Get all gins for this tenant
	filter := &models.GinFilter{
		TenantID:  tenantID,
		Limit:     0, // Get all
	}

	allGins, err := s.ginRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list gins: %w", err)
	}

	// Calculate similarity scores
	var similarGins []*SimilarGin

	for _, gin := range allGins {
		// Skip the source gin itself
		if gin.ID == sourceGin.ID {
			continue
		}

		// Calculate match score and reasons
		score, reasons := s.calculateSimilarity(sourceGin, gin)

		if score > 0 {
			similarGins = append(similarGins, &SimilarGin{
				Gin:        gin,
				MatchScore: score,
				Reasons:    reasons,
			})
		}
	}

	// Sort by match score (descending)
	for i := 0; i < len(similarGins); i++ {
		for j := i + 1; j < len(similarGins); j++ {
			if similarGins[j].MatchScore > similarGins[i].MatchScore {
				similarGins[i], similarGins[j] = similarGins[j], similarGins[i]
			}
		}
	}

	// Limit results
	if limit > 0 && len(similarGins) > limit {
		similarGins = similarGins[:limit]
	}

	logger.Info("Found similar gins", "count", len(similarGins))

	return similarGins, nil
}

// calculateSimilarity calculates similarity score between two gins
func (s *Service) calculateSimilarity(source, candidate *models.Gin) (float64, []string) {
	var score float64
	var reasons []string

	// Country match (30% weight)
	if source.Country != nil && candidate.Country != nil && *source.Country == *candidate.Country {
		score += 0.3
		reasons = append(reasons, fmt.Sprintf("Same country (%s)", *source.Country))
	}

	// Gin type match (25% weight)
	if source.GinType != nil && candidate.GinType != nil && *source.GinType == *candidate.GinType {
		score += 0.25
		reasons = append(reasons, fmt.Sprintf("Same type (%s)", *source.GinType))
	}

	// Rating similarity (20% weight)
	if source.Rating != nil && candidate.Rating != nil {
		ratingDiff := float64(abs(*source.Rating - *candidate.Rating))
		ratingScore := (5.0 - ratingDiff) / 5.0 * 0.2
		if ratingScore > 0.15 {
			score += ratingScore
			reasons = append(reasons, fmt.Sprintf("Similar rating (%d vs %d)", *source.Rating, *candidate.Rating))
		}
	}

	// ABV similarity (15% weight)
	if source.ABV != nil && candidate.ABV != nil {
		abvDiff := abs64(*source.ABV - *candidate.ABV)
		if abvDiff < 5.0 {
			abvScore := (5.0 - abvDiff) / 5.0 * 0.15
			score += abvScore
			if abvDiff < 2.0 {
				reasons = append(reasons, fmt.Sprintf("Similar ABV (%.1f%% vs %.1f%%)", *source.ABV, *candidate.ABV))
			}
		}
	}

	// Brand match (10% weight)
	if source.Brand != nil && candidate.Brand != nil && *source.Brand == *candidate.Brand {
		score += 0.1
		reasons = append(reasons, fmt.Sprintf("Same brand (%s)", *source.Brand))
	}

	return score, reasons
}

// GetSuggestionsByCountry retrieves gins from the same country
func (s *Service) GetSuggestionsByCountry(ctx context.Context, tenantID, ginID int64, limit int) ([]*models.Gin, error) {
	// Get the source gin
	sourceGin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source gin: %w", err)
	}

	if sourceGin.Country == nil {
		return []*models.Gin{}, nil
	}

	// Get gins from same country
	filter := &models.GinFilter{
		TenantID:  tenantID,
		Country:   sourceGin.Country,
		Limit:     limit,
		SortBy:    "rating",
		SortOrder: "desc",
	}

	gins, err := s.ginRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list gins: %w", err)
	}

	// Remove source gin from results
	var filtered []*models.Gin
	for _, gin := range gins {
		if gin.ID != sourceGin.ID {
			filtered = append(filtered, gin)
		}
	}

	return filtered, nil
}

// GetSuggestionsByType retrieves gins of the same type
func (s *Service) GetSuggestionsByType(ctx context.Context, tenantID, ginID int64, limit int) ([]*models.Gin, error) {
	// Get the source gin
	sourceGin, err := s.ginRepo.GetByID(ctx, tenantID, ginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source gin: %w", err)
	}

	if sourceGin.GinType == nil {
		return []*models.Gin{}, nil
	}

	// Get gins of same type
	filter := &models.GinFilter{
		TenantID:  tenantID,
		GinType:   sourceGin.GinType,
		Limit:     limit,
		SortBy:    "rating",
		SortOrder: "desc",
	}

	gins, err := s.ginRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list gins: %w", err)
	}

	// Remove source gin from results
	var filtered []*models.Gin
	for _, gin := range gins {
		if gin.ID != sourceGin.ID {
			filtered = append(filtered, gin)
		}
	}

	return filtered, nil
}

// abs returns absolute value of int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// abs64 returns absolute value of float64
func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
