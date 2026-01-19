package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/yourusername/gin-collection-saas/internal/infrastructure/cache"
)

const blacklistPrefix = "blacklist:"
const userRevokedPrefix = "user_revoked:"

// TokenBlacklist manages revoked JWT tokens
type TokenBlacklist struct {
	redis *cache.RedisClient
}

// NewTokenBlacklist creates a new blacklist manager
func NewTokenBlacklist(redis *cache.RedisClient) *TokenBlacklist {
	return &TokenBlacklist{redis: redis}
}

// RevokeToken adds a token to the blacklist
// jti: JWT ID (unique identifier)
// expiresAt: When the token naturally expires
func (b *TokenBlacklist) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error {
	if b.redis == nil {
		return nil // Graceful degradation if Redis unavailable
	}

	key := blacklistPrefix + jti
	ttl := time.Until(expiresAt)

	if ttl <= 0 {
		return nil // Token already expired, no need to blacklist
	}

	return b.redis.Set(ctx, key, "revoked", ttl)
}

// IsRevoked checks if a token is blacklisted
func (b *TokenBlacklist) IsRevoked(ctx context.Context, jti string) bool {
	if b.redis == nil {
		return false // Graceful degradation
	}

	key := blacklistPrefix + jti
	exists, _ := b.redis.Exists(ctx, key)
	return exists
}

// RevokeAllUserTokens revokes all tokens for a user (for password change)
// Stores user ID with timestamp, check in middleware
func (b *TokenBlacklist) RevokeAllUserTokens(ctx context.Context, userID int64, since time.Time) error {
	if b.redis == nil {
		return nil
	}

	key := fmt.Sprintf("%s%d", userRevokedPrefix, userID)
	ttl := 24 * time.Hour // JWT expiration time

	return b.redis.Set(ctx, key, since.Unix(), ttl)
}

// IsUserTokenRevoked checks if tokens issued before a certain time are revoked
func (b *TokenBlacklist) IsUserTokenRevoked(ctx context.Context, userID int64, issuedAt time.Time) bool {
	if b.redis == nil {
		return false
	}

	key := fmt.Sprintf("%s%d", userRevokedPrefix, userID)
	val, err := b.redis.Get(ctx, key)
	if err != nil || val == "" {
		return false
	}

	revokedSince, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false
	}

	// Token is revoked if it was issued before the revocation timestamp
	return issuedAt.Unix() < revokedSince
}
