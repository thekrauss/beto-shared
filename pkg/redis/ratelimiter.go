package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

// implementes a distributed rate limiter (fixed window)
func AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {

	//increments a Redis counter
	val, err := Client.Incr(ctx, key).Result()
	if err != nil {
		return false, errors.Wrap(err, errors.CodeInternal, "failed to incr Redis key for rate limiting")
	}

	//if it is the first increment -> set TTL
	if val == 1 {
		if err := Client.Expire(ctx, key, window).Err(); err != nil {
			return false, errors.Wrap(err, errors.CodeInternal, "failed to set TTL on Redis key for rate limiting")
		}
	}

	//authorizes if counter <= limit
	return val <= int64(limit), nil
}

// returns how many requests are left before blocking
func GetRemainingRequests(ctx context.Context, key string, limit int) (int, error) {
	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return limit, nil //key non-existent (not yet used)
	} else if err != nil {
		return 0, errors.Wrap(err, errors.CodeInternal, "failed to get key for rate limiting")
	}

	count, _ := strconv.Atoi(val)
	if count >= limit {
		return 0, nil
	}
	return limit - count, nil
}
