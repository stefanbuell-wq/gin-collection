package gin

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/domain/models"
)

// ExportJSON exports gins as JSON
func (s *Service) ExportJSON(ctx context.Context, tenantID int64) ([]byte, error) {
	// Get all gins for tenant
	filter := &models.GinFilter{
		TenantID: tenantID,
		Limit:    0, // No limit for export
	}

	gins, err := s.ginRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list gins for export: %w", err)
	}

	// Convert to JSON
	data, err := json.MarshalIndent(gins, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gins to JSON: %w", err)
	}

	return data, nil
}

// ExportCSV exports gins as CSV
func (s *Service) ExportCSV(ctx context.Context, tenantID int64) (string, error) {
	// Get all gins for tenant
	filter := &models.GinFilter{
		TenantID: tenantID,
		Limit:    0, // No limit for export
	}

	gins, err := s.ginRepo.List(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("failed to list gins for export: %w", err)
	}

	// Build CSV
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write header
	header := []string{
		"ID", "Name", "Brand", "Country", "Region", "Gin Type", "ABV", "Bottle Size (ml)",
		"Fill Level (%)", "Price", "Market Value", "Purchase Date", "Purchase Location",
		"Barcode", "Rating", "Nose Notes", "Palate Notes", "Finish Notes", "General Notes",
		"Description", "Recommended Tonic", "Recommended Garnish", "Is Finished", "Created At",
	}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, gin := range gins {
		row := []string{
			strconv.FormatInt(gin.ID, 10),
			gin.Name,
			ptrToString(gin.Brand),
			ptrToString(gin.Country),
			ptrToString(gin.Region),
			ptrToString(gin.GinType),
			ptrFloat64ToString(gin.ABV),
			ptrIntToString(gin.BottleSize),
			ptrIntToString(gin.FillLevel),
			ptrFloat64ToString(gin.Price),
			ptrFloat64ToString(gin.CurrentMarketValue),
			ptrTimeToString(gin.PurchaseDate),
			ptrToString(gin.PurchaseLocation),
			ptrToString(gin.Barcode),
			ptrIntToString(gin.Rating),
			ptrToString(gin.NoseNotes),
			ptrToString(gin.PalateNotes),
			ptrToString(gin.FinishNotes),
			ptrToString(gin.GeneralNotes),
			ptrToString(gin.Description),
			ptrToString(gin.RecommendedTonic),
			ptrToString(gin.RecommendedGarnish),
			boolToString(gin.IsFinished),
			gin.CreatedAt.Format("2006-01-02"),
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return builder.String(), nil
}

// ImportJSON imports gins from JSON
func (s *Service) ImportJSON(ctx context.Context, tenantID int64, data []byte) (int, error) {
	var gins []*models.Gin
	if err := json.Unmarshal(data, &gins); err != nil {
		return 0, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	imported := 0
	for _, gin := range gins {
		// Set tenant ID
		gin.TenantID = tenantID
		gin.ID = 0 // Reset ID for new creation

		// Skip if barcode already exists
		if gin.Barcode != nil && *gin.Barcode != "" {
			exists, _ := s.ginRepo.CheckBarcodeExists(ctx, tenantID, *gin.Barcode)
			if exists {
				continue // Skip duplicate
			}
		}

		// Create gin
		if err := s.ginRepo.Create(ctx, gin); err != nil {
			// Log error but continue with next gin
			continue
		}

		imported++
	}

	// Update usage metrics
	if imported > 0 {
		if err := s.usageRepo.IncrementMetric(ctx, tenantID, "gin_count", imported); err != nil {
			// Log error but don't fail
		}
	}

	return imported, nil
}

// Helper functions
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrIntToString(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}

func ptrFloat64ToString(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *f)
}

func ptrTimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func boolToString(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
