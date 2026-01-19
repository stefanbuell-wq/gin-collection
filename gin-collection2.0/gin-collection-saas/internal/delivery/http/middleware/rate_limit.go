package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/cache"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// inMemoryEntry represents a rate limit counter with expiration
type inMemoryEntry struct {
	count     int64
	expiresAt time.Time
}

// inMemoryStore provides fallback rate limiting when Redis is unavailable
type inMemoryStore struct {
	mu      sync.RWMutex
	entries map[string]*inMemoryEntry
}

// newInMemoryStore creates a new in-memory rate limit store
func newInMemoryStore() *inMemoryStore {
	store := &inMemoryStore{
		entries: make(map[string]*inMemoryEntry),
	}
	// Start cleanup goroutine
	go store.cleanup()
	return store
}

// incr increments a counter and returns the new value
func (s *inMemoryStore) incr(key string, expiry time.Duration) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	entry, exists := s.entries[key]

	if !exists || now.After(entry.expiresAt) {
		// Create new entry or reset expired one
		s.entries[key] = &inMemoryEntry{
			count:     1,
			expiresAt: now.Add(expiry),
		}
		return 1
	}

	entry.count++
	return entry.count
}

// cleanup periodically removes expired entries
func (s *inMemoryStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, entry := range s.entries {
			if now.After(entry.expiresAt) {
				delete(s.entries, key)
			}
		}
		s.mu.Unlock()
	}
}

// Global in-memory store instance
var memoryStore = newInMemoryStore()

// RateLimitMiddleware handles API rate limiting
type RateLimitMiddleware struct {
	redis       *cache.RedisClient
	memoryStore *inMemoryStore
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(redis *cache.RedisClient) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redis:       redis,
		memoryStore: memoryStore,
	}
}

// incrementWithFallback increments a rate limit counter, falling back to in-memory if Redis fails
func (m *RateLimitMiddleware) incrementWithFallback(ctx context.Context, key string, expiry time.Duration) (int64, error) {
	// Try Redis first
	if m.redis != nil {
		count, err := m.redis.Incr(ctx, key)
		if err == nil {
			// Set expiry on first request
			if count == 1 {
				if err := m.redis.Expire(ctx, key, expiry); err != nil {
					logger.Error("Failed to set rate limit expiry", "error", err.Error())
				}
			}
			return count, nil
		}
		logger.Warn("Redis unavailable for rate limiting, falling back to in-memory", "error", err.Error())
	}

	// Fallback to in-memory
	count := m.memoryStore.incr(key, expiry)
	return count, nil
}

