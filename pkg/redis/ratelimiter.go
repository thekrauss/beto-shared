package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

// implémente un rate limiter distribué (fixed window)
func AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {

	// incrémente un compteur Redis
	val, err := Client.Incr(ctx, key).Result()
	if err != nil {
		return false, errors.Wrap(err, errors.CodeInternal, "failed to incr Redis key for rate limiting")
	}

	// si c'est le premier incrément -> definir TTL
	if val == 1 {
		if err := Client.Expire(ctx, key, window).Err(); err != nil {
			return false, errors.Wrap(err, errors.CodeInternal, "failed to set TTL on Redis key for rate limiting")
		}
	}

	// autorise si compteur <= limite
	return val <= int64(limit), nil
}

// retourne combien il reste de requêtes avant le blocage
func GetRemainingRequests(ctx context.Context, key string, limit int) (int, error) {
	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return limit, nil // clé inexistante () pas encore utilisé)
	} else if err != nil {
		return 0, errors.Wrap(err, errors.CodeInternal, "failed to get key for rate limiting")
	}

	count, _ := strconv.Atoi(val)
	if count >= limit {
		return 0, nil
	}
	return limit - count, nil
}
