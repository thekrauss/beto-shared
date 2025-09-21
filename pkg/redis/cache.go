package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

// stocke une string avec TTL
func SetString(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := Client.Set(ctx, key, value, ttl).Err(); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to set string in Redis")
	}
	return nil
}

// récupère une string
func GetString(ctx context.Context, key string) (string, error) {
	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(errors.CodeNotFound, "key not found in Redis")
	} else if err != nil {
		return "", errors.Wrap(err, errors.CodeInternal, "failed to get string from Redis")
	}
	return val, nil
}

// stocke une struct en JSON
func SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to marshal JSON for Redis")
	}
	return SetString(ctx, key, string(data), ttl)
}

// récupère une struct depuis JSON
func GetJSON(ctx context.Context, key string, v interface{}) error {
	val, err := GetString(ctx, key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), v); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to unmarshal JSON from Redis")
	}
	return nil
}

// supprime une clé
func Delete(ctx context.Context, key string) error {
	if err := Client.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to delete key from Redis")
	}
	return nil
}
