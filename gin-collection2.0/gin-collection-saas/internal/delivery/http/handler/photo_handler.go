package handler

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/middleware"
	"github.com/yourusername/gin-collection-saas/internal/delivery/http/response"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	photoUsecase "github.com/yourusername/gin-collection-saas/internal/usecase/photo"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// PhotoHandler handles photo HTTP requests
type PhotoHandler struct {
	photoService *photoUsecase.Service
}

// NewPhotoHandler creates a new photo handler
func NewPhotoHandler(photoService *photoUsecase.Service) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

// GetPhotos handles GET /api/v1/gins/:gin_id/photos
func (h *PhotoHandler) GetPhotos(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	ginIDStr := c.Param("gin_id")
	ginID, err := strconv.ParseInt(ginIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	photos, err := h.photoService.GetPhotosByGinID(c.Request.Context(), tenantID, ginID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"photos": photos,
		"count":  len(photos),
	})
}

// Upload handles POST /api/v1/gins/:gin_id/photos
func (h *PhotoHandler) Upload(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	ginIDStr := c.Param("gin_id")
	ginID, err := strconv.ParseInt(ginIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		logger.Debug("Failed to get form file", "error", err.Error())
		response.ValidationError(c, map[string]string{
			"error": "Photo file is required",
		})
		return
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file", "error", err.Error())
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	// Get optional parameters
	photoTypeStr := c.PostForm("photo_type")
	if photoTypeStr == "" {
		photoTypeStr = "bottle"
	}

	photoType := models.PhotoType(photoTypeStr)

	// Validate photo type
	if photoType != models.PhotoTypeBottle &&
		photoType != models.PhotoTypeLabel &&
		photoType != models.PhotoTypeMoment &&
		photoType != models.PhotoTypeTasting {
		response.ValidationError(c, map[string]string{
			"error": "Invalid photo_type. Allowed: bottle, label, moment, tasting",
		})
		return
	}

	// Get caption (optional)
	var caption *string
	if captionStr := c.PostForm("caption"); captionStr != "" {
		caption = &captionStr
	}

	// Upload photo
	photo, err := h.photoService.UploadPhoto(
		c.Request.Context(),
		tenantID,
		ginID,
		header.Filename,
		data,
		photoType,
		caption,
	)

	if err != nil {
		logger.Error("Failed to upload photo", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Created(c, photo)
}

// Delete handles DELETE /api/v1/gins/:gin_id/photos/:photo_id
func (h *PhotoHandler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	photoIDStr := c.Param("photo_id")
	photoID, err := strconv.ParseInt(photoIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid photo ID"})
		return
	}

	// Delete photo
	if err := h.photoService.DeletePhoto(c.Request.Context(), tenantID, photoID); err != nil {
		logger.Error("Failed to delete photo", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Photo deleted successfully",
	})
}

// SetPrimary handles PUT /api/v1/gins/:gin_id/photos/:photo_id/primary
func (h *PhotoHandler) SetPrimary(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		c.JSON(400, gin.H{"error": "Tenant not found"})
		return
	}

	ginIDStr := c.Param("gin_id")
	ginID, err := strconv.ParseInt(ginIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid gin ID"})
		return
	}

	photoIDStr := c.Param("photo_id")
	photoID, err := strconv.ParseInt(photoIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid photo ID"})
		return
	}

	// Set primary
	if err := h.photoService.SetPrimaryPhoto(c.Request.Context(), tenantID, ginID, photoID); err != nil {
		logger.Error("Failed to set primary photo", "error", err.Error())
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Primary photo set successfully",
	})
}
