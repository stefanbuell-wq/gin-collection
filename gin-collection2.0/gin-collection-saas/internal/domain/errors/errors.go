package errors

import "errors"

// Domain-specific errors
var (
	// General errors
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden - insufficient permissions")
	ErrInvalidInput      = errors.New("invalid input")
	ErrConflict          = errors.New("resource already exists")
	ErrInternalServer    = errors.New("internal server error")

	// Authentication errors
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrEmailNotVerified    = errors.New("email not verified")

	// Tenant errors
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrTenantSuspended     = errors.New("tenant account is suspended")
	ErrSubdomainTaken      = errors.New("subdomain already taken")

	// Subscription/Tier errors
	ErrLimitReached        = errors.New("tier limit reached - upgrade required")
	ErrFeatureNotAvailable = errors.New("feature not available in current tier")
	ErrSubscriptionInactive = errors.New("subscription is not active")

	// Gin-specific errors
	ErrGinNotFound         = errors.New("gin not found")
	ErrBarcodeAlreadyExists = errors.New("barcode already exists in your collection")
	ErrInvalidRating       = errors.New("rating must be between 1 and 5")

	// Photo errors
	ErrPhotoLimitReached   = errors.New("photo limit reached for this gin")
	ErrStorageLimitReached = errors.New("storage limit reached")
	ErrInvalidFileType     = errors.New("invalid file type - only images allowed")
	ErrFileTooLarge        = errors.New("file too large")

	// Multi-user errors (Enterprise)
	ErrMultiUserNotAllowed = errors.New("multi-user feature requires Enterprise tier")
	ErrUserNotInTenant     = errors.New("user does not belong to this tenant")

	// API errors (Enterprise)
	ErrAPIAccessNotAllowed = errors.New("API access requires Enterprise tier")
	ErrInvalidAPIKey       = errors.New("invalid API key")
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error           string `json:"error"`
	UpgradeRequired bool   `json:"upgrade_required,omitempty"`
	Code            string `json:"code,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error, upgradeRequired bool) *ErrorResponse {
	return &ErrorResponse{
		Error:           err.Error(),
		UpgradeRequired: upgradeRequired,
	}
}
