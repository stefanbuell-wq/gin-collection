package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// MaxImageUploadSize is the maximum allowed size for image uploads (50 MB)
	MaxImageUploadSize = 50 << 20 // 50 MB

	// MaxJSONBodySize is the maximum allowed size for JSON request bodies (1 MB)
	MaxJSONBodySize = 1 << 20 // 1 MB
)

// LimitUploadSize returns a middleware that limits the request body size.
// It uses http.MaxBytesReader which:
// - Stops reading once the limit is reached
// - Returns an error "http: request body too large"
// - Protects against memory exhaustion attacks
func LimitUploadSize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// LimitImageUpload is a convenience middleware for image uploads (50 MB limit)
func LimitImageUpload() gin.HandlerFunc {
	return LimitUploadSize(MaxImageUploadSize)
}

// LimitJSONBody is a convenience middleware for JSON bodies (1 MB limit)
func LimitJSONBody() gin.HandlerFunc {
	return LimitUploadSize(MaxJSONBodySize)
}
