package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainErrors "github.com/yourusername/gin-collection-saas/internal/domain/errors"
)

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// Error sends an error response based on the error type
func Error(c *gin.Context, err error) {
	switch err {
	case domainErrors.ErrNotFound, domainErrors.ErrGinNotFound, domainErrors.ErrTenantNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	case domainErrors.ErrUnauthorized, domainErrors.ErrInvalidCredentials, domainErrors.ErrInvalidToken:
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	case domainErrors.ErrForbidden, domainErrors.ErrTenantSuspended:
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	case domainErrors.ErrLimitReached, domainErrors.ErrFeatureNotAvailable:
		c.JSON(http.StatusForbidden, gin.H{
			"success":          false,
			"error":            err.Error(),
			"upgrade_required": true,
		})
	case domainErrors.ErrConflict, domainErrors.ErrEmailAlreadyExists, domainErrors.ErrSubdomainTaken, domainErrors.ErrBarcodeAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	case domainErrors.ErrInvalidInput, domainErrors.ErrInvalidRating, domainErrors.ErrInvalidFileType:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	case domainErrors.ErrRateLimitExceeded:
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
	}
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, errs map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Validation failed",
		"errors":  errs,
	})
}