// RateLimitByTenant enforces rate limits based on tenant tier
// Uses a sliding window algorithm with Redis
func (m *RateLimitMiddleware) RateLimitByTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant from context (set by auth middleware)
		tenantVal, exists := c.Get("tenant")
		if !exists {
			// No tenant context - skip rate limiting (handled by auth middleware)
			c.Next()
			return
		}

		tenant, ok := tenantVal.(*models.Tenant)
		if !ok {
			c.Next()
			return
		}

		// Get rate limit for tenant's tier
		limits := models.PlanLimitsMap[tenant.Tier]
		rateLimit := limits.APIRateLimit

		// Skip rate limiting for tiers without API access
		if !limits.HasAPIAccess {
			c.Next()
			return
		}

		// Skip if rate limit is 0 (unlimited)
		if rateLimit == 0 {
			c.Next()
			return
		}

		// Create rate limit key (per tenant, per hour)
		now := time.Now()
		windowKey := fmt.Sprintf("ratelimit:tenant:%d:%s", tenant.ID, now.Format("2006010215"))

		ctx := context.Background()

		// Increment counter with fallback
		count, err := m.incrementWithFallback(ctx, windowKey, 65*time.Minute)
		if err != nil {
			logger.Error("Rate limit error", "error", err.Error(), "tenant_id", tenant.ID)
			c.Next()
			return
		}

		// Calculate remaining requests
		remaining := rateLimit - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Calculate reset time (end of current hour window)
		resetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rateLimit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		// Check if rate limit exceeded
		if int(count) > rateLimit {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			logger.Warn("Rate limit exceeded",
				"tenant_id", tenant.ID,
				"tier", tenant.Tier,
				"count", count,
				"limit", rateLimit,
			)

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(429, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "RATE_LIMIT_EXCEEDED",
					"message":     "Rate limit exceeded. Please wait before making more requests.",
					"retry_after": retryAfter,
					"limit":       rateLimit,
					"reset":       resetTime.Unix(),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByIP enforces rate limits by IP address (for unauthenticated endpoints)
func (m *RateLimitMiddleware) RateLimitByIP(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Create rate limit key (per IP, per minute)
		now := time.Now()
		windowKey := fmt.Sprintf("ratelimit:ip:%s:%s", clientIP, now.Format("200601021504"))

		ctx := context.Background()

		// Increment counter with fallback
		count, err := m.incrementWithFallback(ctx, windowKey, 2*time.Minute)
		if err != nil {
			logger.Error("Rate limit error", "error", err.Error(), "ip", clientIP)
			c.Next()
			return
		}

		// Calculate remaining
		remaining := requestsPerMinute - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Reset time (end of current minute)
		resetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())

		// Set headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		// Check if exceeded
		if int(count) > requestsPerMinute {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			logger.Warn("IP rate limit exceeded", "ip", clientIP, "count", count, "limit", requestsPerMinute)

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(429, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "RATE_LIMIT_EXCEEDED",
					"message":     "Too many requests. Please slow down.",
					"retry_after": retryAfter,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitLogin enforces stricter rate limits for login attempts
// 5 attempts per 15 minutes per IP
func (m *RateLimitMiddleware) RateLimitLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Create rate limit key (per IP, 15-minute window)
		now := time.Now()
		windowStart := now.Truncate(15 * time.Minute)
		windowKey := fmt.Sprintf("ratelimit:login:%s:%d", clientIP, windowStart.Unix())

		ctx := context.Background()

		// Increment counter with fallback
		count, err := m.incrementWithFallback(ctx, windowKey, 16*time.Minute)
		if err != nil {
			logger.Error("Login rate limit error", "error", err.Error(), "ip", clientIP)
			c.Next()
			return
		}

		maxAttempts := 5 // Reduced from 10 for better security

		// Calculate remaining
		remaining := maxAttempts - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Reset time
		resetTime := windowStart.Add(15 * time.Minute)

		// Set headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(maxAttempts))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		// Check if exceeded
		if int(count) > maxAttempts {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			logger.Warn("Login rate limit exceeded", "ip", clientIP, "count", count)

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(429, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "TOO_MANY_LOGIN_ATTEMPTS",
					"message":     "Too many login attempts. Please try again later.",
					"retry_after": retryAfter,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitRegistration enforces rate limits for registration attempts
// 3 registrations per hour per IP
func (m *RateLimitMiddleware) RateLimitRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Create rate limit key (per IP, hourly window)
		now := time.Now()
		windowKey := fmt.Sprintf("ratelimit:register:%s:%s", clientIP, now.Format("2006010215"))

		ctx := context.Background()

		// Increment counter with fallback
		count, err := m.incrementWithFallback(ctx, windowKey, 65*time.Minute)
		if err != nil {
			logger.Error("Registration rate limit error", "error", err.Error(), "ip", clientIP)
			c.Next()
			return
		}

		maxAttempts := 3

		// Calculate remaining
		remaining := maxAttempts - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Reset time (end of current hour)
		resetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

		// Set headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(maxAttempts))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		// Check if exceeded
		if int(count) > maxAttempts {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			logger.Warn("Registration rate limit exceeded", "ip", clientIP, "count", count)

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(429, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "TOO_MANY_REGISTRATIONS",
					"message":     "Too many registration attempts. Please try again later.",
					"retry_after": retryAfter,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitPasswordReset enforces rate limits for password reset requests
// 3 requests per hour per Email
func (m *RateLimitMiddleware) RateLimitPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read and restore request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read request body for rate limiting", "error", err.Error())
			c.Next()
			return
		}
		// Restore body for downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse email from body
		var reqBody struct {
			Email string `json:"email"`
		}
		if err := json.Unmarshal(bodyBytes, &reqBody); err != nil || reqBody.Email == "" {
			// If we can't parse email, fall back to IP-based limiting
			clientIP := c.ClientIP()
			m.rateLimitByIPHourly(c, "password_reset_ip", clientIP, 3)
			return
		}

		// Normalize email (lowercase)
		email := strings.ToLower(strings.TrimSpace(reqBody.Email))

		// Create rate limit key (per email, hourly window)
		now := time.Now()
		// Hash email for privacy in Redis keys
		emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))[:16]
		windowKey := fmt.Sprintf("ratelimit:pwreset:%s:%s", emailHash, now.Format("2006010215"))

		ctx := context.Background()

		// Increment counter with fallback
		count, err := m.incrementWithFallback(ctx, windowKey, 65*time.Minute)
		if err != nil {
			logger.Error("Password reset rate limit error", "error", err.Error())
			c.Next()
			return
		}

		maxAttempts := 3

		// Calculate remaining
		remaining := maxAttempts - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// Reset time
		resetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

		// Set headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(maxAttempts))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		// Check if exceeded
		if int(count) > maxAttempts {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			logger.Warn("Password reset rate limit exceeded", "email_hash", emailHash, "count", count)

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			// Return success to prevent email enumeration (security best practice)
			// But internally rate limited - request is not processed
			c.JSON(200, gin.H{
				"success": true,
				"message": "If an account with this email exists, a password reset link has been sent.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByIPHourly enforces rate limits by IP address with hourly window
// Useful for token validation endpoints
func (m *RateLimitMiddleware) RateLimitByIPHourly(requestsPerHour int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		m.rateLimitByIPHourly(c, "hourly", clientIP, requestsPerHour)
	}
}

// rateLimitByIPHourly is a helper method for hourly IP-based rate limiting
func (m *RateLimitMiddleware) rateLimitByIPHourly(c *gin.Context, prefix string, identifier string, maxRequests int) {
	now := time.Now()
	windowKey := fmt.Sprintf("ratelimit:%s:%s:%s", prefix, identifier, now.Format("2006010215"))

	ctx := context.Background()

	// Increment counter with fallback
	count, err := m.incrementWithFallback(ctx, windowKey, 65*time.Minute)
	if err != nil {
		logger.Error("Rate limit error", "error", err.Error(), "identifier", identifier)
		c.Next()
		return
	}

	// Calculate remaining
	remaining := maxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	// Reset time
	resetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

	// Set headers
	c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

	// Check if exceeded
	if int(count) > maxRequests {
		retryAfter := int(time.Until(resetTime).Seconds())
		if retryAfter < 0 {
			retryAfter = 0
		}

		logger.Warn("Hourly IP rate limit exceeded", "identifier", identifier, "count", count, "limit", maxRequests)

		c.Header("Retry-After", strconv.Itoa(retryAfter))
		c.JSON(429, gin.H{
			"success": false,
			"error": gin.H{
				"code":        "RATE_LIMIT_EXCEEDED",
				"message":     "Too many requests. Please try again later.",
				"retry_after": retryAfter,
			},
		})
		c.Abort()
		return
	}

	c.Next()
}
