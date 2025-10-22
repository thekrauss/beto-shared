package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

var Client *redis.Client

func Init(host string, port int, password string, db int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to ping Redis")
	}

	return nil
}

func Ping(ctx context.Context) error {
	if Client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	status, err := Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	if status != "PONG" {
		return fmt.Errorf("redis returned unexpected status: %s", status)
	}
	return nil
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
