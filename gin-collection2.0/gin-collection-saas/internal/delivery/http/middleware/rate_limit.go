package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-collection-saas/internal/domain/models"
	"github.com/yourusername/gin-collection-saas/internal/infrastructure/cache"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// RateLimitMiddleware handles API rate limiting
type RateLimitMiddleware struct {
	redis *cache.RedisClient
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(redis *cache.RedisClient) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redis: redis,
	}
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

		// Increment counter
		count, err := m.redis.Incr(ctx, windowKey)
		if err != nil {
			logger.Error("Rate limit Redis error", "error", err.Error(), "tenant_id", tenant.ID)
			// On Redis error, allow request but log warning
			c.Next()
			return
		}

		// Set expiry on first request of the window (1 hour + buffer)
		if count == 1 {
			if err := m.redis.Expire(ctx, windowKey, 65*time.Minute); err != nil {
				logger.Error("Failed to set rate limit expiry", "error", err.Error())
			}
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

		// Increment counter
		count, err := m.redis.Incr(ctx, windowKey)
		if err != nil {
			logger.Error("Rate limit Redis error", "error", err.Error(), "ip", clientIP)
			c.Next()
			return
		}

		// Set expiry on first request
		if count == 1 {
			if err := m.redis.Expire(ctx, windowKey, 2*time.Minute); err != nil {
				logger.Error("Failed to set rate limit expiry", "error", err.Error())
			}
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
func (m *RateLimitMiddleware) RateLimitLogin() gin.HandlerFunc {
	// 10 login attempts per 15 minutes per IP
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Create rate limit key (per IP, 15-minute window)
		now := time.Now()
		windowStart := now.Truncate(15 * time.Minute)
		windowKey := fmt.Sprintf("ratelimit:login:%s:%d", clientIP, windowStart.Unix())

		ctx := context.Background()

		// Increment counter
		count, err := m.redis.Incr(ctx, windowKey)
		if err != nil {
			logger.Error("Login rate limit Redis error", "error", err.Error(), "ip", clientIP)
			c.Next()
			return
		}

		// Set expiry on first request (16 minutes to cover the window)
		if count == 1 {
			if err := m.redis.Expire(ctx, windowKey, 16*time.Minute); err != nil {
				logger.Error("Failed to set login rate limit expiry", "error", err.Error())
			}
		}

		maxAttempts := 10

		// Check if exceeded
		if int(count) > maxAttempts {
			resetTime := windowStart.Add(15 * time.Minute)
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
