package utils

import (
	"bytes"
	"fmt"
)

// Allowed image content types
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// ValidateImageMagicBytes validates file content by checking magic bytes (file signature)
// This prevents attacks where malicious files are renamed with image extensions
func ValidateImageMagicBytes(data []byte) (string, error) {
	if len(data) < 12 {
		return "", fmt.Errorf("file too small to determine type (minimum 12 bytes required)")
	}

	// JPEG: FF D8 FF (Start of Image marker)
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "image/jpeg", nil
	}

	// PNG: 89 50 4E 47 0D 0A 1A 0A (PNG signature)
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "image/png", nil
	}

	// GIF: 47 49 46 38 (GIF87a or GIF89a)
	if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}) {
		return "image/gif", nil
	}

	// WebP: RIFF....WEBP (RIFF container with WEBP format)
	if bytes.HasPrefix(data, []byte{0x52, 0x49, 0x46, 0x46}) && // RIFF
		len(data) >= 12 && string(data[8:12]) == "WEBP" {
		return "image/webp", nil
	}

	return "", fmt.Errorf("unsupported or invalid image file type")
}

// GetExtensionForContentType returns the appropriate file extension for a content type
func GetExtensionForContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

// IsAllowedImageType checks if the content type is in the allowed list
func IsAllowedImageType(contentType string) bool {
	return allowedImageTypes[contentType]
}
