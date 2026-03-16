// Purpose: Redis client — connection setup + JWT blacklist operations
// Ref: ADR-002 — JWT blacklist via Redis; 30-minute idle timeout enforcement
// Dev must implement: Connect(url string), IsBlacklisted(token), Blacklist(token, ttl)
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect creates a Redis client and verifies connectivity.
func Connect(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis url parse error: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}

// BlacklistToken adds a JWT to the blacklist with the given TTL.
// Key format: "blacklist:{token_jti}"
func BlacklistToken(ctx context.Context, rdb *redis.Client, jti string, ttl time.Duration) error {
	return rdb.Set(ctx, "blacklist:"+jti, "1", ttl).Err()
}

// IsBlacklisted checks if a JWT jti is in the blacklist.
func IsBlacklisted(ctx context.Context, rdb *redis.Client, jti string) (bool, error) {
	val, err := rdb.Exists(ctx, "blacklist:"+jti).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}
