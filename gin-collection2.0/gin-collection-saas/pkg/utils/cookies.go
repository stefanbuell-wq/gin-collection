package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// AccessTokenCookieName is the name of the access token cookie
	AccessTokenCookieName = "access_token"
	// RefreshTokenCookieName is the name of the refresh token cookie
	RefreshTokenCookieName = "refresh_token"
)

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite http.SameSite
}

// SetAuthCookies sets both access and refresh token cookies
func SetAuthCookies(c *gin.Context, accessToken, refreshToken string, cfg *CookieConfig, accessExpiry, refreshExpiry time.Duration) {
	// Set SameSite mode
	c.SetSameSite(cfg.SameSite)

	// Access Token Cookie - sent with all requests
	c.SetCookie(
		AccessTokenCookieName,
		accessToken,
		int(accessExpiry.Seconds()),
		"/",          // Path: all routes
		cfg.Domain,   // Domain
		cfg.Secure,   // Secure: HTTPS only in production
		true,         // HttpOnly: not accessible via JavaScript
	)

	// Refresh Token Cookie - only sent to auth endpoints
	c.SetCookie(
		RefreshTokenCookieName,
		refreshToken,
		int(refreshExpiry.Seconds()),
		"/api/v1/auth", // Path: only auth routes
		cfg.Domain,
		cfg.Secure,
		true,
	)
}

// SetAccessTokenCookie sets only the access token cookie (for refresh)
func SetAccessTokenCookie(c *gin.Context, accessToken string, cfg *CookieConfig, accessExpiry time.Duration) {
	c.SetSameSite(cfg.SameSite)
	c.SetCookie(
		AccessTokenCookieName,
		accessToken,
		int(accessExpiry.Seconds()),
		"/",
		cfg.Domain,
		cfg.Secure,
		true,
	)
}

// ClearAuthCookies removes both auth cookies by setting them to expire immediately
func ClearAuthCookies(c *gin.Context, cfg *CookieConfig) {
	c.SetSameSite(cfg.SameSite)

	// Clear access token cookie
	c.SetCookie(
		AccessTokenCookieName,
		"",
		-1, // Negative MaxAge deletes the cookie
		"/",
		cfg.Domain,
		cfg.Secure,
		true,
	)

	// Clear refresh token cookie
	c.SetCookie(
		RefreshTokenCookieName,
		"",
		-1,
		"/api/v1/auth",
		cfg.Domain,
		cfg.Secure,
		true,
	)
}

// GetAccessTokenFromCookie extracts the access token from the request cookie
func GetAccessTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie(AccessTokenCookieName)
}

// GetRefreshTokenFromCookie extracts the refresh token from the request cookie
func GetRefreshTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie(RefreshTokenCookieName)
}
